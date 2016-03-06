// Package sumex provides a stream idiomatic api for go using goroutines and
// channels whilst still allowing to maintain syncronicity of operations
// using channel pipeline strategies.
// Unlike most situations, panics in sumex are caught and sent out as errors.
package sumex

import (
	"sync"
	"sync/atomic"

	"github.com/influx6/faux/panics"
	"github.com/satori/go.uuid"
)

// Stat defines the current capacity workings of
type Stat struct {
	TotalWorkersRunning int64
	TotalWorkers        int64
	Pending             int64
	Completed           int64
	Closed              int64
}

// Streams define a pipeline operation for applying operations to
// data streams.
type Streams interface {
	Inject(interface{})
	InjectError(error)
	Stream(Streams) Streams
	Stats() Stat
	UUID() string
	Shutdown()
	CloseNotify() <-chan struct{}
}

// Proc defines a interface for processors for sumex streams.
type Proc interface {
	Do(interface{}, error) (interface{}, error)
}

// dataSink provides a interface{} channel type.
type dataSink chan interface{}

// errorSink provies an error channel type.
type errorSink chan error

// Stream defines a structure that implements the Streams interface, providing
// a basic building block for stream operation.
type Stream struct {
	uuid string

	closed               int64
	processed            int64
	pending              int64
	shutdownAfterpending int64
	workersUp            int64
	workers              int
	proc                 Proc

	data dataSink
	err  errorSink
	nc   chan struct{}

	wg   sync.WaitGroup
	pl   sync.RWMutex
	pubs []Streams // list of listeners.
}

// New returns a new Streams compliant instance.
func New(w int, p Proc) Streams {
	if w <= 0 {
		w = 1
	}

	sm := Stream{
		uuid:    uuid.NewV4().String(),
		workers: w,
		proc:    p,
		data:    make(dataSink),
		err:     make(errorSink),
		nc:      make(chan struct{}),
	}

	// startup the error worker. We only need one.
	go sm.initEW()

	// initialize the total data workers needed.
	for i := 0; i < w; i++ {
		go sm.initDW()
	}

	return &sm
}

// Stats reports the current operational status of the streamer
func (s *Stream) Stats() Stat {
	return Stat{
		TotalWorkersRunning: atomic.LoadInt64(&s.workersUp),
		TotalWorkers:        int64(s.workers + 1),
		Pending:             atomic.LoadInt64(&s.pending),
		Completed:           atomic.LoadInt64(&s.processed),
		Closed:              atomic.LoadInt64(&s.closed),
	}
}

// Shutdown closes the data and error channels.
func (s *Stream) Shutdown() {
	if atomic.LoadInt64(&s.closed) == 0 {
		return
	}

	if atomic.LoadInt64(&s.pending) != 0 {
		atomic.StoreInt64(&s.shutdownAfterpending, 1)
		atomic.StoreInt64(&s.closed, 1)
	}

	s.wg.Wait()

	close(s.data)
	close(s.err)
	close(s.nc)
	atomic.StoreInt64(&s.closed, 1)
}

// CloseNotify returns a chan used to shutdown the close of a stream.
func (s *Stream) CloseNotify() <-chan struct{} {
	return s.nc
}

// Stream provides a pipe which adds a new receiver of data to the provided
// stream. Returns the supplied stream.
func (s *Stream) Stream(sm Streams) Streams {
	s.pl.Lock()
	defer s.pl.Unlock()
	s.pubs = append(s.pubs, sm)
	return sm
}

// UUID returns a UUID string for the given stream.
func (s *Stream) UUID() string {
	return s.uuid
}

// Inject pipes in a new data for execution by the stream into its data channel.
func (s *Stream) Inject(d interface{}) {
	if atomic.LoadInt64(&s.closed) != 0 {
		return
	}

	atomic.AddInt64(&s.pending, 1)
	{
		s.data <- d
	}
	atomic.AddInt64(&s.pending, -1)
}

// InjectError pipes in a new data for execution by the stream
// into its err channel.
func (s *Stream) InjectError(e error) {
	if atomic.LoadInt64(&s.closed) != 0 {
		return
	}

	atomic.AddInt64(&s.pending, 1)
	{
		s.err <- e
	}
	atomic.AddInt64(&s.pending, -1)
}

