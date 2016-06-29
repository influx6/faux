// Package pub provides a functional reactive pubsub structure to leverage a
// pure function style reactive behaviour. Originally pulled from pub.Publisher.
package pub

import (
	"errors"
	"io"
	"sync"

	"github.com/satori/go.uuid"
)

// ErrFailedBind represent a failure in binding two Publishers
var ErrFailedBind = errors.New("Failed to Bind Publishers")

// ErrPublisherClosed returned when Publisher is closed
var ErrPublisherClosed = errors.New("Publisher is Closed")

// Handler provides a signal function type:
/*
  It takes three arguments:
		- Publisher:(Publisher) the Publisher itself for reply processing
		- failure:(error) the current error being returned when a data is nil
		- data:(interface{}) the current data being returned,nil when theres an error
*/
type Handler func(Publisher, error, interface{})

// IdiomaticHandler provides a more go-idiomatic handler type which
// depends on the return types of function, using the appropriate
// rule that when the error is nil, a pub.Reply is used else
// a pub.ReplyError. It reduces the typing complexity of Pub.
type IdiomaticHandler func(error, interface{}) (interface{}, error)

// Publisher provides an interface definition for the Publisher type, to allow compatibility by future extenders when composing with other structs.
type Publisher interface {
	io.Closer
	CloseIndicator
	Connector
	Sender
	Replier
	Detacher

	UseRoot(Publisher)
	UUID() string
}

// Pub provides a pure functional Publisher, which uses an internal wait group to ensure if close is called that call values where delivered
type Pub struct {
	op                      Handler
	branches, enders, roots *mapReact
	csignal                 chan bool
	wo                      sync.Mutex
	wg                      sync.WaitGroup
	closed                  bool
	uuid                    string
}

// Pubb returns a new Publisher
func Pubb(fx Handler) Publisher {
	return NewPub(fx)
}

// NewPub returns a new functional Publisher
func NewPub(op Handler) *Pub {
	fr := Pub{
		op:       op,
		branches: newMapReact(),
		enders:   newMapReact(),
		roots:    newMapReact(),
		csignal:  make(chan bool),
		uuid:     uuid.NewV4().String(),
	}

	return &fr
}

// Identity returns Pub that resends its inputs as outputs with no changes.
func Identity() *Pub {
	return NewPub(IdentityMuxer())
}

// Idiomatic returns a Publisher using the IdiomaticMuxer.
func Idiomatic(fx IdiomaticHandler) Publisher {
	return Pubb(IdiomaticMuxer(fx))
}

// Simple returns a Publisher using the SimpleMuxer as a mux generator.
func Simple(fx func(Publisher, interface{})) Publisher {
	return NewPub(SimpleMuxer(fx))
}

// Always returns a Publisher which consistently returns one value
// i.e (the supplied argument).
func Always(v interface{}) Publisher {
	return NewPub(IdentityValueMuxer(v))
}

// UseRoot adds this Publisher as a root Publisher.
func (f *Pub) UseRoot(rx Publisher) {
	f.roots.Add(rx)
}

// LiftOnly calls the Lift function to lift the Publishers and sets the close
// bool to false to prevent closing each other
func LiftOnly(rs ...Publisher) {
	Lift(false, rs...)
}

// LiftOut uses Lift to chain a set of Publishers and returns a new Publisher which is the last in the chain
func LiftOut(conClose bool, rs ...Publisher) Publisher {
	lr := Identity()
	rs = append(rs, lr)
	Lift(conClose, rs...)
	return lr
}

// Lift takes a set of Connectors and pipes the data from one to the next
func Lift(conClose bool, rs ...Publisher) {
	if len(rs) <= 1 {
		return
	}

	var cur Publisher

	cur = rs[0]
	rs = rs[1:]

	for _, co := range rs {
		func(cl Publisher) {
			cur.Bind(cl, conClose)
			cur = cl
		}(co)
	}
}

