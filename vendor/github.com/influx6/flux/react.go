package flux

import (
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"sync"
	"sync/atomic"
)

// ErrFailedBind represent a failure in binding two Reactors
var ErrFailedBind = errors.New("Failed to Bind Reactors")

// ErrReactorClosed returned when reactor is closed
var ErrReactorClosed = errors.New("Reactor is Closed")

// SignalMuxHandler provides a signal function type:
/*
  It takes three arguments:
		- reactor:(Reactor) the reactor itself for reply processing
		- failure:(error) the current error being returned when a data is nil
		- data:(interface{}) the current data being returned,nil when theres an error
*/
type SignalMuxHandler func(reactor Reactor, failure error, signal interface{})

// Reactor provides an interface definition for the reactor type to allow compatibility by future extenders when composing with other structs.
type Reactor interface {
	io.Closer
	CloseIndicator
	Connector
	Sender
	Replier
	Detacher

	UseRoot(Reactor)
}

// FlatReactor provides a pure functional reactor which uses an internal wait group to ensure if close is called that call values where delivered
type FlatReactor struct {
	op SignalMuxHandler
	// root,next Reactor
	branches, enders, roots *mapReact
	csignal                 chan bool
	wo                      sync.Mutex
	wg                      sync.WaitGroup
	closed                  bool
}

// Reactive returns a ReactiveStacks
func Reactive(fx SignalMuxHandler) Reactor {
	return FlatReactive(fx)
}

// FlatReactive returns a new functional reactor
func FlatReactive(op SignalMuxHandler) *FlatReactor {
	fr := FlatReactor{
		op:       op,
		branches: newMapReact(),
		enders:   newMapReact(),
		roots:    newMapReact(),
		csignal:  make(chan bool),
	}

	return &fr
}

// FlatStack returns a flat reactor
func FlatStack() Reactor {
	return ReactStack(ReactIdentity())
}

// IdentityValueMuxer provides the handler for a providing a pure piping behaviour where data passed in is used as the return data value
func IdentityValueMuxer(v interface{}) SignalMuxHandler {
	return func(r Reactor, err error, _ interface{}) {
		if err != nil {
			r.ReplyError(err)
			return
		}
		r.Reply(v)
	}
}

// IdentityMuxer provides the handoler for a providing a pure piping behaviour where data is left untouched as it comes in and goes out
func IdentityMuxer() SignalMuxHandler {
	return func(r Reactor, err error, data interface{}) {
		if err != nil {
			r.ReplyError(err)
			return
		}
		r.Reply(data)
	}
}

// SimpleMuxer provides the handoler for a providing a pure piping behaviour where data is left untouched as it comes in and goes out
func SimpleMuxer(fx func(Reactor, interface{})) SignalMuxHandler {
	return func(r Reactor, err error, data interface{}) {
		if err != nil {
			r.ReplyError(err)
			return
		}

		fx(r, data)
	}
}

// FlatIdentity returns flatreactor that resends its inputs as outputs with no changes
func FlatIdentity() *FlatReactor {
	return FlatReactive(IdentityMuxer())
}

// FlatSimple returns a reactor using the SimpleMuxer as a mux generator
func FlatSimple(fx func(Reactor, interface{})) Reactor {
	return Reactive(SimpleMuxer(fx))
}

// FlatAlways returns a reactor with consistently returns the provided value
func FlatAlways(v interface{}) Reactor {
	return FlatReactive(IdentityValueMuxer(v))
}

// ReactIdentity is more written to provide a backward compatibility for cold using the old
//channel based reactor
func ReactIdentity() Reactor {
	return FlatIdentity()
}

// UseRoot Adds this reactor as a root of the called reactor
func (f *FlatReactor) UseRoot(rx Reactor) {
	f.roots.Add(rx)
}

// LiftOnly calls the Lift function to lift the Reactors and sets the close
// bool to false to prevent closing each other
func LiftOnly(rs ...Reactor) {
	Lift(false, rs...)
}

// LiftOut uses Lift to chain a set of reactors and returns a new reactor which is the last in the chain
func LiftOut(conClose bool, rs ...Reactor) Reactor {
	lr := FlatIdentity()
	rs = append(rs, lr)
	Lift(conClose, rs...)
	return lr
}

// Lift takes a set of Connectors and pipes the data from one to the next
func Lift(conClose bool, rs ...Reactor) {
	if len(rs) <= 1 {
		return
	}

	var cur Reactor

	cur = rs[0]
	rs = rs[1:]

	for _, co := range rs {
		func(cl Reactor) {
			cur.Bind(cl, conClose)
			cur = cl
		}(co)
	}
}

