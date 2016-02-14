// Package mque provides a argument aware queue that allows creating pubsub
// type queues where unless the type expected by the callback is published to
// the queue or unless the subscriber expects no argument will it be runned.
package mque

import (
	"errors"
	"reflect"
	"sync"
	"sync/atomic"

	"github.com/influx6/faux/loop"
	"github.com/influx6/faux/reflection"
)

//==============================================================================

// Qu defines a queue interface that defines the general queue registery.
type Qu interface {
	Q(interface{}) loop.Looper
	Run(interface{})
	Flush()
}

// New returns a new implementer of Qu.
func New() Qu {
	mq := MQue{}
	return &mq
}

//==============================================================================

// MQue defines a callback queue, that accept only one argument functions.
type MQue struct {
	l      sync.RWMutex
	muxers []*mqueSub
}

// Flush ends the queue listeners
func (m *MQue) Flush() {
	m.l.Lock()
	m.muxers = nil
	m.l.Unlock()
}

// Run applies the argument against the queues callbacks.
func (m *MQue) Run(mx interface{}) {
	m.l.RLock()
	defer m.l.RUnlock()
	for _, mux := range m.muxers {
		if atomic.LoadInt64(&mux.alive) > 0 {
			mux.Run(mx)
		}
	}
}

//==============================================================================

// ErrInvalid is returned when the provided type is not a function type.
var ErrInvalid = errors.New("Invalid Func Type")

// ErrTooMuchArgs is returned when the argument of the function is more than one.
var ErrTooMuchArgs = errors.New("Expected One argument Function")

// Q adds a new function type into the queue.
func (m *MQue) Q(mx interface{}) loop.Looper {
	if !reflection.IsFuncType(mx) {
		return nil
	}

	m.l.RLock()
	id := len(m.muxers)
	m.l.RUnlock()

	tm, _ := reflection.FuncValue(mx)

	var tu reflect.Type
	var bit int

	args, _ := reflection.GetFuncArgumentsType(mx)
	if size := len(args); size > 0 {
		tu = args[0]
		bit = 1
	}

	mq := mqueSub{
		id:    id,
		mx:    mx,
		tm:    tm,
		qu:    m,
		am:    tu,
		has:   bit,
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
	mx    interface{}
	tm    reflect.Value
	am    reflect.Type
	qu    *MQue
	has   int
	alive int64
}

// End removes this subscriber for the queue.
func (m *mqueSub) End(f ...func()) {
	atomic.StoreInt64(&m.alive, 0)
	for _, fx := range f {
		fx()
	}
}

// Run recevies the argument and
func (m *mqueSub) Run(d interface{}) {
	if m.has == 0 {
		m.tm.Call(nil)
		return
	}

	var configVal reflect.Value
	ctype := reflect.TypeOf(d)

	if !ctype.AssignableTo(m.am) {
		if !ctype.ConvertibleTo(m.am) {
			return
		}

		vum := reflect.ValueOf(d)
		configVal = vum.Convert(m.am)
	} else {
		configVal = reflect.ValueOf(d)
	}

	m.tm.Call([]reflect.Value{configVal})
}

//==============================================================================
