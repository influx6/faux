// Package pub provides a functional reactive pubsub structure to leverage a
// pure function style reactive behaviour. Originally pulled from pub.Node.
// NOTE: Any use of "asynchronouse" actually means to "run within a goroutine",
// and inversely, the use of "synchronouse" means to run it within the current
// goroutine, generally referred to as "main", or in other words.
package pub

import (
	"reflect"
	"sync"

	"github.com/influx6/faux/context"
	"github.com/influx6/faux/reflection"
	"github.com/satori/go.uuid"
)

var (
	errorType = reflect.TypeOf((*error)(nil)).Elem()
	ctxType   = reflect.TypeOf((*Ctx)(nil)).Elem()
	uType     = reflect.TypeOf((*interface{})(nil)).Elem()

	hlType = reflect.TypeOf((*Handler)(nil)).Elem()
	dlType = reflect.TypeOf((*DataHandler)(nil)).Elem()
	elType = reflect.TypeOf((*ErrorHandler)(nil)).Elem()
)

//==============================================================================

// Ctx defines a type which is passed into all Handlers to provide access
// to an underline context.Context provider and the source Read and Write methods.
type Ctx interface {
	Ctx() context.Context
	RW() ReadWriter
}

// Handler defines a function type which processes data and accepts a ReadWriter
// through which it sends its reply.
type Handler func(Ctx, error, interface{})

// MagicHandler returns a new Handler wrapping the provided value as needed if
// it matches its DataHandler, ErrorHandler, Handler or magic function type.
// MagicFunction type is a function which follows this type form:
// func(Ctx, error, <CustomType>).
func MagicHandler(node interface{}) Handler {
	var hl Handler

	switch node.(type) {
	case func(Ctx, error, interface{}):
		hl = node.(func(Ctx, error, interface{}))
	case func(Ctx, error):
		hl = WrapError(node.(func(Ctx, error)))
	case func(Ctx, interface{}):
		hl = WrapData(node.(func(Ctx, interface{})))
	case func(interface{}) interface{}:
		hl = WrapDataOnly(node.(func(interface{}) interface{}))
	case func(interface{}) error:
		hl = WrapErrorOnly(node.(func(interface{}) error))
	default:
		if !reflection.IsFuncType(node) {
			return nil
		}

		tm, _ := reflection.FuncValue(node)
		args, _ := reflection.GetFuncArgumentsType(node)

		if alen := len(args); alen < 3 || alen > 3 {
			return nil
		}

		// Check if this first item is a Ctx type.
		if ok, _ := reflection.CanSetForType(ctxType, args[0]); !ok {
			return nil
		}

		// Check if this second item is a error type.
		if ok, _ := reflection.CanSetForType(errorType, args[1]); !ok {
			return nil
		}

		data := args[2]
		dZero := reflect.Zero(data)
		dZeroError := reflect.Zero(errorType)

		hl = func(ctx Ctx, err error, val interface{}) {
			ma := reflect.ValueOf(ctx)

			if err != nil {
				tm.Call([]reflect.Value{ma, reflect.ValueOf(err), dZero})
				return
			}

			mVal := reflect.ValueOf(val)
			ok, convert := reflection.CanSetFor(data, mVal)
			if !ok {
				return
			}

			if convert {
				mVal, err = reflection.Convert(data, mVal)
				if err != nil {
					return
				}
			}

			dArgs := []reflect.Value{ma, dZeroError, mVal}
			tm.Call(dArgs)
		}
	}

	return hl
}

// IdentityHandler returns a new Handler which forwards it's errors or data to
// its subscribers.
func IdentityHandler() Handler {
	return func(ctx Ctx, err error, data interface{}) {
		if err != nil {
			ctx.RW().Write(ctx, err)
			return
		}
		ctx.RW().Write(ctx, data)
	}
}

// DataHandler defines a function type that concentrates on handling only data
// replies alone.
type DataHandler func(Ctx, interface{})

// WrapData returns a Handler which wraps a DataHandler within it, but
// passing forward all errors it receives.
func WrapData(dh DataHandler) Handler {
	return func(m Ctx, err error, data interface{}) {
		if err != nil {
			m.RW().Write(m, err)
			return
		}
		dh(m, data)
	}
}

// DataOnlyHandler defines a function type that concentrates on handling only data
// replies alone.
type DataOnlyHandler func(interface{}) interface{}

// WrapDataOnly returns a Handler which wraps a DataHandler within it, but
// passing forward all errors it receives.
func WrapDataOnly(dh DataOnlyHandler) Handler {
	return func(m Ctx, err error, data interface{}) {
		if err != nil {
			m.RW().Write(m, err)
			return
		}

		res := dh(data)

		if ers, ok := res.(error); ok {
			m.RW().Write(m, ers)
			return
		}

		m.RW().Write(m, res)
	}
}

// ErrorHandler defines a function type that concentrates on handling only data
// errors alone.
type ErrorHandler func(Ctx, error)

