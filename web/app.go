package web

import (
	"net/http"

	"github.com/dimfeld/httptreemux"
	"golang.org/influx/faux/context"
)

// Handle provides the signature for handler providers.
type Handle func(context.Context, http.ResponseWriter, *http.Request) error

// Middleware defines the middleware signature for creating middlewares.
type Middleware func(Handler) Handler

// App provides the core provider for creating a server provider using
// http.Router and the sumex.Stream management system.
type App struct {
	CORS bool
	tree *httptreemux.TreeMux
}

// Do performs the needed operation for handling a app-server.
func (a *App) Do(ctx context.Context, err error, data interface{}) (interface{}, error) {

	return nil, nil
}