// MergeReactors merges data from a set of Senders into a new reactor stream
func MergeReactors(rs ...Reactor) Reactor {
	if len(rs) == 0 {
		return nil
	}

	if len(rs) == 1 {
		return rs[0]
	}

	ms := ReactIdentity()
	var endwrap func()

	for _, si := range rs {
		func(so Reactor) {
			so.Bind(ms, false)
			if endwrap != nil {
				oc := endwrap
				endwrap = func() {
					<-so.CloseNotify()
					oc()
				}
			} else {
				endwrap = func() {
					<-so.CloseNotify()
				}
			}
		}(si)
	}

	GoSilent("MergeClose", func() {
		defer ms.Close()
		endwrap()
	})

	return ms
}

// Manage manages the operations of reactor
// func (r *ReactiveStack) Manage() {
// 	defer func() {
// 		r.enders.Close()
// 		r.roots.Clean()
// 		r.branch.Clean()
// 	}()
//
// 	for {
// 		select {
// 		case err, ok := <-r.ps.Errors:
// 			if !ok {
// 				return
// 			}
// 			r.wg.Done()
// 			// go r.op(r, err.(error), nil)
// 			r.op(r, err.(error), nil)
// 		case signal, ok := <-r.ps.Signals:
// 			if !ok {
// 				return
// 			}
// 			r.wg.Done()
// 			// go r.op(r, nil, signal)
// 			r.op(r, nil, signal)
// 		}
// 	}
// }

//CloseIndicator was created as a later means of providing a simply indicator of the close state of a Reactor
type CloseIndicator interface {
	CloseNotify() <-chan bool
}

// CloseNotify provides a channel for notifying a close event
func (f *FlatReactor) CloseNotify() <-chan bool {
	return f.csignal
}

// Close closes the reactor and removes all connections
func (f *FlatReactor) Close() error {
	f.wo.Lock()
	ok := f.closed
	f.wo.Unlock()

	if ok {
		return nil
	}

	f.wg.Wait()

	f.wo.Lock()
	f.closed = true
	f.wo.Unlock()

	f.branches.Close()

	f.roots.Do(func(rm SenderDetachCloser) {
		go rm.Detach(f)
	})

	f.roots.Close()
	f.enders.Close()

	close(f.csignal)
	return nil
}

// Detacher details the detach interface used by the Reactor
type Detacher interface {
	Detach(Reactor)
}

// Detach removes the given reactor from its connections
func (f *FlatReactor) Detach(rm Reactor) {
	f.enders.Disable(rm)
	f.branches.Disable(rm)
}

// Sender defines the delivery methods used to deliver data into Reactor process
type Sender interface {
	Send(v interface{})
	SendError(v error)
}

// Send applies a message value to the handler
func (f *FlatReactor) Send(b interface{}) {
	// if b == nil { return }
	f.wg.Add(1)
	defer f.wg.Done()
	f.op(f, nil, b)
}

// SendError applies a error value to the handler
func (f *FlatReactor) SendError(err error) {
	// if err == nil { return }
	f.wg.Add(1)
	defer f.wg.Done()
	f.op(f, err, nil)
}

// DistributeSignals provide a function that takes a React and other multiple Reactors and distribute the data from the first reactor to others
func DistributeSignals(rs Reactor, ms ...Sender) Reactor {
	return rs.React(func(r Reactor, err error, data interface{}) {
		for _, mi := range ms {
			go func(rec Sender) {
				if err != nil {
					rec.SendError(err)
					return
				}
				rec.Send(data)
			}(mi)
		}
	}, true)
}

// Replier defines reply methods to reply to requests
type Replier interface {
	Reply(v interface{})
	ReplyError(v error)
}

//Reply allows the reply of an data message
func (f *FlatReactor) Reply(v interface{}) {
	f.branches.Deliver(nil, v)
}

//ReplyError allows the reply of an error message
func (f *FlatReactor) ReplyError(err error) {
	f.branches.Deliver(err, nil)
}

// Connector defines the core connecting methods used for binding with a Reactor
type Connector interface {
	// Bind provides a convenient way of binding 2 reactors
	Bind(r Reactor, closeAlong bool)
	// React generates a reactor based off its caller
	React(s SignalMuxHandler, closeAlong bool) Reactor
}

// Bind connects two reactors
func (f *FlatReactor) Bind(rx Reactor, cl bool) {
	f.branches.Add(rx)
	rx.UseRoot(f)

	if cl {
		f.enders.Add(rx)
	}
}

// React builds a new reactor from this one
func (f *FlatReactor) React(op SignalMuxHandler, cl bool) Reactor {
	nx := FlatReactive(op)
	nx.UseRoot(f)

	f.branches.Add(nx)

	if cl {
		f.enders.Add(nx)
	}

	return nx
}

