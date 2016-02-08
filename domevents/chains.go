package domevents

import (
	"errors"
	"sync"
	"sync/atomic"

	"github.com/influx6/faux/loop"
)

//==============================================================================

// Qu defines a queue interface that defines the general queue registery.
type Qu interface {
	Q(func(Event)) loop.Looper
	Run(Event)
	Flush()
}

// NewQu returns a new implementer of Qu.
func NewQu() Qu {
	mq := MQue{}
	return &mq
}

//==============================================================================

// MQue defines a callback queue, that accept only one argument functions.
type MQue struct {
	l      sync.RWMutex
	muxers []*mqueSub
}

// Run applies the argument against the queues callbacks.
func (m *MQue) Run(e Event) {
	m.l.RLock()
	defer m.l.RUnlock()
	for _, mux := range m.muxers {
		if atomic.LoadInt64(&mux.alive) > 0 {
			mux.Run(e)
		}
	}
}

// Flush ends the queue listeners
func (m *MQue) Flush() {
	m.l.Lock()
	m.muxers = nil
	m.l.Unlock()
}

//==============================================================================

// ErrInvalid is returned when the provided type is not a function type.
var ErrInvalid = errors.New("Invalid Func Type")

// ErrTooMuchArgs is returned when the argument of the function is more than one.
var ErrTooMuchArgs = errors.New("Expected One argument Function")

// Q adds a new function type into the queue.
func (m *MQue) Q(mx func(Event)) loop.Looper {
	m.l.RLock()
	id := len(m.muxers)
	m.l.RUnlock()

	mq := mqueSub{
		id:    id,
		mx:    mx,
		alive: 1,
	}

	m.l.Lock()
	m.muxers = append(m.muxers, &mq)
	m.l.Unlock()

	return &mq
}

//==============================================================================

// mqueSub defines a queue subscriber attached to a specific queue.
type mqueSub struct {
	id    int
	mx    func(Event)
	alive int64
}

// End removes this subscriber for the queue.
func (m *mqueSub) End() {
	atomic.StoreInt64(&m.alive, 0)
}

// Run recevies the argument and
func (m *mqueSub) Run(e Event) {
	m.mx(e)
}

//==============================================================================