// Merge merges data from a set of Senders into a new Publisher stream
func Merge(rs ...Publisher) Publisher {
	if len(rs) == 0 {
		return nil
	}

	if len(rs) == 1 {
		return rs[0]
	}

	ms := Identity()
	var endwrap func()

	for _, si := range rs {
		func(so Publisher) {
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

	go func() {
		defer func() { recover(); ms.Close() }()
		endwrap()
	}()

	return ms
}

// Detacher details the detach interface used by the Publisher
type Detacher interface {
	Detach(Publisher)
}

// Detach removes the given Publisher from its connections
func (f *Pub) Detach(rm Publisher) {
	f.enders.Disable(rm)
	f.branches.Disable(rm)
}

//CloseIndicator was created as a later means of providing a simply indicator of the close state of a Publisher
type CloseIndicator interface {
	CloseNotify() <-chan bool
}

// UUID returns a unique id for the publisher.
func (f *Pub) UUID() string {
	return f.uuid
}

// CloseNotify provides a channel for notifying a close event
func (f *Pub) CloseNotify() <-chan bool {
	return f.csignal
}

// Close closes the Publisher and removes all connections
func (f *Pub) Close() error {
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

// Sender defines the delivery methods used to deliver data into Publisher process
type Sender interface {
	Send(v interface{})
	SendError(v error)
}

// Send applies a message value to the handler
func (f *Pub) Send(b interface{}) {
	// if b == nil { return }
	f.wg.Add(1)
	defer f.wg.Done()
	f.op(f, nil, b)
}

// SendError applies a error value to the handler
func (f *Pub) SendError(err error) {
	// if err == nil { return }
	f.wg.Add(1)
	defer f.wg.Done()
	f.op(f, err, nil)
}

// Distribute provide a function that takes a React and other multiple Publishers and distribute the data from the first Publisher to others
func Distribute(rs Publisher, ms ...Sender) Publisher {
	return rs.React(func(r Publisher, err error, data interface{}) {
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
func (f *Pub) Reply(v interface{}) {
	f.branches.Deliver(nil, v)
}

//ReplyError allows the reply of an error message
func (f *Pub) ReplyError(err error) {
	f.branches.Deliver(err, nil)
}

// Connector defines the core connecting methods used for binding with a Publisher
type Connector interface {
	// Bind provides a convenient way of binding 2 Publishers
	Bind(r Publisher, closeAlong bool)
	// React generates a Publisher based off its caller
	React(s Handler, closeAlong bool) Publisher
}

// Bind connects two Publishers
func (f *Pub) Bind(rx Publisher, cl bool) {
	f.branches.Add(rx)
	rx.UseRoot(f)

	if cl {
		f.enders.Add(rx)
	}
}

// React builds a new Publisher from this one
func (f *Pub) React(op Handler, cl bool) Publisher {
	nx := NewPub(op)
	nx.UseRoot(f)

	f.branches.Add(nx)

	if cl {
		f.enders.Add(nx)
	}

	return nx
}

// Stackers provides a construct for providing a strict top-down method call for the Bind,React and BindControl for Publishers,it allows passing these function requests to the last Publisher in the stack while still passing data from the top
type Stackers struct {
	Publisher
	stacks []Connector
	ro     sync.Mutex
}

// ErrEmptyStack is returned when a stack is empty
var ErrEmptyStack = errors.New("Stack Empty")

// PublisherStack returns a Stacker as a Publisher with an identity Publisher as root
func PublisherStack() Publisher {
	return Stack(Identity())
}

// Stack returns a new Publisher based off the Stacker struct which is safe for concurrent use
func Stack(root Publisher) *Stackers {
	sr := Stackers{Publisher: root}
	return &sr
}

// Close wraps the internal close method of the root
func (sr *Stackers) Close() error {
	err := sr.Publisher.Close()
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

// Last returns the last Publishers stacked
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

// Length returns the total stack Publishers
func (sr *Stackers) Length() int {
	var l int
	sr.ro.Lock()
	{
		l = len(sr.stacks)
		sr.ro.Unlock()
	}
	return l
}

// Bind wraps the bind method of the Publisher,if no Publisher has been stack then it binds with the root else gets the last Publisher and binds with that instead
func (sr *Stackers) Bind(r Publisher, cl bool) {
	var lr Connector
	var err error

	if lr, err = sr.Last(); err != nil {
		sr.Publisher.Bind(r, cl)
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

// React wraps the root React() method and stacks the return Publisher or passes it to the last stacked Publisher and stacks that returned Publisher for next use
func (sr *Stackers) React(s Handler, cl bool) Publisher {
	var lr Connector
	var err error

	if lr, err = sr.Last(); err != nil {
		co := sr.Publisher.React(s, cl)
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

// IdentityValueMuxer provides the handler for a providing a pure piping behaviour where data passed in is used as the return data value
func IdentityValueMuxer(v interface{}) Handler {
	return func(r Publisher, err error, _ interface{}) {
		if err != nil {
			r.ReplyError(err)
			return
		}
		r.Reply(v)
	}
}

// IdentityMuxer provides the handoler for a providing a pure piping behaviour where data is left untouched as it comes in and goes out
func IdentityMuxer() Handler {
	return func(r Publisher, err error, data interface{}) {
		if err != nil {
			r.ReplyError(err)
			return
		}
		r.Reply(data)
	}
}

// SimpleMuxer provides the handoler for a providing a pure piping behaviour where data is left untouched as it comes in and goes out
func SimpleMuxer(fx func(Publisher, interface{})) Handler {
	return func(r Publisher, err error, data interface{}) {
		if err != nil {
			r.ReplyError(err)
			return
		}

		fx(r, data)
	}
}

// IdiomaticMuxer returns a Handler that uses a more function return style handler.
func IdiomaticMuxer(h IdiomaticHandler) Handler {
	return func(p Publisher, err error, data interface{}) {
		reply, err2 := h(err, data)
		if err2 != nil {
			p.ReplyError(err2)
			return
		}
		p.Reply(reply)
	}
}

// MapReact provides a nice way of adding multiple reacts into a Publisher reply list.
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
