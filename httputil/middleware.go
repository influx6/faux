package httputil

import (
	"net/http"
	"strings"

	"github.com/dimfeld/httptreemux"
	"github.com/influx6/faux/metrics"
)

// Handler defines a function type to process a giving request.
type Handler func(*Context) error

// ErrorHandler defines a function type which sets giving respnse to a Response object.
type ErrorHandler func(error, *Context) error

// HandlerMW defines a function which wraps a provided http.handlerFunc
// which encapsulates the original for a underline operation.
type HandlerMW func(Handler, ...Middleware) Handler

// HandlerFuncMW defines a function which wraps a provided http.handlerFunc
// which encapsulates the original for a underline operation.
type HandlerFuncMW func(Handler, ...Middleware) http.HandlerFunc

// TreemuxHandlerMW defines a function which wraps a provided http.handlerFunc
// which encapsulates the original for a underline operation.
type TreemuxHandlerMW func(Handler, ...Middleware) httptreemux.HandlerFunc

// TreeMuxHandler defines a function type for the httptreemux.Handler type.
type TreeMuxHandler func(http.ResponseWriter, *http.Request, map[string]string)

// Middleware defines a function type which is used to create a chain
// of handlers for processing giving request.
type Middleware func(next Handler) Handler

// IdentityHandler defines a Handler function that returns a nil error.
func IdentityHandler(c *Context) error {
	return nil
}

// MetricsMW defines a MW function that adds giving metrics.Metric to all context
// before calling next handler.
func MetricsMW(m metrics.Metrics) Middleware {
	return func(next Handler) Handler {
		return func(ctx *Context) error {
			return next(Apply(ctx, SetMetrics(m)))
		}
	}
}

// NetworkAuthenticationNeeded implements a Handler which returns http.StatusNetworkAuthenticationRequired always.
func NetworkAuthenticationNeeded(ctx *Context) error {
	ctx.Status(http.StatusNetworkAuthenticationRequired)
	return nil
}

// NoContentRequest implements a Handler which returns http.StatusNoContent always.
func NoContentRequest(ctx *Context) error {
	ctx.Status(http.StatusNoContent)
	return nil
}

// OKRequest implements a Handler which returns http.StatusOK always.
func OKRequest(ctx *Context) error {
	ctx.Status(http.StatusOK)
	return nil
}

// BadRequestWithError implements a Handler which returns http.StatusBagRequest always.
func BadRequestWithError(err error, ctx *Context) error {
	if err != nil {
		if httperr, ok := err.(HTTPError); ok {
			http.Error(ctx.Response(), httperr.Error(), httperr.Code)
			return nil
		}
		http.Error(ctx.Response(), err.Error(), http.StatusBadRequest)
	}
	return nil
}

// BadRequest implements a Handler which returns http.StatusBagRequest always.
func BadRequest(ctx *Context) error {
	ctx.Status(http.StatusBadRequest)
	return nil
}

// NotFound implements a Handler which returns http.StatusNotFound always.
func NotFound(ctx *Context) error {
	ctx.Status(http.StatusNotFound)
	return nil
}

// StripPrefixMW returns a middleware which strips the URI of the request of
// the provided Prefix. All prefix must come in /prefix/ format.
func StripPrefixMW(prefix string) Middleware {
	if !strings.HasPrefix(prefix, "/") {
		prefix = "/" + prefix
	}

	return func(next Handler) Handler {
		return func(ctx *Context) error {

			req := ctx.Request()
			reqURL := req.URL.Path

			if !strings.HasPrefix(reqURL, "/") {
				reqURL = "/" + reqURL
			}

			req.URL.Path = strings.TrimPrefix(reqURL, prefix)
			return next(ctx)
		}
	}
}

// HTTPRedirect returns a Handler which always redirect to the given path.
func HTTPRedirect(to string, code int) Handler {
	return func(ctx *Context) error {
		return ctx.Redirect(code, to)
	}
}

