package httputil

import (
	"net/http"
	"sync"
)

// Handler defines a function type to process a giving request.
type Handler func(*Context) error

// HandlerMW defines a function which wraps a provided http.handlerFunc
// which encapsulates the original for a underline operation.
type HandlerMW func(Handler, ...Middleware) http.HandlerFunc

// Middlware defines a function type which is used to create a chain
// of handlers for processing giving request.
type Middleware func(next Handler) Handler

// IdentityHandler defines a Handler function that returns a nil error.
func IdentityHandler(c *Context) error {
	return nil
}

// WrapHandler attempts to wrap provided http.HandlerFunc to return a httputil.Context.
func WrapHandler(fx http.HandlerFunc) Handler {
	return func(ctx *Context) error {
		fx(ctx.Response(), ctx.Request())
		return nil
	}
}

// MixHandler wraps two provided Handler and returns a new Middleware.
func MixHandler(mo, mi Handler) Handler {
	return func(c *Context) error {
		if err := mo(c); err != nil {
			return err
		}

		return mi(c)
	}
}

type handlerHost struct {
	fn http.HandlerFunc
}

// ServeHTTP services the giving request using the underline http.handlerFunc.
func (h handlerHost) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.fn(w, r)
}

// HandlerFuncToHandler defines a function which returns a http.Handler from the
// provided http.HandlerFunc.
func HandlerFuncToHandler(hl http.HandlerFunc) http.Handler {
	return &handlerHost{
		fn: hl,
	}
}

// Pool defines a function which will return a http.HandlerFunc which will
// receive new Context objects with the provided options applied and it generated
// from a sync.Pool which will be used to retrieve and create new Context objects.
// WARNING: When the http.handlerFunc returned by the returned HandlerX function,
// the Context created will be reset and put back into the pull. So ensure calls
// do not escape the http.HandlerFunc returned.
func Pool(errHandler ErrorHandler, ops ...Options) HandlerMW {
	contextPool := sync.Pool{
		New: func() interface{} {
			return NewContext(ops...)
		},
	}

	return func(handle Handler, mw ...Middleware) http.HandlerFunc {
		middleware := MW(mw...)

		return func(w http.ResponseWriter, r *http.Request) {
			ctx, ok := contextPool.Get().(*Context)
			if !ok {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			// Reset Request and Response for context.
			ctx.Reset(r, &Response{Writer: w})

			defer ctx.Reset(nil, nil)
			defer contextPool.Put(ctx)

			ctx.InitForms()

			if err := middleware(ctx); err != nil && errHandler != nil {
				errHandler(err, ctx)
				return
			}

			if err := handle(ctx); err != nil && errHandler != nil {
				errHandler(err, ctx)
				return
			}
		}
	}
}

// MWi combines multiple Middleware to return a new Middleware.
func MWi(mos ...Middleware) Middleware {
	var initial Middleware

	for _, mw := range mos {
		if initial == nil {
			initial = mw
			continue
		}

		initial = DMW(initial, mw)
	}

	return initial
}

// MW combines multiple Middleware to return a single Handler.
func MW(mos ...Middleware) Handler {
	return MWi(mos...)(IdentityHandler)
}

// DMW combines two middleware and returns a single Handler.
func DMW(mo, mi Middleware) Middleware {
	return func(next Handler) Handler {
		handler := mo(mi(IdentityHandler))

		return func(ctx *Context) error {
			if err := handler(ctx); err != nil {
				return err
			}

			return next(ctx)
		}
	}
}
