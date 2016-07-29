// Package pub provides a functional reactive pubsub structure to leverage a
// pure function style reactive behaviour. Originally pulled from pub.Node.
// NOTE: Any use of "asynchronouse" actually means to "run within a goroutine",
// and inversely, the use of "synchronouse" means to run it within the current
// goroutine, generally referred to as "main", or in other words.
package pub

import (
	"errors"
	"reflect"
	"sync"
	"sync/atomic"

	"github.com/influx6/faux/context"
	"github.com/influx6/faux/reflection"
)

var (
	errorType = reflect.TypeOf((*error)(nil)).Elem()
	ctxType   = reflect.TypeOf((*Ctx)(nil)).Elem()
	uType     = reflect.TypeOf((*interface{})(nil)).Elem()

	hlType = reflect.TypeOf((*Handler)(nil)).Elem()
	dlType = reflect.TypeOf((*DataHandler)(nil)).Elem()
	elType = reflect.TypeOf((*ErrorHandler)(nil)).Elem()

	dZeroError = reflect.Zero(errorType)
)

//==============================================================================

// Ctx defines a context interface with the ability to signal ending of a
// stream to other recievers.
type Ctx interface {
	context.Context
	End()
	Ended() bool
	Renew() Ctx
}

// NewCtx returns an instance of a new Ctx.
func NewCtx() Ctx {
	return &ctxn{
		Context: context.New(),
	}
}

// NewCtxWith returns an instance of a new Ctx using the provided context.Context.
func NewCtxWith(ctx context.Context) Ctx {
	return &ctxn{
		Context: ctx,
	}
}

type ctxn struct {
	context.Context
	streamEnded int64
}

// Renew returns a new context with the underline context.Context for this
// context.
func (c *ctxn) Renew() Ctx {
	return &ctxn{Context: c.Context}
}

// End sets the context as done. A context can only end once,and provides a
// means to signaify the end of a stream.
func (c *ctxn) End() {
	atomic.StoreInt64(&c.streamEnded, 1)
}

// Ended returns true/false if the current context has ended.
func (c *ctxn) Ended() bool {
	return atomic.LoadInt64(&c.streamEnded) != 0
}

//==============================================================================

// Handler defines a function type which processes data and accepts a ReadWriter
// through which it sends its reply.
type Handler func(Ctx, error, interface{}) (interface{}, error)

// MustWrap returns the Handler else panics if it fails to create the Handler
// from the provided function type.
func MustWrap(node interface{}) Handler {
	dh := Wrap(node)
	if dh != nil {
		return dh
	}

	panic("Invalid type provided for Handler")
}