// sendError delivers the an error to all registered streamers.
func (s *Stream) sendError(e error) {
	s.pl.RLock()
	for _, sm := range s.pubs {
		sm.InjectError(e)
	}
	s.pl.RUnlock()
}

// send delivers the data to all registered streamers.
func (s *Stream) send(d interface{}) {
	s.pl.RLock()
	for _, sm := range s.pubs {
		sm.Inject(d)
	}
	s.pl.RUnlock()
}

// initDW initializes data workers for stream.
func (s *Stream) initDW() {
	defer func() {
		s.wg.Done()
		atomic.AddInt64(&s.workersUp, -1)
	}()

	s.wg.Add(1)
	atomic.AddInt64(&s.workersUp, 1)

	for {

		if atomic.LoadInt64(&s.closed) != 0 && atomic.LoadInt64(&s.shutdownAfterpending) == 0 {
			break
		}

		// if we are to shutdown after pending and pending is now zero then
		if atomic.LoadInt64(&s.shutdownAfterpending) != 0 && atomic.LoadInt64(&s.pending) == 0 {
			break
		}

		d, ok := <-s.data
		if !ok {
			break
		}

		panics.Guard(func() error {
			res, err := s.proc.Do(d, nil)
			if err != nil {
				atomic.AddInt64(&s.processed, 1)
				s.sendError(err)
				return nil
			}

			if res == nil {
				return nil
			}

			atomic.AddInt64(&s.processed, 1)
			s.send(res)
			return nil
		})

	}
}

// initEW initializes error workers for stream.
func (s *Stream) initEW() {
	defer func() {
		s.wg.Done()
		atomic.AddInt64(&s.workersUp, -1)
	}()

	s.wg.Add(1)
	atomic.AddInt64(&s.workersUp, 1)

	for {

		if atomic.LoadInt64(&s.closed) != 0 && atomic.LoadInt64(&s.shutdownAfterpending) == 0 {
			break
		}

		// if we are to shutdown after pending and pending is now zero then
		if atomic.LoadInt64(&s.shutdownAfterpending) != 0 && atomic.LoadInt64(&s.pending) == 0 {
			break
		}

		e, ok := <-s.err
		if !ok {
			break
		}

		panics.Guard(func() error {
			res, err := s.proc.Do(nil, e)
			if err != nil {
				atomic.AddInt64(&s.processed, 1)
				s.sendError(err)
				return nil
			}

			if res == nil {
				return nil
			}

			atomic.AddInt64(&s.processed, 1)
			s.send(res)
			return nil
		})

	}
}

// ProcHandler defines a base function type for sumex streams.
type ProcHandler func(interface{}, error) (interface{}, error)

// Do creates a new stream from a function provided
func Do(sm Streams, workers int, h ProcHandler) Streams {
	if h == nil {
		panic("nil ProcHandler")
	}
	ms := New(workers, doworker{h})
	sm.Stream(ms)
	return ms
}

// doworker provides a proxy structure matching the Proc interface.
// Calls the provided ProcHandler underneath its Do method.
type doworker struct {
	p ProcHandler
}

// Do provide a proxy caller for the handler registered to it.
func (do doworker) Do(d interface{}, err error) (interface{}, error) {
	return do.p(d, err)
}

// Receive returns a receive-only blocking chan which returns
// all data items from a giving stream.
// The returned channel gets closed along  with the stream
func Receive(sm Streams) <-chan interface{} {
	mc := make(chan interface{})

	ms := Do(sm, 1, func(d interface{}, _ error) (interface{}, error) {
		mc <- d
		return nil, nil
	})

	go func() {
		<-sm.CloseNotify()
		ms.Shutdown()
		close(mc)
	}()

	return mc
}

// ReceiveError returns a receive-only blocking chan which returns all error
// items from a giving stream.
// The returned channel gets closed along  with the stream
func ReceiveError(sm Streams) <-chan error {
	mc := make(chan error)

	ms := Do(sm, 1, func(_ interface{}, e error) (interface{}, error) {
		if e != nil {
			mc <- e
		}
		return nil, nil
	})

	go func() {
		<-sm.CloseNotify()
		ms.Shutdown()
		close(mc)
	}()

	return mc
}
