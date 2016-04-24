package middleware

import (
	"github.com/influx6/faux/context"
	"github.com/influx6/faux/web/app"
)

// DBKey defines the key used to retrieve the configuration for creating a mongodb
// session.
var DBKey = "db"

// MongoDB returns the middleware which creates a mongo database session for
// the giving context.
func MongoDB(h app.Handler) app.Handler {
	return func(ctx context.Context, w *app.ResponseRequest, params app.Param) error {

		return h(ctx, w, params)
	}
}