// Wrap returns a new Handler wrapping the provided value as needed if
// it matches its DataHandler, ErrorHandler, Handler or magic function type.
// MagicFunction type is a function which follows this type form:
// func(context.Context, error, <CustomType>).
func Wrap(node interface{}) Handler {
	var hl Handler

	switch node.(type) {
	case func(Ctx, error, interface{}) (interface{}, error):
		hl = node.(func(Ctx, error, interface{}) (interface{}, error))
	case func(Ctx, interface{}):
		hl = wrapDataWithNoReturn(node.(func(Ctx, interface{})))
	case func(Ctx, interface{}) interface{}:
		hl = wrapDataWithReturn(node.(func(Ctx, interface{}) interface{}))
	case func(Ctx, interface{}) (interface{}, error):
		hl = wrapData(node.(func(Ctx, interface{}) (interface{}, error)))
	case func(Ctx, error) (interface{}, error):
		hl = wrapError(node.(func(Ctx, error) (interface{}, error)))
	case func(interface{}) interface{}:
		hl = wrapDataOnly(node.(func(interface{}) interface{}))
	case func(interface{}):
		hl = wrapJustData(node.(func(interface{})))
	case func(error):
		hl = wrapJustError(node.(func(error)))
	case func(error) error:
		hl = wrapErrorReturn(node.(func(error) error))
	case func() interface{}:
		hl = wrapNoData(node.(func() interface{}))
	case func(interface{}) error:
		hl = wrapErrorOnly(node.(func(interface{}) error))
	default:
		if !reflection.IsFuncType(node) {
			return nil
		}

		tm, _ := reflection.FuncValue(node)
		args, _ := reflection.GetFuncArgumentsType(node)

		dLen := len(args)

		if dLen < 2 {
			return nil
		}

		// Check if this first item is a Ctx type.
		if ok, _ := reflection.CanSetForType(ctxType, args[0]); !ok {
			return nil
		}

		var data reflect.Type
		var isCustorm bool

		if dLen > 2 {

			// Check if this second item is a error type.
			if ok, _ := reflection.CanSetForType(errorType, args[1]); !ok {
				return nil
			}

			data = args[2]
		} else {
			data = args[1]
			isCustorm = true
		}

		dZero := reflect.Zero(data)

		hl = func(ctx Ctx, err error, val interface{}) (interface{}, error) {
			ma := reflect.ValueOf(ctx)
			me := dZeroError

			if err != nil {
				me = reflect.ValueOf(err)

				if !isCustorm {
					resArgs := tm.Call([]reflect.Value{ma, me, dZero})
					if len(resArgs) < 1 {
						return nil, nil
					}

					if len(resArgs) == 1 {
						rVal := resArgs[0].Interface()
						if dx, ok := rVal.(error); ok {
							return nil, dx
						}

						return rVal, nil
					}

					mr1 := resArgs[0].Interface()
					mr2 := resArgs[1].Interface()

					if emr2, ok := mr2.(error); ok {
						return mr1, emr2
					}

					return mr1, nil
				}

				return nil, err
			}

			mVal := dZero

			if val != nil {
				mVal = reflect.ValueOf(val)

				ok, convert := reflection.CanSetFor(data, mVal)
				if !ok {
					return nil, errors.New("Invalid Type Received")
				}

				if convert {
					mVal, err = reflection.Convert(data, mVal)
					if err != nil {
						return nil, errors.New("Type Conversion Failed")
					}
				}

			}

			if !isCustorm {
				dArgs := []reflect.Value{ma, me, mVal}
				resArgs := tm.Call(dArgs)
				if len(resArgs) < 1 {
					return nil, nil
				}

				if len(resArgs) == 1 {
					rVal := resArgs[0].Interface()
					if dx, ok := rVal.(error); ok {
						return nil, dx
					}

					return rVal, nil
				}

				mr1 := resArgs[0].Interface()
				mr2 := resArgs[1].Interface()

				if emr2, ok := mr2.(error); ok {
					return mr1, emr2
				}

				return mr1, nil
			}

			dArgs := []reflect.Value{ma, mVal}

			resArgs := tm.Call(dArgs)
			if len(resArgs) < 1 {
				return nil, nil
			}

			if len(resArgs) == 1 {
				rVal := resArgs[0].Interface()
				if dx, ok := rVal.(error); ok {
					return nil, dx
				}

				return rVal, nil
			}

			mr1 := resArgs[0].Interface()
			mr2 := resArgs[1].Interface()

			if emr2, ok := mr2.(error); ok {
				return mr1, emr2
			}

			return mr1, nil
		}
	}

	return hl
}

// IdentityHandler returns a new Handler which forwards it's errors or data to
// its subscribers.
func IdentityHandler() Handler {
	return func(ctx Ctx, err error, data interface{}) (interface{}, error) {
		if err != nil {
			return nil, err
		}
		return data, nil
	}
}

// DataHandler defines a function type that concentrates on handling only data
// replies alone.
type DataHandler func(Ctx, interface{}) (interface{}, error)

// wrapData returns a Handler which wraps a DataHandler within it, but
// passing forward all errors it receives.
func wrapData(dh DataHandler) Handler {
	return func(m Ctx, err error, data interface{}) (interface{}, error) {
		if err != nil {
			return nil, err
		}

		return dh(m, data)
	}
}

// DataWithNoReturnHandler defines a function type that concentrates on handling only data
// replies alone.
type DataWithNoReturnHandler func(Ctx, interface{})

// wrapDataWithNoReturn returns a Handler which wraps a DataHandler within it, but
// passing forward all errors it receives.
func wrapDataWithNoReturn(dh DataWithNoReturnHandler) Handler {
	return func(m Ctx, err error, data interface{}) (interface{}, error) {
		if err != nil {
			return nil, err
		}

		dh(m, data)
		return data, nil
	}
}

// DataWithReturnHandler defines a function type that concentrates on handling only data
// replies alone.
type DataWithReturnHandler func(Ctx, interface{}) interface{}

// wrapDataWithReturn returns a Handler which wraps a DataHandler within it, but
// passing forward all errors it receives.
func wrapDataWithReturn(dh DataWithReturnHandler) Handler {
	return func(m Ctx, err error, data interface{}) (interface{}, error) {
		if err != nil {
			return nil, err
		}

		return dh(m, data), nil
	}
}

// NoDataHandler defines an Handler which allows a return value when called
// but has no data passed in.
type NoDataHandler func() interface{}

