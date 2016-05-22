package middleware

import (
	"github.com/influx6/faux/context"
	"github.com/influx6/faux/db/mongo"
	"github.com/influx6/faux/web/app"
)

//==============================================================================

// Log defines event logger that allows us to record events for a specific
// action that occured.
type Log interface {
	Dev(context interface{}, name string, message string, data ...interface{})
	User(context interface{}, name string, message string, data ...interface{})
	Error(context interface{}, name string, err error, message string, data ...interface{})
}

var events eventlog

// logg provides a concrete implementation of a logger.
type eventlog struct{}

// Dev logs all developer log reports.
func (l eventlog) Dev(context interface{}, name string, message string, data ...interface{}) {}

// Log logs all standard user log reports.
func (l eventlog) User(context interface{}, name string, message string, data ...interface{}) {}

// Error logs all error reports.
func (l eventlog) Error(context interface{}, name string, err error, message string, data ...interface{}) {
}

//==============================================================================

// MongoDBKey defines the key used to retrieve the configuration for creating a
// mongodb session in a context.
var MongoDBKey = "MONGO_DB"

// MongoDB returns the middleware which creates a mongo database session for
// the giving context.
func MongoDB(c mongo.Config, l Log) app.Middleware {
	md := &mongo.Mongnod{C: c, EventLog: l}

	return func(h app.Handler) app.Handler {
		return func(ctx context.Context, w *app.ResponseRequest) error {
			return h(ctx, w)
		}
	}
}
