package flux

//FlatChains define a simple flat chain
type FlatChains interface {
	Call(error, interface{})
	Next(FlatChains)
}

// NextHandler provides next call for flat chains
type NextHandler func(error, interface{})

// FlatHandler provides a handler for flatchain
type FlatHandler func(error, interface{}, NextHandler)

// FlatChain provides a simple middleware like
type FlatChain struct {
	op   FlatHandler
	next FlatChains
}

//NewFlatChain returns a new flatchain instance
func NewFlatChain(fx FlatHandler) *FlatChain {
	return &FlatChain{
		op: fx,
	}
}

// Next sets the next flat chains else passes it down to the last chain to set as next chain
func (r *FlatChain) Next(rx FlatChains) {
	if r.next == nil {
		r.next = rx
		return
	}
	r.next.Next(rx)
}

// Call calls the next chain if any
func (r *FlatChain) Call(err error, d interface{}) {
	r.op(err, d, func(ex error, v interface{}) {
		if r.next != nil {
			r.next.Call(err, d)
		}
	})
}
