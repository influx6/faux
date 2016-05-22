package middleware

import (
	"github.com/influx6/faux/context"
	"github.com/influx6/faux/web/app"
)

// MongoDBKey defines the key used to retrieve the configuration for creating a
// mongodb session in a context.
var MongoDBKey = "MONGO_DB"

// MongoDB returns the middleware which creates a mongo database session for
// the giving context.
func MongoDB(h app.Handler) app.Handler {
	return func(ctx context.Context, w *app.ResponseRequest) error {
		return h(ctx, w)
	}
}
