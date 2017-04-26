// Package workers provides a worker idiomatic api for go using goroutines and
// channels whilst still allowing to maintain syncronicity of operations
// using channel pipeline strategies.
// Unlike most situations, panics in sumex are caught and sent out as errors.
package workers

import (
	"bytes"
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

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

// Handler defines a interface for processors for sumex workers.
type Handler interface {
	Do(context.Context, error, interface{}) (interface{}, error)
}

//==============================================================================

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

// Schedule defines a type which takes a time.Duration and returns a new
// duration.
type Schedule func(time.Duration) time.Duration

// BasicSchedule defines a simple scheduler which returns either the passed in
// duration if not less than zero else returns 1s.
func BasicSchedule(dt time.Duration) time.Duration {
	if dt <= 0 {
		return 1000 * time.Millisecond
	}

	return dt
}

//==============================================================================

// Config defines a configuration interface for workers.workers.
type Config struct {
	Max              int           // Minimum allowed workers.
	Min              int           // Maximum allowed workers.
	Log              Log           // Log to use for logging.
	SkipError        bool          // Used to tell the worker to pass down errors without calling Handler.
	RelaxScheduler   Schedule      // Used normally or when low pressure to keep relaxed check times.
	ChokeScheduler   Schedule      // Used when there is high work pressure to increase check times.
	CheckDuration    time.Duration // Initial duration before manager checks state of workers.
	MaxCheckDuration time.Duration // Maximum allow duration growth for state checks.
}

// Worker define a pipeline operation for applying operations to
// data workers.
type Worker interface {
	Stats() Stat
	UUID() string
	Logs() Log
	Shutdown()
	Next(Worker) Worker
	CloseNotify() <-chan struct{}
	Error(context.Context, error)
	Data(context.Context, interface{})
}

// New returns a new worker compliant instance.
func New(c Config, p Handler) Worker {
	if c.Min <= 0 {
		c.Min = 1
	}

	if c.Max <= 0 {
		c.Max = 2
	}

	if c.Log == nil {
		c.Log = events
	}

	if c.CheckDuration <= 0 {
		c.CheckDuration = 1 * time.Millisecond
	}

	if c.ChokeScheduler == nil {
		c.ChokeScheduler = BasicSchedule
	}

	if c.RelaxScheduler == nil {
		c.RelaxScheduler = BasicSchedule
	}

	sm := worker{
		lastStat:      time.Now(),
		config:        &c,
		uuid:          uuid.NewV4().String(),
		Handler:       p,
		data:          make(dataSink),
		ender:         make(chan struct{}),
		nc:            make(chan struct{}),
		mn:            make(chan struct{}),
		startDuration: c.CheckDuration,
		ctx:           context.New(),
	}

	// initialize the total data workers needed.
	for i := 0; i < sm.config.Min; i++ {
		go sm.worker()
	}

	go sm.manage()

	return &sm
}

//==============================================================================

// worker defines a structure that implements the worker interface, providing
// a basic building block for worker operation.
type worker struct {
	config *Config
	uuid   string

	closed               int64
	active               int64
	processed            int64
	pending              int64
	shutdownAfterpending int64
	workersUp            int64
	Handler              Handler
	ctx                  context.Context
	startDuration        time.Duration

	data     dataSink
	ender    chan struct{}
	nc       chan struct{}
	mn       chan struct{}
	lastStat time.Time

	workerGroup sync.WaitGroup

	pl   sync.RWMutex
	pubs []Worker // list of listeners.
}

// Stat defines the current capacity workings of
type Stat struct {
	TotalWorkersRunning int64         `json:"total_workers_running"`
	TotalWorkers        int64         `json:"total_workers"`
	Pending             int64         `json:"pending_tasks"`
	Completed           int64         `json:"total_completed_tasks"`
	Closed              int64         `json:"total_removed_workers"`
	ElapsedStat         time.Duration `json:"elapsed_stat"`
	Time                time.Time     `json:"time"`
}

// String returns a representation of the giving stat.
func (s Stat) String() string {
	return fmt.Sprintf(`
		 Current Time: %s
		 Total Elapsed Time from Last Stat: %s
		 Total Current Workers: %d
		 Total Active Workers: %d
		 Total Pending Task: %d
		 Total Completed Task: %d
		 Total Closed Workers: %d
	`, s.Time.UTC(), s.ElapsedStat, s.TotalWorkers, s.TotalWorkersRunning, s.Pending, s.Completed, s.Closed)
}

// Stats reports the current operational status of the streamer
func (s *worker) Stats() Stat {
	now := time.Now()

	elpased := now.Sub(s.lastStat)
	s.lastStat = now

	return Stat{
		TotalWorkersRunning: atomic.LoadInt64(&s.active),
		TotalWorkers:        atomic.LoadInt64(&s.workersUp),
		Pending:             atomic.LoadInt64(&s.pending),
		Completed:           atomic.LoadInt64(&s.processed),
		Closed:              atomic.LoadInt64(&s.closed),
		ElapsedStat:         elpased,
		Time:                now,
	}
}

// Log returns the internal logger for this worker.
func (s *worker) Logs() Log {
	return s.config.Log
}

// Shutdown closes the data and error channels.
func (s *worker) Shutdown() {
	s.config.Log.Log(s.uuid, "Shutdown", "Started : Shutdown Requested")
	if atomic.LoadInt64(&s.closed) > 0 {
		s.config.Log.Log(s.uuid, "Stats", "Completed : Shutdown Request : Previously Done")
		return
	}

	atomic.StoreInt64(&s.closed, 1)
	defer close(s.nc)

	close(s.mn)
	close(s.ender)

	s.workerGroup.Wait()

	s.config.Log.Log(s.uuid, "Shutdown", "Completed : Shutdown Requested")
}

// CloseNotify returns a chan used to shutdown the close of a worker.
func (s *worker) CloseNotify() <-chan struct{} {
	return s.nc
}

// Next provides a pipe which adds a new receiver of data to the provided
// worker. Returns the supplied worker receiver.
func (s *worker) Next(sm Worker) Worker {
	s.pl.Lock()
	defer s.pl.Unlock()
	s.pubs = append(s.pubs, sm)
	return sm
}

// UUID returns a UUID string for the given worker.
func (s *worker) UUID() string {
	return s.uuid
}

// Data sends in data for execution by the worker into its data channel.
// It allows providing an optional context which would be passed into the
// internal processor else using the default context of the worker.
func (s *worker) Data(ctx context.Context, d interface{}) {
	if atomic.LoadInt64(&s.closed) > 0 {
		return
	}

	if ctx == nil {
		ctx = s.ctx
	}

	// s.dataGroup.Add(1)

	s.config.Log.Log(s.uuid, "Data", "Started : Data Recieved : %s", fmt.Sprintf("%+v", d))
	atomic.AddInt64(&s.pending, 1)
	{
		s.data <- &payload{ctx: ctx, d: d}
	}
	atomic.AddInt64(&s.pending, -1)
	s.config.Log.Log(s.uuid, "Data", "Completed")

	// s.dataGroup.Done()
}

// Error pipes in a new data for execution by the worker
// into its err channel.
func (s *worker) Error(ctx context.Context, e error) {
	if atomic.LoadInt64(&s.closed) != 0 {
		return
	}

	if ctx == nil {
		ctx = s.ctx
	}

	// s.dataGroup.Add(1)

	s.config.Log.Error(s.uuid, "Error", e, "Started : Error Recieved : %s", fmt.Sprintf("%+v", e))
	atomic.AddInt64(&s.pending, 1)
	{
		s.data <- &payload{ctx: ctx, err: e}
	}
	atomic.AddInt64(&s.pending, -1)
	s.config.Log.Log(s.uuid, "Error", "Completed")
	// s.dataGroup.Done()

}

func (s *worker) manage() {
	{

		defer func() {
			s.config.Log.Log(s.uuid, "worker", "Info : Worker Manager Shutdown")
		}()

		clock := time.NewTimer(s.config.CheckDuration)

	manageloop:
		for {

			select {
			case <-s.mn:
				break manageloop
			case <-clock.C:

				// Collect the current stats.
				stat := s.Stats()
				s.config.Log.Log(s.uuid, "worker", "Info : Stat : {%+s}", stat)

				// If we have more workers more than requests then shave off extra luggage.
				if stat.TotalWorkers > stat.Pending {

					unUsed := int(stat.Pending - stat.TotalWorkers)
					if unUsed < 0 {
						unUsed = unUsed * -1
					}

					// If pending is zero, just remove all the needed
					if stat.Pending < 1 && unUsed > s.config.Min {
						rmCount := unUsed - s.config.Min
						s.config.Log.Log(s.uuid, "worker", "Info : Removing Total Workers[%d]", rmCount)

						if atomic.LoadInt64(&s.closed) > 0 {
							break
						}

						for i := 0; i < rmCount; i++ {
							s.ender <- struct{}{}
						}

						continue
					}

					var wasted int

					// If we have more pending that mimumium workers, then just remove
					// the un-used workers instead. No needed killing proficiency here.
					if int(stat.Pending) > s.config.Min {

						// Reduce the uneeded workers.
						wasted = int(stat.TotalWorkers) - unUsed
					} else {

						// We practically have workers than even minimum needed for tasks.
						// Return all workers to minimum needed.
						wasted = int(stat.TotalWorkers) - s.config.Min
					}

					if wasted > 0 {
						s.config.Log.Log(s.uuid, "worker", "Info : Removing Total Workers[%d]", wasted)

						if atomic.LoadInt64(&s.closed) > 0 {
							break
						}

						for i := 0; i < wasted; i++ {
							s.ender <- struct{}{}
						}
					}

					continue
				}

				// Sleep for 1ms then call another stat to see if any change.
				// time.Sleep(20 * time.Nanosecond)

				newStat := s.Stats()
				s.config.Log.Log(s.uuid, "worker", "Info : New Stat : {%+s}", newStat)

				var newWorkersAdded bool

				// We check if there is an increase in the workers or if the workers
				// are still unable to handle the workload, when this is the case we
				// fire out that the needed workers double of the load difference to
				// account for the case that more requesters could be increasing.
				if newStat.Pending >= stat.Pending || newStat.TotalWorkers < newStat.Pending {
					newWorkersAdded = true

					load := int(newStat.TotalWorkers - newStat.Pending)

					// Incase we are dealing with a -1 negative range number, then mux it
					// by -1.
					if load < 0 {
						load = load * -1
					}

					s.config.Log.Log(s.uuid, "worker", "Info : Current Total Load[%d]", load)

					// We only ever allow a maximum growth range regardless of load.
					if load > int(s.config.Max) {
						load = s.config.Max
					} else {
						// If growing at twice the load rate is lower than max workers,
						// then grow accoding to that.
						if doubleLoad := load * 2; doubleLoad < s.config.Max {
							load = doubleLoad
						}
					}

					s.config.Log.Log(s.uuid, "worker", "Info : Add Total Workers[%d]", load)
					for i := 0; i < load; i++ {
						go s.worker()
					}
				}

				var scheduler Schedule

				// If we added new workers then we have a large task pool. Use the
				// Relax scheduler to allow a longer wait check to allow current max
				// pool to reduce.
				if newWorkersAdded {
					scheduler = s.config.ChokeScheduler
				} else {
					scheduler = s.config.RelaxScheduler
				}

				// Recheck the duration clock and resets the clock.
				if s.config.CheckDuration >= s.config.MaxCheckDuration {
					s.config.CheckDuration = s.startDuration
				} else {
					s.config.CheckDuration = scheduler(s.config.CheckDuration)
				}

				s.config.Log.Log(s.uuid, "worker", "Info : Using New Check Duration[%s]", s.config.CheckDuration)
				if !clock.Reset(s.config.CheckDuration) {
					clock = time.NewTimer(s.config.CheckDuration)
				}

				continue

			}
		}
	}
}

// worker initializes data workers for worker.
func (s *worker) worker() {
	defer s.workerGroup.Done()

	atomic.AddInt64(&s.workersUp, 1)
	s.workerGroup.Add(1)

loop:
	for {
		select {
		case <-s.ender:
			atomic.AddInt64(&s.workersUp, -1)
			break loop
		case load, ok := <-s.data:
			if !ok {
				break loop
			}

			atomic.AddInt64(&s.active, 1)
			{
				panics.Defer(func() {
					s.pl.RLock()
					defer s.pl.RUnlock()

					if s.config.SkipError {
						if load.err != nil {
							for _, sm := range s.pubs {
								go sm.Error(load.ctx, load.err)
							}
							return
						}
					}

					res, err := s.Handler.Do(load.ctx, load.err, load.d)
					s.config.Log.Log(s.uuid, "worker", "Info : Res : { Response: %+s, Error: %+s}", res, err)

					atomic.AddInt64(&s.processed, 1)

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
					s.Logs().Error(s.uuid, "worker", errors.New("Panic"), "Panic : %+s", d.Bytes())
				})
			}
			atomic.AddInt64(&s.active, -1)
		}
	}

	s.config.Log.Log(s.uuid, "worker", "Info : Goroutine : Shutdown")
}

//==========================================================================================

// Handle defines a base function type for sumex workers.
type Handle func(context.Context, error, interface{}) (interface{}, error)

// Do creates a new worker from a function provided
func Do(sm Worker, w Config, h Handle) Worker {

	if h == nil {
		panic("nil ProcHandler")
	}

	ms := New(w, doworker{h})

	if sm != nil {
		sm.Next(ms)
	}

	return ms
}

//==========================================================================================

// Identity returns a worker which returns what it recieves as output.
func Identity(w Config, l Log) Worker {
	return Do(nil, w, func(ctx context.Context, err error, d interface{}) (interface{}, error) {
		return d, err
	})
}

//==========================================================================================

// doworker provides a proxy structure matching the Handler interface.
// Calls the provided ProcHandler underneath its Do method.
type doworker struct {
	p Handle
}

// Do provide a proxy caller for the handler registered to it.
func (do doworker) Do(ctx context.Context, err error, d interface{}) (interface{}, error) {
	return do.p(ctx, err, d)
}

//==========================================================================================

// Receive returns a receive-only blocking chan which returns
// all data items from a giving worker.
// The returned channel gets closed along  with the worker
func Receive(sm Worker) (<-chan interface{}, Worker) {
	mc := make(chan interface{})

	ms := Do(sm, Config{}, func(ctx context.Context, _ error, d interface{}) (interface{}, error) {
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
// items from a giving worker.
// The returned channel gets closed along  with the worker
func ReceiveError(sm Worker) (<-chan error, Worker) {
	mc := make(chan error)

	ms := Do(sm, Config{}, func(ctx context.Context, e error, _ interface{}) (interface{}, error) {
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