// WrapError returns a Handler which wraps a DataOnlyHandler within it, but
// passing forward all errors it receives.
func WrapError(dh ErrorHandler) Handler {
	return func(m Ctx, err error, data interface{}) {
		if err == nil {
			dh(m, err)
			return
		}
		m.RW().Write(m, data)
	}
}

// ErrorOnlyHandler defines a function type that concentrates on handling only error
// replies alone.
type ErrorOnlyHandler func(interface{}) error

// WrapErrorOnly returns a Handler which wraps a ErrorOnlyHandler within it, but
// passing forward all errors it receives.
func WrapErrorOnly(dh ErrorOnlyHandler) Handler {
	return func(m Ctx, err error, data interface{}) {
		if err != nil {
			m.RW().Write(m, err)
			return
		}

		m.RW().Write(m, dh(data))
	}
}

//==============================================================================

// ReadWriter defines a type which defines a Reader and Writer interface conforming
// methods.
type ReadWriter interface {
	Reader
	Writer
}

// Node provides an interface definition for the Node type, to allow
// compatibility by future extenders when composing with other structs.
type Node interface {
	ReadWriter
	Reactor
	Inversion

	UUID() string
}

// Magic returns a new functional Node.
func Magic(op interface{}) Node {
	hl := MagicHandler(op)
	if hl == nil {
		panic("invalid type provided")
	}
	return nSync(hl, false)
}

// AsyncMagic returns a new functional Node.
func AsyncMagic(op interface{}) Node {
	hl := MagicHandler(op)
	if hl == nil {
		panic("invalid type provided")
	}
	return aSync(hl, false)
}

// InverseMagic returns a new functional Node.
func InverseMagic(op interface{}) Node {
	hl := MagicHandler(op)
	if hl == nil {
		panic("invalid type provided")
	}
	return nSync(hl, true)
}

// InverseAsyncMagic returns a new functional Node.
func InverseAsyncMagic(op interface{}) Node {
	hl := MagicHandler(op)
	if hl == nil {
		panic("invalid type provided")
	}
	return aSync(hl, true)
}

// pub provides a pure functional Node, which uses an internal wait group to
// ensure if close is called that call values where delivered.
type pub struct {
	uuid string
	op   Handler
	root Node

	async   bool
	inverse bool
	rw      sync.RWMutex
	subs    []Node
}

// nSync returns a new functional Node.
func nSync(op Handler, inverse bool) Node {
	node := pub{
		op:      op,
		inverse: inverse,
		uuid:    uuid.NewV4().String(),
	}

	return &node
}

// aSync returns a new functional Node.
func aSync(op Handler, inverse bool) Node {
	node := pub{
		op:      op,
		async:   true,
		inverse: inverse,
		uuid:    uuid.NewV4().String(),
	}

	return &node
}

// UUID returns the Node unique identification.
func (p *pub) UUID() string {
	return p.uuid
}

// Reader defines the delivery methods used to deliver data into Node process.
type Reader interface {
	Read(v interface{}, ctx ...context.Context)
}

// context defines a struct which composes both a context.Ctx and a
type contxt struct {
	ctx context.Context
	rw  ReadWriter
}

// Ctx returns the context.Context for this struct.
func (c contxt) Ctx() context.Context {
	return c.ctx
}

// RW returns the ReadWriter for this struct.
func (c contxt) RW() ReadWriter {
	return c.rw
}

// Send applies a message value to the handler.
func (p *pub) Read(b interface{}, ctxs ...context.Context) {
	var ctx context.Context

	if len(ctxs) < 1 {
		ctx = context.New()
	} else {
		ctx = ctxs[0]
	}

	ctxn := &contxt{
		ctx: ctx,
		rw:  p,
	}

	if err, ok := b.(error); ok {
		if p.async {
			go p.op(ctxn, err, nil)
			return
		}

		p.op(ctxn, err, nil)
		return
	}

	if p.async {
		go p.op(ctxn, nil, b)
		return
	}

	p.op(ctxn, nil, b)
}

// NthFinder defines a function type which takes the length and index to
// return a new index value.
type NthFinder func(index int, length int) (NewIndex int)

// Writer defines reply methods to reply to requests
type Writer interface {
	Write(Ctx, interface{})
	WriteEvery(Ctx, interface{}, NthFinder)
}

// Write allows the reply of an data message.
// Note: We use the variadic format for the context but only one is used.
func (p *pub) Write(ctx Ctx, v interface{}) {
	ctxn := &contxt{
		ctx: ctx.Ctx(),
		rw:  p,
	}

	var isErr bool

	// Grab the error if it indeed is an error once.
	err, ok := v.(error)
	if ok {
		isErr = true
	}

	p.rw.RLock()
	{
		for _, node := range p.subs {
			if isErr {
				node.Read(err, ctxn.Ctx())
				continue
			}

			node.Read(v, ctxn.Ctx())
		}
	}
	p.rw.RUnlock()

}

func defaultFinder(index int, length int) int {
	return index
}

