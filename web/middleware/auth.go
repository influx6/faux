package middleware

import (
	"github.com/influx6/faux/context"
	"github.com/influx6/faux/web/app"
)

// AuthKey defines the auth key which is used to store the authentication details
// for the usage of the BasicAuth middleware.
var AuthKey = "BASICAUTH"

// BasicAuth defines a middleware for adding BasicAuth into the response for a request.
func BasicAuth(h app.Handler) app.Handler {
	return func(ctx context.Context, w *app.ResponseRequest, params app.Param) error {

		return h(ctx, w, params)
	}
}