// Stackers provides a construct for providing a strict top-down method call for the Bind,React and BindControl for Reactors,it allows passing these function requests to the last Reactor in the stack while still passing data from the top
type Stackers struct {
	Reactor
	stacks []Connector
	ro     sync.Mutex
}

// ErrEmptyStack is returned when a stack is empty
var ErrEmptyStack = errors.New("Stack Empty")

// ReactorStack returns a Stacker as a reactor with an identity reactor as root
func ReactorStack() Reactor {
	return ReactStack(ReactIdentity())
}

// ReactStack returns a new Reactor based off the Stacker struct which is safe for concurrent use
func ReactStack(root Reactor) *Stackers {
	sr := Stackers{Reactor: root}
	return &sr
}

// Close wraps the internal close method of the root
func (sr *Stackers) Close() error {
	err := sr.Reactor.Close()
	sr.Clear()
	return err
}

// Clear clears the stacks and resolves back to root
func (sr *Stackers) Clear() {
	sr.ro.Lock()
	{
		sr.stacks = nil
	}
	sr.ro.Unlock()
}

// Last returns the last Reactors stacked
func (sr *Stackers) Last() (Connector, error) {
	var r Connector
	sr.ro.Lock()
	{
		l := len(sr.stacks)
		if l > 0 {
			r = sr.stacks[l-1]
		}
	}
	sr.ro.Unlock()

	if r == nil {
		return nil, ErrEmptyStack
	}

	return r, nil
}

// Length returns the total stack Reactors
func (sr *Stackers) Length() int {
	var l int
	sr.ro.Lock()
	{
		l = len(sr.stacks)
		sr.ro.Unlock()
	}
	return l
}

// Bind wraps the bind method of the Reactor,if no Reactor has been stack then it binds with the root else gets the last Reactor and binds with that instead
func (sr *Stackers) Bind(r Reactor, cl bool) {
	var lr Connector
	var err error

	if lr, err = sr.Last(); err != nil {
		sr.Reactor.Bind(r, cl)
		sr.ro.Lock()
		{
			sr.stacks = append(sr.stacks, r)
		}
		sr.ro.Unlock()
		return
	}

	lr.Bind(r, cl)
	sr.ro.Lock()
	{
		sr.stacks = append(sr.stacks, r)
	}
	sr.ro.Unlock()
}

// React wraps the root React() method and stacks the return Reactor or passes it to the last stacked Reactor and stacks that returned reactor for next use
func (sr *Stackers) React(s SignalMuxHandler, cl bool) Reactor {
	var lr Connector
	var err error

	if lr, err = sr.Last(); err != nil {
		co := sr.Reactor.React(s, cl)
		sr.ro.Lock()
		{
			sr.stacks = append(sr.stacks, co)
		}
		sr.ro.Unlock()
		return co
	}

	co := lr.React(s, cl)
	sr.ro.Lock()
	{
		sr.stacks = append(sr.stacks, co)
	}
	sr.ro.Unlock()
	return co
}

// SendBinder defines the combination of the Sender and Binding interfaces
type SendBinder interface {
	Sender
	Connector
}

// SendReplier provides the interface for the combination of senders and repliers
type SendReplier interface {
	Replier
	Sender
}

// SendReplyCloser provides the interface for the combination of closers,senders and repliers
type SendReplyCloser interface {
	io.Closer
	Replier
	Sender
}

// SendReplyDetacher provides the interface for the combination of senders,detachers and repliers
type SendReplyDetacher interface {
	Replier
	Sender
	Detacher
}

// SendReplyDetachCloser provides the interface for the combination of closers, senders,detachers and repliers
type SendReplyDetachCloser interface {
	io.Closer
	Replier
	Sender
	Detacher
}

// ReplyCloser provides an interface that combines Replier and Closer interfaces
type ReplyCloser interface {
	Replier
	io.Closer
}

// SendCloser provides an interface that combines Sender and Closer interfaces
type SendCloser interface {
	Sender
	io.Closer
}

// SenderDetachCloser provides an interface that combines Sender and Closer interfaces
type SenderDetachCloser interface {
	Sender
	Detacher
	io.Closer
}

// MapReact provides a nice way of adding multiple reacts into a reactor reply list.
type mapReact struct {
	ro sync.RWMutex
	ma map[SenderDetachCloser]bool
}

// NewMapReact returns a new MapReact store
func newMapReact() *mapReact {
	ma := mapReact{ma: make(map[SenderDetachCloser]bool)}
	return &ma
}

// Clean resets the map
func (m *mapReact) Clean() {
	m.ro.Lock()
	m.ma = make(map[SenderDetachCloser]bool)
	m.ro.Unlock()
}

