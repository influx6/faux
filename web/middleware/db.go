package middleware

import (
	"github.com/ardanlabs/kit/cfg"
	"github.com/influx6/faux/context"
	"github.com/influx6/faux/db/mongo"
	"github.com/influx6/faux/web/app"
)

//==============================================================================

// Logs defines event logger that allows us to record events for a specific
// action that occured.
type Logs interface {
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

// MongoSessionKey defines the key which locks the mongo session provided into
// the passed context using the mongo middleware.
const MongoSessionKey = "MONGO_DB_SESSION"

// MongoDatabaseKey defines the key which locks the mongo database to be used
// directly from the giving mongo configuration. It lets you access the proper
// db without knowledge of the name by using the config DB attribute.
const MongoDatabaseKey = "MONGO_DB_DATABASE"

// MongoDB returns the middleware which creates a mongo database session for
// the giving context.
func MongoDB(c mongo.Config, l Logs) app.Middleware {
	if l == nil {
		l = events
	}

	md := &mongo.Mongnod{C: c, Log: l}

	return func(h app.Handler) app.Handler {
		return func(ctx context.Context, w *app.ResponseRequest) error {
			db, ses, err := md.New(nil)
			if err != nil {
				l.Error("app.Middleware", "MongoDB", err, "Completed")
				return err
			}

			// End the session once this is done.
			defer ses.Close()

			if err := h(ctx.WithValue(MongoSessionKey, ses).WithValue(MongoDatabaseKey, db), w); err != nil {
				l.Error("app.Middleware", "MongoDB", err, "Completed")
				return err
			}

			return nil
		}
	}
}

// MongoDBEnv returns a middleware which loads its configuration from the
// environment variables of the host system.
func MongoDBEnv(configName string, l Logs) app.Middleware {
	if l == nil {
		l = events
	}

	// Initialize the configuration system to retrieve environment variavles.
	cfg.Init(cfg.EnvProvider{Namespace: configName})

	return MongoDB(mongo.Config{
		Host:     cfg.MustString("MONGO_HOST"),
		AuthDB:   cfg.MustString("MONGO_AUTHDB"),
		DB:       cfg.MustString("MONGO_DB"),
		User:     cfg.MustString("MONGO_USER"),
		Password: cfg.MustString("MONGO_PASS"),
	}, l)
}
