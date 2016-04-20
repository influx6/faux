package flux

// Inspired and Learned from William Kennedy's(Ardan Labs) Infinite Queue

import (
	"errors"
	"sync"
)

// PressureStream provides a higher api for handling pressure requests of two levels (data and errors),there by providing a simple but strong foundation for higher level constructs
type PressureStream struct {
	Signals, Errors chan interface{}

	// internal queue system for back-pressuing and delivery guarantee
	signalQueue *Queue
	errorQueue  *Queue
}

// NewPressureStream returns a new PressureStream
func NewPressureStream() *PressureStream {
	data := make(chan interface{})
	errs := make(chan interface{})
	return BuildPressureStream(data, errs)
}

// BuildPressureStream returns a new PressureStream instance
func BuildPressureStream(dataSignals, errorSignals chan interface{}) *PressureStream {
	ps := PressureStream{
		Signals: dataSignals,
		Errors:  errorSignals,

		signalQueue: NewQueue(dataSignals),
		errorQueue:  NewQueue(errorSignals),
	}

	return &ps
}

// RemainingSignals returns the current size of signal queue
func (ps *PressureStream) RemainingSignals() int {
	return ps.errorQueue.Length()
}

// RemainingErrors returns the current size of error queue
func (ps *PressureStream) RemainingErrors() int {
	return ps.errorQueue.Length()
}

// SendSignal delivers data into the data channel
func (ps *PressureStream) SendSignal(d interface{}) {
	ps.signalQueue.Enqueue(d)
}

// SendError delivers data into the data channel
func (ps *PressureStream) SendError(d error) {
	ps.errorQueue.Enqueue(d)
}

// Close closes the PressureStream but ensures the function only returns control when the both back queues have been closed
func (ps *PressureStream) Close() {
	ps.signalQueue.Close()
	ps.errorQueue.Close()
}

// ErrQueueEmpty is returned when the queue has no more elements
var ErrQueueEmpty = errors.New("Queue: Is Empty")

// Queue is a simple queue with the capability of handling infinite receivals and retrieving with a more control system
type Queue struct {
	Deq chan interface{}

	// internal manager channels
	enq  chan interface{}
	done chan bool

	// wait-group to ensure proper closing of operations/no ophaned go-routines
	wg sync.WaitGroup

	ro sync.Mutex

	// internal buffer
	bu *Buffer
}

// NewQueue returns a new pressure queue
func NewQueue(deq chan interface{}) *Queue {
	q := Queue{
		Deq: deq,

		enq:  make(chan interface{}),
		done: make(chan bool),
		bu:   NewBuffer(),
	}

	q.wg.Add(1)

	go q.manage()

	return &q
}

//Length returns the length of items within the queue's buffer
func (q *Queue) Length() int {
	return q.bu.Length()
}

// Close sets the internal operations channels to a close state and waits till the buffer operations are complete to return, hence ensuring the last items where delivered
func (q *Queue) Close() {
	close(q.enq)
	q.wg.Wait()
	close(q.Deq)
}

// Enqueue adds a new item into the queue's buffer with a guarantee that it was received before returned and manages concurrency by locking the function until the data as being received
func (q *Queue) Enqueue(item interface{}) {
	q.ro.Lock()
	{
		q.enq <- item
		<-q.done
	}
	q.ro.Unlock()
}

// manage entails the operations of the infinite queue ensuring all data is received and properly delivered to the queue outpoint
func (q *Queue) manage() {
	defer q.wg.Done()

	for {

		if q.bu.Length() == 0 {
			select {
			case item, open := <-q.enq:
				if !open {
					// possibly closed,return
					return
				}
				q.bu.Enqueue(item)
				q.done <- true
			}
		}

		if q.bu.Length() > 0 {

			// cache the first item to preseve it else if we dequeue it will be lost when an enqueue is received
			firstItem, err := q.bu.Peek()

			if err != nil {
				// we know there is was an item in the queue,hence we retake the operation to ensure we didnt miss
				continue
			}

			// we order up the queueu-dequeue operations and ensure the first item was delivered to the out-channel provided
			select {
			// wait for an item to be delivered into the queue or till the queue is closed
			case item, open := <-q.enq:
				if !open {
					// possibly closed,return
					return
				}
				q.bu.Enqueue(item)
				q.done <- true

			// send out the first item as guaranteed
			case q.Deq <- firstItem:
				q.bu.Dequeue()
			}
		}

	}
}

// ErrBufferEmpty is returned when a op is performed with an empty buffer
var ErrBufferEmpty = errors.New("Buffer: Is Empty")

// Buffer provides a infinite receive space for handling incoming data, providing a good back pressure mechanism
type Buffer struct {
	rw     sync.RWMutex
	buffer []interface{}
}

// NewBuffer returns a new Buffer instance
func NewBuffer() *Buffer {
	return &Buffer{}
}

// Length returns the current size of element in the buffer
func (b *Buffer) Length() int {
	var ln int

	b.rw.RLock()
	{
		ln = len(b.buffer)
	}
	b.rw.RUnlock()

	return ln
}

// Peek returns the first value in the buffer or an error if empty
func (b *Buffer) Peek() (interface{}, error) {
	var ln int
	var v interface{}

	b.rw.RLock()
	{
		if ln = len(b.buffer); ln > 0 {
			v = b.buffer[0]
		}
	}
	b.rw.RUnlock()

	if ln == 0 {
		return nil, ErrBufferEmpty
	}

	return v, nil
}

// Clear emties the buffer of current content
func (b *Buffer) Clear() {
	b.rw.Lock()
	{
		b.buffer = nil
	}
	b.rw.Unlock()
}

// Enqueue adds an item into the buffers end
func (b *Buffer) Enqueue(item interface{}) {
	b.rw.Lock()
	{
		b.buffer = append(b.buffer, item)
	}
	b.rw.Unlock()
}

// Dequeue removes an item into the buffers front
func (b *Buffer) Dequeue() (interface{}, error) {
	var ln int
	var item interface{}

	b.rw.Lock()
	{
		if ln = len(b.buffer); ln > 0 {
			item = b.buffer[0]
			b.buffer = b.buffer[1:]
		}
	}
	b.rw.Unlock()

	if ln == 0 {
		return nil, ErrBufferEmpty
	}

	return item, nil
}
