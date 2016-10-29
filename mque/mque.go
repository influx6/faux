// Package mque provides a argument aware queue that allows creating pubsub
// type queues where unless the type expected by the callback is published to
// the queue or unless the subscriber expects no argument will it be runned.
package mque

import (
	"reflect"
	"sync"

	"github.com/influx6/faux/reflection"
)

//==============================================================================

// End defines an interface which exposes a End function.
type End interface {
	End()
	AddEnd(func())
}

// New returns a new implementer of Qu.
func New() *MQue {
	var any mqueSub
	return &MQue{any: &any}
}

//==============================================================================

// MQue defines a callback queue, that accept only one argument functions.
type MQue struct {
	l      sync.RWMutex
	muxers []*mqueSub
	any    *mqueSub
}

// Flush ends the queue listeners.
func (m *MQue) Flush() {
	m.l.Lock()
	m.muxers = nil
	m.l.Unlock()
}

// Run applies the argument against the queues callbacks.
func (m *MQue) Run(val interface{}) {
	m.l.RLock()
	defer m.l.RUnlock()

	ctype := reflect.TypeOf(val)

	// Run the any callbacks.
	m.any.Run(val, ctype)

	// Check and Run those who match.
	for _, mux := range m.muxers {
		if mux.CanRun(ctype) {
			mux.Run(val, ctype)
		}
	}
}

// Q adds a new function type into the queue.
func (m *MQue) Q(mx interface{}, rmx ...func()) End {
	if !reflection.IsFuncType(mx) {
		return nil
	}

	tm, _ := reflection.FuncValue(mx)

	var hasArgs bool
	var tu reflect.Type

	args, _ := reflection.GetFuncArgumentsType(mx)
	if size := len(args); size != 0 {
		tu = args[0]
		hasArgs = true
	}

	if !hasArgs {
		index := len(m.any.tms)
		m.any.tms = append(m.any.tms, tm)

		return &mqueSubIndex{
			index:  index,
			queue:  m.any,
			ending: rmx,
		}
	}

	var sub *mqueSub

	m.l.RLock()
	{
		for _, tSub := range m.muxers {
			if tSub.CanRun(tu) {
				sub = tSub
				break
			}
		}
	}
	m.l.RUnlock()

	if sub != nil {
		tindex := len(sub.tms)
		sub.tms = append(sub.tms, tm)

		return &mqueSubIndex{
			index:  tindex,
			queue:  sub,
			ending: rmx,
		}
	}

	var mq mqueSub
	mq.has = true
	mq.am = tu
	mq.tms = []reflect.Value{tm}

	m.l.Lock()
	m.muxers = append(m.muxers, &mq)
	m.l.Unlock()

	return &mqueSubIndex{
		index:  0,
		queue:  &mq,
		ending: rmx,
	}
}

//==============================================================================

type mqueSubIndex struct {
	index  int
	queue  *mqueSub
	ending []func()
}

// AddEnd adds the giving function to the end target.
func (m *mqueSubIndex) AddEnd(fx func()) {
	m.ending = append(m.ending, fx)
}

// End calls removes the listener type from the subscription queue.
func (m *mqueSubIndex) End() {
	if m.queue == nil {
		return
	}

	m.queue.tms = append(m.queue.tms[:m.index], m.queue.tms[m.index+1:]...)
	for _, fx := range m.ending {
		fx()
	}

	m.ending = nil
	m.queue = nil
}

// mqueSub defines a queue subscriber attached to a specific queue.
type mqueSub struct {
	has bool
	am  reflect.Type
	tms []reflect.Value
}

func (m *mqueSub) Flush() {
	m.tms = nil
}

// CanRun returns whether the argument can be used with this subscriber.
func (m *mqueSub) CanRun(d reflect.Type) bool {
	if !d.AssignableTo(m.am) {
		// if d.ConvertibleTo(m.am) {
		// 	return true
		// }

		return false
	}

	return true
}

// Run recevies the argument and
func (m *mqueSub) Run(d interface{}, ctype reflect.Type) {
	if !m.has {
		for _, tm := range m.tms {
			tm.Call([]reflect.Value{})
		}

		return
	}

	var configVal reflect.Value

	if !ctype.AssignableTo(m.am) {
		if !ctype.ConvertibleTo(m.am) {
			return
		}

		vum := reflect.ValueOf(d)
		configVal = vum.Convert(m.am)
	} else {
		configVal = reflect.ValueOf(d)
	}

	for _, tm := range m.tms {
		tm.Call([]reflect.Value{configVal})
	}
}

//==============================================================================
