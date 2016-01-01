// Package sumex provides a stream idiomatic api for go using goroutines and
// channels whilst still allowing to maintain syncronicity of operations
// using channel pipeline strategies.
// Unlike most situations, panics in sumex are caught and sent out as errors.
package sumex

import (
	"sync"
	"sync/atomic"

	"github.com/influx6/faux/panics"
)

// Streams define a pipeline operation for applying operations to
// data streams.
type Streams interface {
	Inject(interface{})
	InjectError(error)
	Stream(Streams) Streams
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
	closed int64
	proc   Proc

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
		proc: p,
		data: make(dataSink),
		err:  make(errorSink),
		nc:   make(chan struct{}),
	}

	// startup the error worker. We only need one.
	go sm.initEW()

	// initialize the total data workers needed.
	for i := 0; i < w; i++ {
		go sm.initDW()
	}

	return &sm
}

// Shutdown closes the data and error channels.
func (s *Stream) Shutdown() {
	if atomic.LoadInt64(&s.closed) == 0 {
		return
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

// Inject pipes in a new data for execution by the stream into its data channel.
func (s *Stream) Inject(d interface{}) {
	if atomic.LoadInt64(&s.closed) == 0 {
		s.wg.Add(1)
		s.data <- d
	}
}

// InjectError pipes in a new data for execution by the stream
// into its err channel.
func (s *Stream) InjectError(e error) {
	if atomic.LoadInt64(&s.closed) == 0 {
		s.wg.Add(1)
		s.err <- e
	}
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
	for {

		if atomic.LoadInt64(&s.closed) != 0 {
			break
		}

		d, ok := <-s.data
		if !ok {
			break
		}

		panics.Guard(func() error {
			defer s.wg.Done()
			res, err := s.proc.Do(d, nil)
			if err != nil {
				s.sendError(err)
				return nil
			}

			s.send(res)
			return nil
		})

	}
}

// initEW initializes error workers for stream.
func (s *Stream) initEW() {
	for {

		if atomic.LoadInt64(&s.closed) != 0 {
			break
		}

		e, ok := <-s.err
		if !ok {
			break
		}

		panics.Guard(func() error {
			defer s.wg.Done()
			res, err := s.proc.Do(nil, e)
			if err != nil {
				s.sendError(err)
				return nil
			}

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