// Then calls the next Handler after the condition handler returns without error.
func Then(condition Handler, nexts ...Handler) Handler {
	if len(nexts) == 0 {
		return condition
	}

	return func(c *Context) error {
		if err := condition(c); err != nil {
			return err
		}

		for _, next := range nexts {
			if err := next(c); err != nil {
				return err
			}
		}
		return nil
	}
}

// HTTPConditionFunc retusn a handler where a Handler is used as a condition where if the handler
// returns an error then the errorAction is called else the noerrorAction gets called with
// context. This allows you create a binary switch where the final action is based on the
// success of the first. Generally if you wish to pass info around, use the context.Bag()
// to do so.
func HTTPConditionFunc(condition Handler, noerrorAction, errorAction Handler) Handler {
	return func(ctx *Context) error {
		if err := condition(ctx); err != nil {
			ctx.Metrics().Emit(metrics.Error(err).WithMessage("HTTPConditionFunc").With("httputil_handler_error", err))
			return errorAction(ctx)
		}
		return noerrorAction(ctx)
	}
}

// HTTPConditionErrorFunc returns a handler where a condition Handler is called whoes result if with an error
// is passed to the errorAction for execution else using the noerrorAction. Differs from HTTPConditionFunc
// due to the assess to the error value.
func HTTPConditionErrorFunc(condition Handler, noerrorAction Handler, errorAction ErrorHandler) Handler {
	return func(ctx *Context) error {
		if err := condition(ctx); err != nil {
			ctx.Metrics().Emit(metrics.Error(err).WithMessage("HTTPConditionFunc").With("httputil_handler_error", err))
			return errorAction(err, ctx)
		}
		return noerrorAction(ctx)
	}
}

// ErrorsAsResponse returns a Handler which will always write out any error that
// occurs as the response for a request if any occurs.
func ErrorsAsResponse(code int, next Handler) Handler {
	return func(ctx *Context) error {
		if err := next(ctx); err != nil {
			ctx.Metrics().Emit(metrics.Error(err).WithMessage("HTTPConditionFunc").With("httputil_handler_error", err))
			if httperr, ok := err.(HTTPError); ok {
				http.Error(ctx.Response(), httperr.Error(), httperr.Code)
				return err
			}

			if code <= 0 {
				code = http.StatusBadRequest
			}

			http.Error(ctx.Response(), err.Error(), code)
			return err
		}
		return nil
	}
}

// HTTPConditionsFunc returns a Handler where if an error occurs would match the returned
// error with a Handler to be runned if the match is found.
func HTTPConditionsFunc(condition Handler, noerrAction Handler, errCons ...ErrConditions) Handler {
	return func(ctx *Context) error {
		if err := condition(ctx); err != nil {
			ctx.Metrics().Emit(metrics.Error(err).WithMessage("HTTPConditionsFunc").With("httputil_handler_error", err))
			for _, errcon := range errCons {
				if errcon.Match(err) {
					return errcon.Handle(ctx)
				}
			}
			return err
		}
		return noerrAction(ctx)
	}
}

// ErrConditions defines a condition which matches expected error
// for performing giving action.
type ErrConditions interface {
	Match(error) bool
	Handle(*Context) error
}

// ErrorCondition defines a type which sets the error that occurs and the handler to be called
// for such an error.
type ErrorCondition struct {
	Err error
	Fn  Handler
}

// ErrCondition returns ErrConditon using provided arguments.
func ErrCondition(err error, fn Handler) ErrorCondition {
	return ErrorCondition{
		Err: err,
		Fn:  fn,
	}
}

// Handler calls the internal Handler with provided Context returning error.
func (ec ErrorCondition) Handler(ctx *Context) error {
	return ec.Fn(ctx)
}

// Match validates the provided error matches expected error.
func (ec ErrorCondition) Match(err error) bool {
	return ec.Err == err
}

// FnErrorCondition defines a type which sets the error that occurs and the handler to be called
// for such an error.
type FnErrorCondition struct {
	Fn  Handler
	Err func(error) bool
}

// FnErrCondition returns ErrConditon using provided arguments.
func FnErrCondition(err func(error) bool, fn Handler) FnErrorCondition {
	return FnErrorCondition{
		Err: err,
		Fn:  fn,
	}
}

