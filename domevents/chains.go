package domevents

import "sync"

// ChainHandler provides a handler for event chain types.
type ChainHandler func(Event, EventHandler)

//Chains define a simple  chain
type Chains interface {
	HandleContext(Event)
	Next(ChainHandler) Chains
	Chain(Chains) Chains
	NChain(Chains) Chains
	Bind(EventHandler) Chains
	UseNext(Chains)
	UsePrev(Chains)
	UnChain()
}

// Chain provides a simple middleware like
type Chain struct {
	op         ChainHandler
	prev, next Chains
}

//ChainIdentity returns a chain that calls the next automatically
func ChainIdentity() Chains {
	return NewChain(func(c Event, nx EventHandler) {
		nx(c)
	})
}

//NewChain returns a new Chain instance
func NewChain(fx ChainHandler) *Chain {
	return &Chain{
		op: fx,
	}
}

// UnChain unlinks the current chain from the set and reconnects the others
func (r *Chain) UnChain() {
	prev := r.prev
	next := r.next

	if prev != nil && next != nil {
		prev.UseNext(next)
		next.UsePrev(prev)
		return
	}

	prev.UseNext(nil)
}

// Bind provides a wrapper over the Next binder function call
func (r *Chain) Bind(rnx EventHandler) Chains {
	return r.Next(func(ev Event, nx EventHandler) {
		rnx(ev)
		nx(ev)
	})
}

// Next allows the chain of the function as a Handler
func (r *Chain) Next(rnx ChainHandler) Chains {
	nx := NewChain(rnx)
	return r.NChain(nx)
}

// Chain sets the next  chains else passes it down to the last chain to set as next chain,returning itself
func (r *Chain) Chain(rx Chains) Chains {
	if r.next == nil {
		rx.UsePrev(r)
		r.UseNext(rx)
	} else {
		r.next.Chain(rx)
	}
	return r
}

// NChain sets the next  chains else passes it down to the last chain to set as next chain,returning the the supplied chain
func (r *Chain) NChain(rx Chains) Chains {
	if r.next == nil {
		r.UseNext(rx)
		return rx
	}

	return r.next.NChain(rx)
}

// HandleContext calls the next chain if any
func (r *Chain) HandleContext(c Event) {
	r.op(c, func(c Event) {
		if r.next != nil {
			r.next.HandleContext(c)
		}
	})
}

// UseNext swaps the next chain with the supplied
func (r *Chain) UseNext(fl Chains) {
	r.next = fl
}

// UsePrev swaps the previous chain with the supplied
func (r *Chain) UsePrev(fl Chains) {
	r.prev = fl
}

//Connect chains second set to the first Chain and returns the first Chain
func Connect(mo Chains, so ...Chains) Chains {
	for _, sof := range so {
		func(do Chains) {
			mo.Chain(do)
		}(sof)
	}
	return mo
}

//ElemEventMux represents a stanard callback function for dom events
type ElemEventMux func(Event, func())

//ListenerStack provides addition of functions into a stack
type ListenerStack struct {
	listeners []ElemEventMux
	lock      sync.RWMutex
}

//NewListenerStack returns a new ListenerStack instance
func NewListenerStack() *ListenerStack {
	return &ListenerStack{
		listeners: make([]ElemEventMux, 0),
	}
}

//Clear flushes the stack listener
func (f *ListenerStack) Clear() {
	f.lock.Lock()
	f.listeners = f.listeners[0:0]
	f.lock.Unlock()
}

//Size returns the total number of listeners
func (f *ListenerStack) Size() int {
	f.lock.RLock()
	sz := len(f.listeners)
	f.lock.RUnlock()
	return sz
}

//Add adds a function into the stack
func (f *ListenerStack) Add(fx ElemEventMux) int {
	var ind int

	f.lock.RLock()
	ind = len(f.listeners)
	f.lock.RUnlock()

	f.lock.Lock()
	f.listeners = append(f.listeners, fx)
	f.lock.Unlock()

	return ind
}

// DeleteIndex removes the function at the provided index
func (f *ListenerStack) DeleteIndex(ind int) {

	if ind <= 0 && len(f.listeners) <= 0 {
		return
	}

	f.lock.Lock()
	copy(f.listeners[ind:], f.listeners[ind+1:])
	f.lock.Unlock()

	f.lock.RLock()
	f.listeners[len(f.listeners)-1] = nil
	f.lock.RUnlock()

	f.lock.Lock()
	f.listeners = f.listeners[:len(f.listeners)-1]
	f.lock.Unlock()
}

//Each runs through the function lists and executing with args
func (f *ListenerStack) Each(d Event) {
	if f.Size() <= 0 {
		return
	}

	f.lock.RLock()

	// var ro sync.Mutex
	var stop bool

	for _, fx := range f.listeners {
		if stop {
			break
		}
		//TODO: is this critical that we send it into a goroutine with a mutex?
		fx(d, func() {
			// ro.Lock()
			stop = true
			// ro.Unlock()
		})
	}

	f.lock.RUnlock()
}
