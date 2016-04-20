// Package sumex provides a stream idiomatic api for go using goroutines and
// channels whilst still allowing to maintain syncronicity of operations
// using channel pipeline strategies.
// Unlike most situations, panics in sumex are caught and sent out as errors.
package sumex

import (
	"bytes"
	"errors"
	"fmt"
	"sync"
	"sync/atomic"

	"github.com/influx6/faux/context"
	"github.com/influx6/faux/panics"
	"github.com/satori/go.uuid"
)

//==============================================================================

// Log defines event logger that allows us to record events for a specific
// action that occured.
type Log interface {
	Log(context interface{}, name string, message string, data ...interface{})
	Error(context interface{}, name string, err error, message string, data ...interface{})
}

//==============================================================================

var events eventlog

// eventlog provides a concrete implementation of a logger.
type eventlog struct{}

// Log logs all standard log reports.
func (l eventlog) Log(context interface{}, name string, message string, data ...interface{}) {}

// Error logs all error reports.
func (l eventlog) Error(context interface{}, name string, err error, message string, data ...interface{}) {
}

//==============================================================================

// Proc defines a interface for processors for sumex streams.
type Proc interface {
	Do(context.Context, error, interface{}) (interface{}, error)
}

//==============================================================================

//==============================================================================

// Stat defines the current capacity workings of
type Stat struct {
	TotalWorkersRunning int64
	TotalWorkers        int64
	Pending             int64
	Completed           int64
	Closed              int64
}

// Stream define a pipeline operation for applying operations to
// data streams.
type Stream interface {
	Stats() Stat
	UUID() string
	Logs() Log
	Shutdown()
	Stream(Stream) Stream
	CloseNotify() <-chan struct{}
	Error(context.Context, error)
	Data(context.Context, interface{})
}

//==============================================================================

// payload defines the payload being injected into the system.
type payload struct {
	err error
	d   interface{}
	ctx context.Context
}

// dataSink provides a interface{} channel type.
type dataSink chan *payload

//==============================================================================

// Worker defines a configuration interface for sumex.Streams.
type Worker struct {
	Max func() int
	Min func() int
	Log Log
}

var minFn = func() int { return 1 }
var maxFn = func() int { return 2 }

// New returns a new Stream compliant instance.
func New(w *Worker, p Proc) Stream {
	if w == nil {
		w = (*Worker)(nil)
	}

	if w.Min == nil || w.Min() <= 0 {
		w.Min = minFn
	}

	if w.Max == nil || w.Max() <= 0 {
		w.Max = maxFn
	}

	if w.Log == nil {
		w.Log = events
	}

	sm := stream{
		Worker: w,
		uuid:   uuid.NewV4().String(),
		proc:   p,
		data:   make(dataSink),
		ender:  make(chan struct{}),
		nc:     make(chan struct{}),
		ctx:    context.New(),
	}

	// initialize the total data workers needed.
	for i := 0; i < sm.Min(); i++ {
		go sm.worker()
	}

	return &sm
}

//==============================================================================

// stream defines a structure that implements the Stream interface, providing
// a basic building block for stream operation.
type stream struct {
	*Worker
	uuid string

	closed               int64
	processed            int64
	pending              int64
	shutdownAfterpending int64
	workersUp            int64
	proc                 Proc
	ctx                  context.Context

	data  dataSink
	ender chan struct{}
	nc    chan struct{}

	wg   sync.WaitGroup
	pl   sync.RWMutex
	pubs []Stream // list of listeners.
}

// Stats reports the current operational status of the streamer
func (s *stream) Stats() Stat {
	return Stat{
		TotalWorkers: atomic.LoadInt64(&s.workersUp),
		Pending:      atomic.LoadInt64(&s.pending),
		Completed:    atomic.LoadInt64(&s.processed),
		Closed:       atomic.LoadInt64(&s.closed),
	}
}

// Log returns the internal logger for this stream.
func (s *stream) Logs() Log {
	return s.Log
}

// Shutdown closes the data and error channels.
func (s *stream) Shutdown() {
	s.Log.Log("Sumex.Stream", "Shutdown", "Started : Shutdown Requested")
	if atomic.LoadInt64(&s.closed) > 0 {
		s.Log.Log("Sumex.Stream", "Stats", "Completed : Shutdown Request : Previously Done")
		return
	}

	defer close(s.nc)

	atomic.StoreInt64(&s.closed, 1)

	atomic.AddInt64(&s.pending, 1)
	{
		close(s.data)
	}
	atomic.AddInt64(&s.pending, -1)

	s.wg.Wait()

	s.Log.Log("Sumex.Stream", "Shutdown", "Completed : Shutdown Requested")
}

// CloseNotify returns a chan used to shutdown the close of a stream.
func (s *stream) CloseNotify() <-chan struct{} {
	return s.nc
}