// wrapNoData returns a Handler which wraps a NoDataHandler within it, but
// forwards all errors it receives. It calls its internal function
// with no arguments taking the response and sending that out.
func wrapNoData(dh NoDataHandler) Handler {
	return func(m Ctx, err error, data interface{}) (interface{}, error) {
		if err != nil {
			return nil, err
		}

		res := dh()
		if erx, ok := res.(error); ok {
			return nil, erx
		}

		return res, nil
	}
}

// DataOnlyHandler defines a function type that concentrates on handling only data
// replies alone.
type DataOnlyHandler func(interface{}) interface{}

// wrapDataOnly returns a Handler which wraps a DataOnlyHandler within it, but
// passing forward all errors it receives.
func wrapDataOnly(dh DataOnlyHandler) Handler {
	return func(m Ctx, err error, data interface{}) (interface{}, error) {
		if err != nil {
			return nil, err
		}

		res := dh(data)
		if erx, ok := res.(error); ok {
			return nil, erx
		}

		return res, nil
	}
}

// JustDataHandler defines a function type which expects one argument.
type JustDataHandler func(interface{})

// wrapJustData wraps a JustDataHandler and returns it as a Handler.
func wrapJustData(dh JustDataHandler) Handler {
	return func(ctx Ctx, err error, d interface{}) (interface{}, error) {
		if err != nil {
			return nil, err
		}

		dh(d)
		return d, nil
	}
}

// JustErrorHandler defines a function type that concentrates on handling only
// errors alone.
type JustErrorHandler func(error)

// wrapJustError returns a Handler which wraps a DataOnlyHandler within it, but
// passing forward all errors it receives.
func wrapJustError(dh JustErrorHandler) Handler {
	return func(ctx Ctx, err error, d interface{}) (interface{}, error) {
		if err != nil {
			dh(err)
			return nil, err
		}

		return d, nil
	}
}

// ErrorReturnHandler defines a function type that concentrates on handling only data
// errors alone.
type ErrorReturnHandler func(error) error

// wrapErrorReturn returns a Handler which wraps a DataOnlyHandler within it, but
// passing forward all errors it receives.
func wrapErrorReturn(dh ErrorReturnHandler) Handler {
	return func(ctx Ctx, err error, d interface{}) (interface{}, error) {
		if err != nil {
			return nil, dh(err)
		}

		return d, nil
	}
}

// ErrorHandler defines a function type that concentrates on handling only data
// errors alone.
type ErrorHandler func(Ctx, error) (interface{}, error)

// wrapError returns a Handler which wraps a DataOnlyHandler within it, but
// passing forward all errors it receives.
func wrapError(dh ErrorHandler) Handler {
	return func(m Ctx, err error, data interface{}) (interface{}, error) {
		if err != nil {
			return dh(m, err)
		}

		return data, nil
	}
}

// ErrorOnlyHandler defines a function type that concentrates on handling only error
// replies alone.
type ErrorOnlyHandler func(interface{}) error

// wrapErrorOnly returns a Handler which wraps a ErrorOnlyHandler within it, but
// passing forward all errors it receives.
func wrapErrorOnly(dh ErrorOnlyHandler) Handler {
	return func(m Ctx, err error, data interface{}) (interface{}, error) {
		if err != nil {
			return nil, dh(data)
		}

		return data, nil
	}
}

//==============================================================================

// WrapHandlers returns a new handler where the first wraps the second with its returned
// values.
func WrapHandlers(h1 Handler, h2 Handler) Handler {
	return func(ctx Ctx, err error, data interface{}) (interface{}, error) {
		m1, e1 := h1(ctx, err, data)
		return h2(ctx, e1, m1)
	}
}

//==============================================================================

// LiftHandler defines the type the Node function returns which allows providers
// to assign handlers to use for
type LiftHandler func(...Handler) Handler

// Lift returns the a serializing handler where the passed in handler is the
// last function to be called. If the value of the argument is not a function,
// then it panics.
func Lift(handle interface{}) LiftHandler {
	mh := Wrap(handle)
	if mh == nil {
		panic("Expected handle passed into be a function")
	}

	// We will stack the handlers where one outputs becomes the input of the next.
	return func(lifts ...Handler) Handler {
		var base Handler

		lifts = append(lifts, mh)

		for i := len(lifts) - 1; i >= 0; i-- {
			if lifts[i] == nil {
				continue
			}

			if base == nil {
				base = lifts[i]
				continue
			}

			base = WrapHandlers(lifts[i], base)
		}

		return func(ctx Ctx, err error, data interface{}) (interface{}, error) {
			if base != nil {
				return base(ctx, err, data)
			}

			return data, err
		}
	}
}

