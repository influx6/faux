package flux

import (
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

const (
	//signals for handling workers
	addWorker    = 1
	removeWorker = 2
)

// Work defines an interface for a work to be performed
type Work interface {
	Work(context interface{}, id int)
}

// doWork defines a work to be routed for execution
type doWork struct {
	Context interface{}
	Work    Work
}

// PoolConfig defines the configuration details for a workpool
type PoolConfig struct {
	MaxWorkers     int64
	MinWorkers     int64
	MetricInterval func() time.Duration
	MetricHandler  func(PoolStat)
}

// PoolStat defines the stat returned when checking health of pool
type PoolStat struct {
	//date of the state
	Stamp time.Time

	MaxWorkers int64
	MinWorkers int64
	//total current workers goroutined
	Workers int64

	//total executed works
	Executed int64
	//pending work in pool
	Pending int64
	//active work in pool
	Active int64
}

var (
	//ErrInvalidMinWorkers defines an error when the minimum worker provided is below or at 0
	ErrInvalidMinWorkers = errors.New("Invalid minimum worker value")
	//ErrInvalidMaxWorkers is returned when the max worker value is incorrect
	ErrInvalidMaxWorkers = errors.New("Invalid maximum worker value")
	//ErrInvalidAddRequest is returned when the pool has reached its maximum worker value and cant add anymore
	ErrInvalidAddRequest = errors.New("Pool is at maximum worker efficiency")
	//ErrWorkRequestDenied is returned when the pool is unable to service a task
	ErrWorkRequestDenied = errors.New("Pool unable to accept task")
)

//WorkPool defines a pool for handling workers
type WorkPool struct {
	*PoolConfig
	name string

	currentWorkers int64

	//total number of updates to goroutines
	updatePending int64

	activeWork   int64
	pendingWork  int64
	executedWork int64

	//kill the pool
	kill     chan struct{}
	shutdown chan struct{}

	//worker management
	commands chan int

	tasks chan doWork

	dohealth sync.Mutex

	//used to managed the current total routines working and allows gracefull shutdown
	man sync.WaitGroup
}

// NewPool creates a new pool instance
func NewPool(context interface{}, name string, config *PoolConfig) (*WorkPool, error) {
	if config.MinWorkers <= 0 {
		return nil, ErrInvalidMinWorkers
	}

	if config.MaxWorkers <= 0 || config.MaxWorkers < config.MinWorkers {
		return nil, ErrInvalidMaxWorkers
	}

	pool := WorkPool{
		PoolConfig: config,
		name:       name,

		shutdown: make(chan struct{}),
		kill:     make(chan struct{}),
		commands: make(chan int),

		tasks: make(chan doWork),
	}

	pool.manage()

	pool.Add(context, int(pool.MinWorkers))

	//  pool.Add()
	return &pool, nil
}

// Stat returns the current stat of the pool
func (w *WorkPool) Stat() PoolStat {
	return PoolStat{
		Stamp:      time.Now(),
		MaxWorkers: atomic.LoadInt64(&w.MaxWorkers),
		MinWorkers: atomic.LoadInt64(&w.MinWorkers),
		Workers:    atomic.LoadInt64(&w.currentWorkers),
		Executed:   atomic.LoadInt64(&w.executedWork),
		Pending:    atomic.LoadInt64(&w.pendingWork),
		Active:     atomic.LoadInt64(&w.activeWork),
	}
}

// DoWait adds a task into the pool within a given duration if not accepted by then it returns with  an error
func (w *WorkPool) DoWait(c interface{}, wo Work, t time.Duration) error {
	dow := doWork{c, wo}

	w.measureHealth()

	atomic.AddInt64(&w.pendingWork, 1)

	select {
	case w.tasks <- dow:
		atomic.AddInt64(&w.pendingWork, -1)
		return nil
	case <-time.After(t):
		atomic.AddInt64(&w.pendingWork, -1)
		return ErrWorkRequestDenied
	}

}

// Do adds a task into the pool and blocks until it is accepted
func (w *WorkPool) Do(c interface{}, wo Work) {
	dow := doWork{c, wo}

	w.measureHealth()

	atomic.AddInt64(&w.pendingWork, 1)
	{
		w.tasks <- dow
	}
	atomic.AddInt64(&w.pendingWork, -1)
}

// Name returns the name of workpool
func (w *WorkPool) Name() string {
	return w.name
}

func (w *WorkPool) measureHealth() {
	//is there still updates pending,if so return
	if atomic.LoadInt64(&w.updatePending) == 0 {
		return
	}

	w.dohealth.Lock()
	defer w.dohealth.Unlock()

	stat := w.Stat()

	if stat.Active == 0 && stat.Pending == 0 && (stat.Workers > w.MinWorkers) {
		w.Reset(w.name, int(w.MinWorkers))
	}

	if stat.Active == stat.Workers && stat.Workers < w.MaxWorkers {

		ko := int(float64(stat.Workers) * .20)

		if ko == 0 {
			ko = 1
		}

		if int(stat.Workers)+ko > int(w.MaxWorkers) {
			ko = int(w.MaxWorkers - stat.Workers)
		}

		w.Add(w.name, ko)
	}

}

// Reset resets the total allowed routines to the given value
func (w *WorkPool) Reset(context interface{}, rmw int) {
	if rmw < 0 {
		rmw = 0
	}

	cur := int(atomic.LoadInt64(&w.currentWorkers))
	w.Add(context, rmw-cur)
}

// Add adds or removes a worker from the active routines
func (w *WorkPool) Add(context interface{}, do int) error {
	if do == 0 {
		return ErrInvalidAddRequest
	}

	cur := int(atomic.LoadInt64(&w.currentWorkers))

	cmd := addWorker
	if do < 0 {
		do = cur - do
		cmd = removeWorker
	}

	atomic.AddInt64(&w.updatePending, int64(do))

	for i := 0; i < do; i++ {
		w.commands <- cmd
	}

	return nil
}

func (w *WorkPool) work(id int) {
	go func() {
		cur := atomic.LoadInt64(&w.currentWorkers)
		maxo := atomic.LoadInt64(&w.MaxWorkers)

		if cur > maxo {
			atomic.StoreInt64(&w.MaxWorkers, cur)
		}

		//decrement that an addition is done
		atomic.AddInt64(&w.updatePending, -1)
	worker:
		for {
			select {
			case do := <-w.tasks:
				atomic.AddInt64(&w.activeWork, 1)
				w.execute(id, do)
				atomic.AddInt64(&w.executedWork, 1)
				atomic.AddInt64(&w.activeWork, -1)
			case <-w.kill:
				break worker
			}
		}

		atomic.AddInt64(&w.currentWorkers, -1)
		//decrement that the removal is done
		atomic.AddInt64(&w.updatePending, -1)
		w.man.Done()
	}()
}

func (w *WorkPool) execute(id int, dw doWork) {
	RecoveryHandler(fmt.Sprintf("WorkPool Worker #Id: %d", id), func() error {
		dw.Work.Work(dw.Context, id)
		return nil
	})
}

// Shutdown ends the workpool and all workers and jobs
func (w *WorkPool) Shutdown() {
	close(w.shutdown)
	w.man.Wait()
}

func (w *WorkPool) nextMetric() time.Duration {
	var metric time.Duration

	if w.MetricInterval != nil {
		metric = w.MetricInterval()
	}

	if metric > 0 && metric < time.Second {
		metric = time.Second
	}

	return metric
}

func (w *WorkPool) manage() {
	w.man.Add(1)
	go func() {

		metrico := w.nextMetric()
		metricTimer := new(time.Timer)

		if metrico > 0 {
			metricTimer = time.NewTimer(metrico)
		}

		for {
			select {
			case <-w.shutdown:
				row := atomic.LoadInt64(&w.currentWorkers)

				//are we already below mark
				if row <= 0 {
					return
				}

				for i := 0; i < int(row); i++ {
					w.kill <- struct{}{}
				}
				w.man.Done()
				return

			case cmd, ok := <-w.commands:
				if !ok {
					return
				}
				switch cmd {
				case addWorker:
					row := atomic.LoadInt64(&w.currentWorkers)
					maxw := atomic.LoadInt64(&w.MaxWorkers)

					//check if we are equal or above watermark
					if row >= maxw {
						return
					}

					id := atomic.AddInt64(&w.currentWorkers, 1)

					//add to the control group for workers
					w.man.Add(1)

					//add the work
					w.work(int(id))

				case removeWorker:
					row := atomic.LoadInt64(&w.currentWorkers)
					minw := atomic.LoadInt64(&w.MinWorkers)

					//check if we re below minimu workers if so,ignore
					if minw > row {
						return
					}

					w.kill <- struct{}{}
				}
			case <-metricTimer.C:
				mo := w.nextMetric()

				if w.MetricHandler != nil {
					w.MetricHandler(w.Stat())
				}
				metricTimer.Reset(mo)
			}
		}
	}()
}