// Stream provides a pipe which adds a new receiver of data to the provided
// stream. Returns the supplied stream receiver.
func (s *stream) Stream(sm Stream) Stream {
	s.pl.Lock()
	defer s.pl.Unlock()
	s.pubs = append(s.pubs, sm)
	return sm
}

// UUID returns a UUID string for the given stream.
func (s *stream) UUID() string {
	return s.uuid
}

// Data sends in data for execution by the stream into its data channel.
// It allows providing an optional context which would be passed into the
// internal processor else using the default context of the stream.
func (s *stream) Data(ctx context.Context, d interface{}) {
	if atomic.LoadInt64(&s.closed) > 0 {
		return
	}

	if ctx == nil {
		ctx = s.ctx
	}

	s.Log.Log("sumex.Stream", "Data", "Started : Data Recieved : %s", fmt.Sprintf("%+v", d))
	atomic.AddInt64(&s.pending, 1)
	{
		s.data <- &payload{ctx: ctx, d: d}
	}
	atomic.AddInt64(&s.pending, -1)
	s.Log.Log("sumex.Stream", "Data", "Completed")
}

// Error pipes in a new data for execution by the stream
// into its err channel.
func (s *stream) Error(ctx context.Context, e error) {
	if atomic.LoadInt64(&s.closed) != 0 {
		return
	}

	if ctx == nil {
		ctx = s.ctx
	}

	s.Log.Error("sumex.Stream", "Error", e, "Started : Error Recieved : %s", fmt.Sprintf("%+v", e))
	atomic.AddInt64(&s.pending, 1)
	{
		s.data <- &payload{ctx: ctx, err: e}
	}
	atomic.AddInt64(&s.pending, -1)
	s.Log.Log("sumex.Stream", "Error", "Completed")
}

// worker initializes data workers for stream.
func (s *stream) worker() {
	defer s.wg.Done()
	defer atomic.AddInt64(&s.workersUp, -1)

	s.wg.Add(1)
	atomic.AddInt64(&s.workersUp, 1)

loop:
	for {
		select {
		case <-s.ender:
			break loop
		case load, ok := <-s.data:
			if !ok {
				break loop
			}

			panics.Defer(func() {
				res, err := s.proc.Do(load.ctx, load.err, load.d)
				s.Log.Log("sumex.Stream", "worker", "Info : Res : { Response: %+s, Error: %+s}", res, err)

				atomic.AddInt64(&s.processed, 1)

				s.pl.RLock()
				defer s.pl.RUnlock()

				if err != nil {
					for _, sm := range s.pubs {
						go sm.Error(load.ctx, err)
					}
					return
				}

				for _, sm := range s.pubs {
					go sm.Data(load.ctx, res)
				}

			}, func(d *bytes.Buffer) {
				s.Logs().Error("sumex.Stream", "worker", errors.New("Panic"), "Panic : %+s", d.Bytes())
			})
		}
	}

	s.Log.Log("sumex.Stream", "worker", "Info : Goroutine : Shutdown")
}

//==========================================================================================

// ProcHandler defines a base function type for sumex streams.
type ProcHandler func(context.Context, error, interface{}) (interface{}, error)

// Do creates a new stream from a function provided
func Do(sm Stream, w *Worker, h ProcHandler) Stream {

	if h == nil {
		panic("nil ProcHandler")
	}

	ms := New(w, doworker{h})

	if sm != nil {
		sm.Stream(ms)
	}

	return ms
}

//==========================================================================================

// Identity returns a stream which returns what it recieves as output.
func Identity(w *Worker, l Log) Stream {
	return Do(nil, w, func(ctx context.Context, err error, d interface{}) (interface{}, error) {
		return d, err
	})
}

//==========================================================================================

// doworker provides a proxy structure matching the Proc interface.
// Calls the provided ProcHandler underneath its Do method.
type doworker struct {
	p ProcHandler
}

// Do provide a proxy caller for the handler registered to it.
func (do doworker) Do(ctx context.Context, err error, d interface{}) (interface{}, error) {
	return do.p(ctx, err, d)
}

//==========================================================================================

// Receive returns a receive-only blocking chan which returns
// all data items from a giving stream.
// The returned channel gets closed along  with the stream
func Receive(sm Stream) (<-chan interface{}, Stream) {
	mc := make(chan interface{})

	ms := Do(sm, nil, func(ctx context.Context, _ error, d interface{}) (interface{}, error) {
		mc <- d
		return nil, nil
	})

	go func() {
		<-sm.CloseNotify()
		ms.Shutdown()
		close(mc)
	}()

	return mc, ms
}

// ReceiveError returns a receive-only blocking chan which returns all error
// items from a giving stream.
// The returned channel gets closed along  with the stream
func ReceiveError(sm Stream) (<-chan error, Stream) {
	mc := make(chan error)

	ms := Do(sm, nil, func(ctx context.Context, e error, _ interface{}) (interface{}, error) {
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

	return mc, ms
}

//==========================================================================================