// WriteEvery allows the delivery/publish of a response to selected index of
// registered nodes using the finder function provided else delivers to all nodes.
// Note: We use the variadic format for the context but only one is used.
func (p *pub) WriteEvery(ctx Ctx, v interface{}, finder NthFinder) {
	ctxn := &contxt{
		ctx: ctx.Ctx(),
		rw:  p,
	}

	if finder == nil {
		finder = defaultFinder
	}

	var isErr bool

	// Grab the error if it indeed is an error once.
	err, ok := v.(error)
	if ok {
		isErr = true
	}

	nlen := len(p.subs)

	p.rw.RLock()
	{
		for index := 0; index < nlen; index++ {
			newIndex := finder(index, nlen)

			if newIndex > 0 && newIndex < nlen {
				node := p.subs[index]

				if isErr {
					node.Write(ctxn, err)
					continue
				}
				node.Write(ctxn, v)
			}

		}
	}
	p.rw.RUnlock()

}

// Inversion defines an interface that allows the creation of an inverter Node
// from another Node, regardless of wether that was inverter or not. Since
// Inversion forcefully forever makes the new Node bind to its last element the
// added element and so on, down the chain, it provides methods suited for ease
// of creation for Nodes.
type Inversion interface {
	Inverse() Node
	InverseWith(interface{}) Node
	InverseAsyncWith(interface{}) Node
}

// Inverse creates a inversed Node with a IdentityHandler which inverts every
// connection you add to it, stacking them serialy down the chain line.
func (p *pub) Inverse() Node {
	node := aSync(IdentityHandler(), true)
	p.Signal(node)
	return node
}

// InverseAsyncWith allows you to create a synchronouse inversed Node.
func (p *pub) InverseWith(node interface{}) Node {
	hl := MagicHandler(node)
	if hl == nil {
		return nil
	}

	return nSync(hl, true)
}

// InverseAsyncWith allows you to create a asynchronouse inversed Node.
func (p *pub) InverseAsyncWith(node interface{}) Node {
	hl := MagicHandler(node)
	if hl == nil {
		return nil
	}

	return aSync(hl, true)
}

// Reactor defines the core connecting methods used for binding with a Node.
type Reactor interface {
	Signal(interface{}) Node
	AsyncSignal(interface{}) Node
}

// Signal sends the response signal from this Node to the provided node.
// If the input is a Node then it is returned, if its a Handler or DataHandler
// then a new Node instance is returned.
func (p *pub) Signal(node interface{}) Node {
	var n Node

	switch node.(type) {
	case Node:
		n = node.(Node)
	default:
		hl := MagicHandler(node)
		if hl == nil {
			return nil
		}
		n = nSync(hl, false)
	}

	// If we have inversed this handler, then return itself, since it
	// wishes to be the first point in entry.
	// To INVERSE, means to redirect how things flows, normally binding
	// flows by returning the next item in the chain, but to invert means
	// to bind to the last item in the subscription list but return itself.
	if p.inverse {
		p.rw.Lock()
		{

			nlen := len(p.subs) - 1
			if nlen > 1 {
				(p.subs[nlen]).Signal(n)
			} else {
				p.subs = append(p.subs, n)
			}

		}
		p.rw.Unlock()

		return p
	}

	p.rw.Lock()
	{
		p.subs = append(p.subs, n)
	}
	p.rw.Unlock()

	return n
}

// AsyncSignal sends the response signal from this Node to the provided node
// within a goroutine. If the input is a Node then it is returned.
func (p *pub) AsyncSignal(node interface{}) Node {
	var n Node

	switch node.(type) {
	case Node:
		n = node.(Node)
	default:
		hl := MagicHandler(node)
		if hl == nil {
			return nil
		}
		n = aSync(hl, false)
	}

	// If we have inversed this handler, then return itself, since it
	// wishes to be the first point in entry.
	// To INVERSE, means to redirect how things flows, normally binding
	// flows by returning the next item in the chain, but to invert means
	// to bind to the last item in the subscription list but return itself.
	if p.inverse {
		p.rw.Lock()
		{

			var newNode Node

			nlen := len(p.subs) - 1
			if nlen > 1 {
				(p.subs[nlen]).AsyncSignal(n)
			} else {
				p.subs = append(p.subs, newNode)
			}

		}
		p.rw.Unlock()

		return p
	}

	p.rw.Lock()
	{
		p.subs = append(p.subs, n)
	}
	p.rw.Unlock()

	return n
}

//==============================================================================

// Lift runs through the giving list of ReadWriters and connects them serialy.
// Chain the next to the previous node.
func Lift(rws ...Reactor) {
	rwsLen := len(rws)

	for index := 0; index < rwsLen; index++ {
		if index < 1 {
			continue
		}

		node := rws[index]
		pnode := rws[index-1]
		pnode.Signal(node)
	}
}

// DeLift runs through the giving list of ReadWriters and connects them
// inversely serialy, chaining the nodes in the inverse order.
func DeLift(rws ...Reactor) {
	rwsLen := len(rws)

	for index := rwsLen - 1; index >= 0; index-- {
		if index >= rwsLen-1 {
			continue
		}

		pnode := rws[index]
		node := rws[index-1]

		pnode.Signal(node)
	}
}
