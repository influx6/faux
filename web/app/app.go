package app

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/dimfeld/httptreemux"
	"github.com/influx6/faux/context"
)

//==============================================================================

// Log defines event logger that allows us to record events for a specific
// action that occured.
type Log interface {
	Log(context interface{}, name string, message string, data ...interface{})
	Error(context interface{}, name string, err error, message string, data ...interface{})
}

var events eventlog

// logg provides a concrete implementation of a logger.
type eventlog struct{}

// Log logs all standard log reports.
func (l eventlog) Log(context interface{}, name string, message string, data ...interface{}) {}

// Error logs all error reports.
func (l eventlog) Error(context interface{}, name string, err error, message string, data ...interface{}) {
}

//==============================================================================

// Param defines the map of values to be handled by the provider.
type Param map[string]string

// Get returns the given value for a specific key and returns a bool to indicate
// if it was found.
func (p Param) Get(key string) (val string, found bool) {
	val, found = p[key]
	return
}

// GetBool returns a bool value if possible from the value of a giving key if it
// exits.
func (p Param) GetBool(key string) (item bool, err error) {
	val, ok := p[key]
	if !ok {
		err = errors.New("Not Found")
		return
	}

	item, err = strconv.ParseBool(val)
	return
}

// GetFloat returns a float value if possible from the value of a giving key if
// it exits.
func (p Param) GetFloat(key string) (item float64, err error) {
	val, ok := p[key]
	if !ok {
		err = errors.New("Not Found")
		return
	}

	item, err = strconv.ParseFloat(val, 64)
	return
}

// GetInt returns a int value if possible from the value of a giving key if it
// exits.
func (p Param) GetInt(key string) (item int, err error) {
	val, ok := p[key]
	if !ok {
		err = errors.New("Not Found")
		return
	}

	item, err = strconv.Atoi(val)
	return
}

//==============================================================================

// Handler provides the signature for handler providers.
type Handler func(context.Context, *ResponseRequest, Param) error

// Middleware defines the middleware signature for creating middlewares.
type Middleware func(Handler) Handler

//==============================================================================

// App provides the core provider for creating a server provider using
// http.Router and the sumex.Stream management system.
type App struct {
	*httptreemux.TreeMux
	log     Log
	gm      []Middleware
	options httptreemux.HandlerFunc
	headers map[string]string
}

// New returns a new App instance.
func New(l Log, cors bool, m map[string]string, mh ...Middleware) *App {
	if m == nil {
		m = make(map[string]string)
	}

	if l == nil {
		l = events
	}

	app := App{
		TreeMux: httptreemux.New(),
		log:     l,
		gm:      mh,
		headers: m,
	}

	if cors {
		app.options = func(w http.ResponseWriter, _ *http.Request, _ map[string]string) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			w.Header().Set("Access-Control-Max-Age", "86400")
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
		}

		app.headers["Access-Control-Allow-Origin"] = "*"
	}

	app.TreeMux.OptionsHandler = app.options

	return &app
}

// Handle decorates the internal TreeMux Handle function to apply global handlers into the system.
func (a *App) Handle(ctx context.Context, verb string, path string, h Handler, m ...Middleware) {

	// Apply the global handlers which calls its next handler in reverse order.
	for i := len(a.gm); i >= 0; i++ {
		h = a.gm[i](h)
	}

	// Apply the local handlers which calls its next handler in reverse order.
	for i := len(m); i >= 0; i++ {
		h = m[i](h)
	}

	a.TreeMux.Handle(verb, path, func(w http.ResponseWriter, r *http.Request, params map[string]string) {
		h(ctx.New(), &ResponseRequest{ResponseWriter: w, R: r}, Param(params))
	})
}

// ServeHTTP implements the http.Handler interface which allows us
// provide a server muxilator.
func (a *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if a.headers != nil {
		for key, val := range a.headers {
			w.Header().Set(key, val)
		}
	}
	a.TreeMux.ServeHTTP(w, r)
}

// Do performs the needed operation for handling a app-server.
func (a *App) Do(ctx context.Context, err error, data interface{}) (interface{}, error) {
	if err != nil {
		return nil, err
	}

	wrap := func(h Handler) Handler {
		// Apply the global handlers which calls its next handler in reverse order.
		for i := len(a.gm); i >= 0; i++ {
			h = a.gm[i](h)
		}

		return h
	}

	switch data.(type) {
	case Route:
		(data.(Route)).Register(ctx, a.TreeMux, wrap)

		return nil, nil
	default:
		return nil, errors.New("Unknwon Action")
	}
}
