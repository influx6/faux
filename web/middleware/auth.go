package middleware

import (
	"errors"

	"github.com/influx6/faux/context"
	"github.com/influx6/faux/web/app"
)

// ErrNotAuthorized is returned when the Authorized is on and the authorization
// token is not found or invalid.
var ErrNotAuthorized = errors.New("Not Authorized")

// AuthKey defines the auth key which is used to store the authentication details
// for the usage of the BasicAuth middleware.
const AuthKey = "BASICAUTH"

// BasicAuth defines a middleware for adding BasicAuth into the response for a request.
func BasicAuth(h app.Handler) app.Handler {
	return func(ctx context.Context, w *app.ResponseRequest) error {
		_, ok := ctx.Get(AuthKey)
		if !ok {
			return h(ctx, w)
		}

		token := w.Header().Get("Authorization")

		if tokenLen := len(token); tokenLen < 5 || token[0:5] != "Basic" {
			return ErrNotAuthorized
		}

		return h(ctx, w)
	}
}