// Distribute takes the output from the provided handle and distribute
// it's returned values to the provided Handlers.
func Distribute(handle interface{}) LiftHandler {
	mh := Wrap(handle)
	if mh == nil {
		panic("Expected handle passed into be a function")
	}

	// We will stack the handlers where one outputs becomes the input of the next.
	return func(lifts ...Handler) Handler {
		return func(ctx Ctx, err error, data interface{}) (interface{}, error) {
			m1, e1 := mh(ctx, err, data)

			for _, lh := range lifts {
				lh(ctx, e1, m1)
			}

			return m1, e1
		}
	}
}

// Response defines a struct for collecting the response from the Handlers.
type Response struct {
	Err   error
	Value interface{}
}

// DistributeButPack takes the output from the provided handle and distribute
// it's returned values to the provided Handlers and packs their responses in a
// slice []Response and returns that as the final response.
func DistributeButPack(handle interface{}) LiftHandler {
	mh := Wrap(handle)
	if mh == nil {
		panic("Expected handle passed into be a function")
	}

	// We will stack the handlers where one outputs becomes the input of the next.
	return func(lifts ...Handler) Handler {
		return func(ctx Ctx, err error, data interface{}) (interface{}, error) {
			var pack []Response

			m1, e1 := mh(ctx, err, data)

			for _, lh := range lifts {
				ld, le := lh(ctx, e1, m1)
				pack = append(pack, Response{
					Err:   le,
					Value: ld,
				})
			}

			return pack, nil
		}
	}
}

// Collect takes all the returned values by passing the recieved arguments and
// applying them to the handle. Where the responses of the handle is packed into
// an array of type []Collected and then returned as the response of the function.
func Collect(handle interface{}) LiftHandler {
	mh := Wrap(handle)
	if mh == nil {
		panic("Expected handle passed into be a function")
	}

	// We will stack the handlers where one outputs becomes the input of the next.
	return func(lifts ...Handler) Handler {
		return func(ctx Ctx, err error, data interface{}) (interface{}, error) {
			var pack []Response

			for _, lh := range lifts {
				m1, e1 := lh(ctx, err, data)
				d1, de := mh(ctx, e1, m1)
				pack = append(pack, Response{
					Err:   de,
					Value: d1,
				})
			}

			return pack, nil
		}
	}
}

//==============================================================================

// Pipeline defines a interface which exposes a stream like pipe which consistently
// delivers values to subscribers when executed.
type Pipeline interface {
	Exec(Ctx, error, interface{}) (interface{}, error)
	Run(Ctx, interface{}) (interface{}, error)
	Flow(Handler) Pipeline
	WithClose(func(Ctx)) Pipeline
	End(Ctx)
}

// New returns a new instance of structure that matches the Pipeline interface.
func New(main Handler) Pipeline {
	p := pipeline{
		main: main,
	}

	return &p
}

type pipeline struct {
	main   Handler
	lw     sync.RWMutex
	lines  []Handler
	closer []func(Ctx)
}

// End calls the close subscription and applies the context.
func (p *pipeline) End(ctx Ctx) {
	p.lw.RLock()
	for _, sub := range p.closer {
		sub(ctx)
	}
	p.lw.RUnlock()
}

// WithClose adds a function into the close notification lines for the pipeline.
// Returns itself for chaining.
func (p *pipeline) WithClose(h func(Ctx)) Pipeline {
	p.lw.RLock()
	p.closer = append(p.closer, h)
	p.lw.RUnlock()
	return p
}

// Flow connects another handler into the subscription list of this pipeline.
// It returns itself to allow chaining.
func (p *pipeline) Flow(h Handler) Pipeline {
	p.lw.RLock()
	p.lines = append(p.lines, h)
	p.lw.RUnlock()
	return p
}

// Run takes a context and val which it applies appropriately to the internal
// handler for the pipeline and applies the result to its subscribers.
func (p *pipeline) Run(ctx Ctx, val interface{}) (interface{}, error) {
	var res interface{}
	var err error

	if eval, ok := val.(error); ok {
		res, err = p.main(ctx, eval, nil)
	} else {
		res, err = p.main(ctx, nil, val)
	}

	p.lw.RLock()
	for _, sub := range p.lines {
		sub(ctx, err, res)
	}
	p.lw.RUnlock()

	return res, err
}

// Run takes a context, error and val which it applies appropriately to the internal
// handler for the pipeline and applies the result to its subscribers.
func (p *pipeline) Exec(ctx Ctx, er error, val interface{}) (interface{}, error) {
	res, err := p.main(ctx, er, val)

	p.lw.RLock()
	for _, sub := range p.lines {
		sub(ctx, err, res)
	}
	p.lw.RUnlock()

	return res, err
}
