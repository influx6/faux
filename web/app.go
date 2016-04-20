package web

import (
	"net/http"

	"github.com/dimfeld/httptreemux"
	"github.com/influx6/faux/context"
	"github.com/influx6/faux/sumex"
)

//==============================================================================

// Log defines event logger that allows us to record events for a specific
// action that occured.
type Log interface {
	Log(context interface{}, name string, message string, data ...interface{})
	Error(context interface{}, name string, err error, message string, data ...interface{})
}

//==============================================================================

// Handler provides the signature for handler providers.
type Handler func(context.Context, http.ResponseWriter, *http.Request) error

// Middleware defines the middleware signature for creating middlewares.
type Middleware func(Handler) Handler

//==============================================================================

// App provides the core provider for creating a server provider using
// http.Router and the sumex.Stream management system.
type App struct {
	sumex.Stream
	tree    *httptreemux.TreeMux
	gm      []Middleware
	options httptreemux.HandlerFunc
	headers map[string]string
}

// New returns a new App instance.
func New(l Log, cors bool, w int, mh ...Middleware) *App {
	app := App{
		gm:      mh,
		tree:    httptreemux.New(),
		headers: make(map[string]string),
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

	app.tree.OptionsHandler = app.options
	app.Stream = sumex.New(l, w, &app)

	return &app
}

// ServeHTTP implements the http.Handler interface which allows us
// provide a server muxilator.
func (a *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.tree.ServeHTTP(w, r)
}

// Do performs the needed operation for handling a app-server.
func (a *App) Do(ctx context.Context, err error, data interface{}) (interface{}, error) {

	return nil, nil
}
