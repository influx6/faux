// Package pub provides a functional reactive pubsub structure to leverage a
// pure function style reactive behaviour. Originally pulled from pub.Node.
// NOTE: Any use of "asynchronouse" actually means to "run within a goroutine",
// and inversely, the use of "synchronouse" means to run it within the current
// goroutine, generally referred to as "main", or in other words.
package pub

import (
	"errors"
	"reflect"
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

// MagicHandler returns a new Handler wrapping the provided value as needed if
// it matches its DataHandler, ErrorHandler, Handler or magic function type.
// MagicFunction type is a function which follows this type form:
// func(context.Context, error, <CustomType>).
func MagicHandler(node interface{}) Handler {
	var hl Handler

	switch node.(type) {
	case func(Ctx, error, interface{}) (interface{}, error):
		hl = node.(func(Ctx, error, interface{}) (interface{}, error))
	case func(Ctx, interface{}) (interface{}, error):
		hl = WrapData(node.(func(Ctx, interface{}) (interface{}, error)))
	case func(interface{}) interface{}:
		hl = WrapDataOnly(node.(func(interface{}) interface{}))
	case func(Ctx, error) (interface{}, error):
		hl = WrapError(node.(func(Ctx, error) (interface{}, error)))
	case func(interface{}):
		hl = WrapJustData(node.(func(interface{})))
	case func(error):
		hl = WrapJustError(node.(func(error)))
	case func(error) error:
		hl = WrapErrorReturn(node.(func(error) error))
	case func() interface{}:
		hl = WrapNoData(node.(func() interface{}))
	case func(interface{}) error:
		hl = WrapErrorOnly(node.(func(interface{}) error))
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

			if err != nil {

				if !isCustorm {
					resArgs := tm.Call([]reflect.Value{ma, reflect.ValueOf(err), dZero})
					if len(resArgs) < 1 {
						return nil, nil
					}

					if len(resArgs) == 1 {
						dx, ok := ((resArgs[0].Interface()).(error))
						if ok {
							return nil, dx
						}

						return dx, nil
					}

					return resArgs[0].Interface(), (resArgs[1].Interface().(error))
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
				dArgs := []reflect.Value{ma, dZeroError, mVal}
				resArgs := tm.Call(dArgs)
				if len(resArgs) < 1 {
					return nil, nil
				}

				if len(resArgs) == 1 {
					dx, ok := ((resArgs[0].Interface()).(error))
					if ok {
						return nil, dx
					}

					return dx, nil
				}

				return resArgs[0].Interface(), (resArgs[1].Interface().(error))
			}

			dArgs := []reflect.Value{ma, mVal}

			resArgs := tm.Call(dArgs)
			if len(resArgs) < 1 {
				return nil, nil
			}

			if len(resArgs) == 1 {
				dx, ok := ((resArgs[0].Interface()).(error))
				if ok {
					return nil, dx
				}

				return dx, nil
			}

			return resArgs[0].Interface(), (resArgs[1].Interface().(error))
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

// WrapData returns a Handler which wraps a DataHandler within it, but
// passing forward all errors it receives.
func WrapData(dh DataHandler) Handler {
	return func(m Ctx, err error, data interface{}) (interface{}, error) {
		if err != nil {
			return nil, err
		}

		return dh(m, data)
	}
}

// NoDataHandler defines an Handler which allows a return value when called
// but has no data passed in.
type NoDataHandler func() interface{}

// WrapNoData returns a Handler which wraps a NoDataHandler within it, but
// forwards all errors it receives. It calls its internal function
// with no arguments taking the response and sending that out.
func WrapNoData(dh NoDataHandler) Handler {
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

// WrapDataOnly returns a Handler which wraps a DataOnlyHandler within it, but
// passing forward all errors it receives.
func WrapDataOnly(dh DataOnlyHandler) Handler {
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

// WrapJustData wraps a JustDataHandler and returns it as a Handler.
func WrapJustData(dh JustDataHandler) Handler {
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

// WrapJustError returns a Handler which wraps a DataOnlyHandler within it, but
// passing forward all errors it receives.
func WrapJustError(dh JustErrorHandler) Handler {
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

// WrapErrorReturn returns a Handler which wraps a DataOnlyHandler within it, but
// passing forward all errors it receives.
func WrapErrorReturn(dh ErrorReturnHandler) Handler {
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

// WrapError returns a Handler which wraps a DataOnlyHandler within it, but
// passing forward all errors it receives.
func WrapError(dh ErrorHandler) Handler {
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

// WrapErrorOnly returns a Handler which wraps a ErrorOnlyHandler within it, but
// passing forward all errors it receives.
func WrapErrorOnly(dh ErrorOnlyHandler) Handler {
	return func(m Ctx, err error, data interface{}) (interface{}, error) {
		if err != nil {
			return nil, dh(data)
		}

		return data, nil
	}
}

//==============================================================================

// Wrap returns a new handler where the first wraps the second with its returned
// values.
func Wrap(h1 Handler, h2 Handler) Handler {
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
	mh := MagicHandler(handle)
	if mh == nil {
		panic("Expected handle passed into be a function")
	}

	// We will stack the handlers where one outputs becomes the input of the next.
	return func(lifts ...Handler) Handler {
		var base Handler

		for i := len(lifts) - 1; i >= 0; i-- {
			if base == nil {
				base = lifts[i]
				continue
			}

			base = Wrap(lifts[i], base)
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
	mh := MagicHandler(handle)
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
	mh := MagicHandler(handle)
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
	mh := MagicHandler(handle)
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