// Deliver either a data or error to the Sender
func (m *mapReact) Deliver(err error, data interface{}) {
	m.ro.RLock()
	for ms, ok := range m.ma {
		if !ok {
			continue
		}

		if err != nil {
			ms.SendError(err)
			continue
		}

		ms.Send(data)
	}
	m.ro.RUnlock()
}

// Add a sender into the map as available
func (m *mapReact) Add(r SenderDetachCloser) {
	m.ro.Lock()
	if !m.ma[r] {
		m.ma[r] = true
	}
	m.ro.Unlock()
}

// Disable a particular sender
func (m *mapReact) Disable(r SenderDetachCloser) {
	var ok bool
	m.ro.RLock()
	_, ok = m.ma[r]
	m.ro.RUnlock()

	if !ok {
		return
	}

	m.ro.Lock()
	m.ma[r] = false
	m.ro.Unlock()
}

// Length returns the length of the map
func (m *mapReact) Length() int {
	var l int
	m.ro.RLock()
	l = len(m.ma)
	m.ro.RUnlock()
	return l
}

// Do performs the function on every item
func (m *mapReact) Do(fx func(SenderDetachCloser)) {
	m.ro.RLock()
	for ms := range m.ma {
		fx(ms)
	}
	m.ro.RUnlock()
}

// DisableAll disables the items in the map
func (m *mapReact) DisableAll() {
	m.ro.Lock()
	for ms := range m.ma {
		m.ma[ms] = false
	}
	m.ro.Unlock()
}

// Close closes all the SenderDetachClosers
func (m *mapReact) Close() {
	m.ro.RLock()
	for ms, ok := range m.ma {
		if !ok {
			continue
		}
		ms.Close()
	}
	m.ro.RUnlock()
	m.Clean()
}

// ChannelStream provides a simple struct for exposing outputs from Reactor to outside
type ChannelStream struct {
	Data   chan interface{}
	Error  chan error
	closed int64
}

// NewChannelStream returns a new channel stream instance with blocked channels, so ensure to fullfill the contract of removing the data you need only
func NewChannelStream() *ChannelStream {
	cs := ChannelStream{
		Data:  make(chan interface{}),
		Error: make(chan error),
	}
	return &cs
}

func (c *ChannelStream) error(err error) {
	if atomic.LoadInt64(&c.closed) <= 0 {
		go func() { c.Error <- err }()
		// c.Error <- err
	}
}

func (c *ChannelStream) data(d interface{}) {
	if atomic.LoadInt64(&c.closed) <= 0 {
		go func() { c.Data <- d }()
		// c.Data <- d
	}
}

// Close ends the capability to use the ChannelStream channels
func (c *ChannelStream) Close() {
	atomic.StoreInt64(&c.closed, 1)
}

// Listen binds into a reactor and will pipe any response into its Data or Error channels, always use this to bind to reactors, to ensure safety in code use i.e dont try to pipe into the channels your own way
func (c *ChannelStream) Listen(m Reactor) {
	if atomic.LoadInt64(&c.closed) > 0 {
		return
	}

	m.React(func(m Reactor, err error, data interface{}) {
		if err != nil {
			c.error(err)
			m.ReplyError(err)
			return
		}
		c.data(data)
		m.Reply(err)
	}, true)
}

// JSONReactor provides a json encoding Reactor,takes any input and tries to transform it into a json using the default json.Marshal function
func JSONReactor() Reactor {
	return Reactive(func(r Reactor, err error, d interface{}) {
		if err != nil {
			r.ReplyError(err)
			return
		}

		jsonbu, err := json.Marshal(d)

		if err != nil {
			r.ReplyError(err)
			return
		}

		r.Reply(jsonbu)
	})
}

// ErrValueType is returned by FileLoader if it gets a type that is not a string
var ErrValueType = errors.New("Value is not a string type")

// FileLoader provides an adaptor to load a file path
func FileLoader() Reactor {
	return Reactive(func(v Reactor, err error, d interface{}) {
		if err != nil {
			v.ReplyError(err)
			return
		}

		var file string
		var ok bool

		if file, ok = d.(string); !ok {
			v.ReplyError(ErrValueType)
			return
		}

		var data []byte

		if data, err = ioutil.ReadFile(file); err != nil {
			v.ReplyError(err)
			return
		}

		v.Reply(string(data))
	})
}

// QueueReactor provides a reactor that listens on the supplied
// queue for data and error messages but if the queue gets closed then
//the reactor is closed along
func QueueReactor(ps *PressureStream) (qr Reactor) {
	qr = ReactIdentity()
	GoDefer("ReactorQueue:Manage", func() {
		defer qr.Close()
		for {
			select {
			case d, ok := <-ps.Signals:
				if !ok {
					return
				}
				qr.Send(d)
			case e, ok := <-ps.Errors:
				if !ok {
					return
				}
				qr.SendError(e.(error))
			}
		}
	})
	return
}