// Handler calls the internal Handler with provided Context returning error.
func (ec FnErrorCondition) Handler(ctx *Context) error {
	return ec.Fn(ctx)
}

// Match validates the provided error matches expected error.
func (ec FnErrorCondition) Match(err error) bool {
	return ec.Err(err)
}

// LogMW defines a log middleware function which wraps a Handler
// and logs what request and response was sent incoming.
func LogMW(next Handler) Handler {
	return func(ctx *Context) error {
		m := ctx.Metrics()
		if m == nil && next == nil {
			return nil
		}
		if m == nil && next != nil {
			return next(ctx)
		}

		var err error
		req := ctx.Request()
		res := ctx.Response()
		res.After(func() {
			if err != nil {
				m.Emit(metrics.Error(err).
					WithMessage("Outgoing HTTP Response").
					With("method", req.Method).
					With("status", res.Status).
					With("header", res.Header()).
					With("path", req.URL.Path).
					With("host", req.Host).
					With("remote", req.RemoteAddr).
					With("agent", req.UserAgent()).
					With("request", req.RequestURI).
					With("outgoing-content-length", res.Size).
					With("incoming-content-length", req.ContentLength).
					With("proto", req.Proto))
				return
			}

			m.Emit(metrics.Info("Outgoing HTTP Response").
				With("method", req.Method).
				With("status", res.Status).
				With("header", res.Header()).
				With("path", req.URL.Path).
				With("host", req.Host).
				With("remote", req.RemoteAddr).
				With("agent", req.UserAgent()).
				With("request", req.RequestURI).
				With("outgoing-content-length", res.Size).
				With("incoming-content-length", req.ContentLength).
				With("proto", req.Proto))
		})

		m.Emit(metrics.Info("Incoming HTTP Request").
			With("method", req.Method).
			With("path", req.URL.Path).
			With("tls", req.TLS != nil).
			With("host", req.Host).
			With("header", req.Header).
			With("remote", req.RemoteAddr).
			With("agent", req.UserAgent()).
			With("request", req.RequestURI).
			With("content-length", req.ContentLength).
			With("proto", req.Proto))

		err = next(ctx)
		return err
	}
}

// IdentityMW defines a Handler function that returns a the next Handler passed to it.
func IdentityMW(next Handler) Handler {
	return next
}

// WrapHandler attempts to wrap provided http.HandlerFunc to return a httputil.Context.
func WrapHandler(fx http.HandlerFunc) Handler {
	return func(ctx *Context) error {
		fx(ctx.Response(), ctx.Request())
		return nil
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

// MuxHandler defines a function which will return a Handler which will
// be used to handle a request.
func MuxHandler(errHandler ErrorHandler, handle Handler, mw ...Middleware) Handler {
	middleware := MW(mw...)

	return func(ctx *Context) error {
		if err := middleware(ctx); err != nil {
			if errHandler != nil {
				errHandler(err, ctx)
				return nil
			}

			return err
		}

		if err := handle(ctx); err != nil {
			if errHandler != nil {
				errHandler(err, ctx)
				return nil
			}

			return err
		}

		return nil
	}
}

// TreemuxHandler defines a function which will return a http.HandlerFunc which will
// receive new Context objects with the provided options applied and it generated
// from a sync.Pool which will be used to retrieve and create new Context objects.
// WARNING: When the http.handlerFunc returned by the returned HandlerX function,
// the Context created will be reset and put back into the pull. So ensure calls
// do not escape the http.HandlerFunc returned.
func TreemuxHandler(errHandler ErrorHandler, ops ...Options) TreemuxHandlerMW {
	return func(handle Handler, mw ...Middleware) httptreemux.HandlerFunc {
		middleware := MW(mw...)

		return func(w http.ResponseWriter, r *http.Request, params map[string]string) {
			ctx := NewContext(ops...)
			ctx.Reset(r, &Response{Writer: w})
			defer ctx.Reset(nil, nil)
			defer ctx.ClearFlashMessages()

			for key, val := range params {
				ctx.Bag().Set(key, val)
			}

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

	if len(mos) == 0 {
		initial = IdentityMW
	}

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
