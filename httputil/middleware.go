package httputil

import "net/http"

// Handler defines a function type to process a giving request.
type Handler func(*Context) error

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

// HandlerMW wraps two provided Handler and returns a new Middleware.
func HandlerMW(mo, mi Handler) Handler {
	return func(c *Context) error {
		if err := mo(c); err != nil {
			return err
		}

		return mi(c)
	}
}

// MW combines multiple Middleware to return a single Handler.
func MW(mos ...Middleware) Handler {
	var initial Middleware

	for _, mw := range mos {
		if initial == nil {
			initial = mw
			continue
		}

		initial = DMW(initial, mw)
	}

	return initial(IdentityHandler)
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
