// Package sumex provides a stream idiomatic api for go using goroutines and
// channels whilst still allowing to maintain syncronicity of operations
// using channel pipeline strategies.
// Unlike most situations, panics in sumex are caught and sent out as errors.
package sumex

import (
	"bytes"
	"fmt"
	"sync"
	"sync/atomic"

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
	Logger() Log
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
	log  Log
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
func New(w int, l Log, p Proc) Streams {
	if w <= 0 {
		w = 1
	}

	if l == nil {
		l = events
	}

	sm := Stream{
		log:     l,
		uuid:    uuid.NewV4().String(),
		workers: w,
		proc:    p,
		data:    make(dataSink),
		err:     make(errorSink),
		nc:      make(chan struct{}),
	}

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

// Logger returns the internal logger for this stream.
func (s *Stream) Logger() Log {
	return s.log
}

// Shutdown closes the data and error channels.
func (s *Stream) Shutdown() {
	s.log.Log("Sumex.Streams", "Shutdown", "Started : Shutdown Requested")
	if atomic.LoadInt64(&s.closed) == 0 {
		s.log.Log("Sumex.Streams", "Stats", "Info : Shutdown Request : Previously Done")
		return
	}

	if pen := atomic.LoadInt64(&s.pending); pen > 0 {
		s.log.Log("Sumex.Streams", "Shutdown", "Info : Pending : %d", pen)
		atomic.StoreInt64(&s.shutdownAfterpending, 1)
		atomic.StoreInt64(&s.closed, 1)
	}

	s.log.Log("Sumex.Streams", "Shutdown", "Started : WaitGroup.Wait()")
	s.wg.Wait()
	s.log.Log("Sumex.Streams", "Shutdown", "Completed : WaitGroup.Wait()")

	close(s.data)
	close(s.err)
	close(s.nc)
	atomic.StoreInt64(&s.closed, 1)
	s.log.Log("Sumex.Streams", "Shutdown", "Completed : Shutdown Requested")
}

// CloseNotify returns a chan used to shutdown the close of a stream.
func (s *Stream) CloseNotify() <-chan struct{} {
	return s.nc
}

// Stream provides a pipe which adds a new receiver of data to the provided
// stream. Returns the supplied stream receiver.
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
	if atomic.LoadInt64(&s.closed) > 0 {
		return
	}

	s.log.Log("sumex.Streams", "Inject", "Started : Data Recieved : %s", fmt.Sprintf("%+v", d))
	atomic.AddInt64(&s.pending, 1)
	{
		s.data <- d
	}
	atomic.AddInt64(&s.pending, -1)
	s.log.Log("sumex.Streams", "Inject", "Completed")
}

// InjectError pipes in a new data for execution by the stream
// into its err channel.
func (s *Stream) InjectError(e error) {
	if atomic.LoadInt64(&s.closed) != 0 {
		return
	}

	s.log.Error("sumex.Streams", "InjectError", e, "Started : Error Recieved : %s", fmt.Sprintf("%+v", e))
	atomic.AddInt64(&s.pending, 1)
	{
		s.err <- e
	}
	atomic.AddInt64(&s.pending, -1)
	s.log.Log("sumex.Streams", "InjectError", "Completed")
}

// sendError delivers the an error to all registered streamers.
func (s *Stream) sendError(e error) {
	s.log.Log("sumex.Streams", "sendError", "Started : %+s", e)
	s.pl.RLock()
	for _, sm := range s.pubs {
		go sm.InjectError(e)
	}
	s.pl.RUnlock()
	s.log.Log("sumex.Streams", "sendError", "Completed")
}

// send delivers the data to all registered streamers.
func (s *Stream) send(d interface{}) {
	s.log.Log("sumex.Streams", "send", "Started : %s", fmt.Sprintf("%+v", d))
	s.pl.RLock()
	for _, sm := range s.pubs {
		go sm.Inject(d)
	}
	s.pl.RUnlock()
	s.log.Log("sumex.Streams", "send", "Completed")
}

// initDW initializes data workers for stream.
func (s *Stream) initDW() {
	defer func() {
		s.log.Log("sumex.Streams", "initDW", "Started : Goroutine : Shutdown")
		s.wg.Done()
		atomic.AddInt64(&s.workersUp, -1)
		s.log.Log("sumex.Streams", "initDW", "Completed : Goroutine : Shutdown")
	}()

	s.wg.Add(1)
	atomic.AddInt64(&s.workersUp, 1)

	for {

		if atomic.LoadInt64(&s.closed) > 0 && atomic.LoadInt64(&s.shutdownAfterpending) < 1 {
			s.log.Log("sumex.Streams", "initDW", "Closed : ShutDown Request Pending")
			break
		}

		// if we are to shutdown after pending and pending is now zero then
		if atomic.LoadInt64(&s.shutdownAfterpending) > 0 && atomic.LoadInt64(&s.pending) < 1 {
			s.log.Log("sumex.Streams", "initDW", "Closed : ShutDown Request : No Pending")
			break
		}

		var d interface{}
		var e error
		var ok bool

		select {
		case d, ok = <-s.data:
		case e, ok = <-s.err:
		}

		if !ok {
			break
		}

		if e != nil {
			s.log.Error("sumex.Streams", "initDW", e, "Received Data")
		}

		panics.Defer(func() {
			res, err := s.proc.Do(d, e)
			if err != nil {
				atomic.AddInt64(&s.processed, 1)
				s.sendError(err)
				return
			}

			if res != nil {
				atomic.AddInt64(&s.processed, 1)
				s.send(res)
			}
		}, func(d *bytes.Buffer) {
			fmt.Println(d.String())
		})

	}
}

//==========================================================================================

// ProcHandler defines a base function type for sumex streams.
type ProcHandler func(interface{}, error) (interface{}, error)

// Do creates a new stream from a function provided
func Do(sm Streams, workers int, h ProcHandler) Streams {

	if h == nil {
		panic("nil ProcHandler")
	}

	var log Log

	if sm == nil {
		log = events
	}

	ms := New(workers, log, doworker{h})

	if sm != nil {
		sm.Stream(ms)
	}

	return ms
}

//==========================================================================================

// Identity returns a stream which returns what it recieves as output.
func Identity(w int, l Log) Streams {
	return Do(nil, w, func(d interface{}, err error) (interface{}, error) {
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
func (do doworker) Do(d interface{}, err error) (interface{}, error) {
	return do.p(d, err)
}

//==========================================================================================

// Receive returns a receive-only blocking chan which returns
// all data items from a giving stream.
// The returned channel gets closed along  with the stream
func Receive(sm Streams) (<-chan interface{}, Streams) {
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

	return mc, ms
}

// ReceiveError returns a receive-only blocking chan which returns all error
// items from a giving stream.
// The returned channel gets closed along  with the stream
func ReceiveError(sm Streams) (<-chan error, Streams) {
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

	return mc, ms
}

//==========================================================================================
