package web

import (
	"fmt"
	"net/http"
	"strings"
	"sync"

	"github.com/dimfeld/httptreemux"
	"github.com/influx6/faux/context"
	"github.com/influx6/faux/sumex"
)

//==============================================================================

// Route defines a route interface for the web package, this allows us to
// register this routes.
type Route interface {
	Register(context.Context, *httptreemux.TreeMux)
}

//==============================================================================

// ContentPipe defines a immutable pipe which registers content-type based
// routes which gets handled by a muxilator, these allows appending
// multiple responders based on content type for a server.
type ContentPipe struct {
	prev    *ContentPipe
	content string
	handle  Handler
}

// Get returns the appropriate handler for the specific content-type if found.
func (c *ContentPipe) Get(req *http.Request) (handler Handler, found bool) {
	found = strings.Contains(req.Header.Get("Content-Type"), c.content)
	if !found && c.prev == nil {
		return c.prev.Get(req)
	}

	if !found {
		return
	}

	handler = c.handle
	return
}

// Append adds a new content type into the change and allows us to provide
// a possible request chain where a request is processed else failed if
// the handler for its content type does not exists.
func (c *ContentPipe) Append(content string, handler Handler) *ContentPipe {
	ch := ContentPipe{
		prev:    c,
		content: content,
		handle:  handler,
	}

	return &ch
}

// ContentResponse provides a concurrently-safe router for handle response
// functions for a routes content-type, it uses a mutex to safe-guard the
// addition and use of new content-providers.
type ContentResponse struct {
	mu   sync.Mutex
	pipe *ContentPipe
}

// Add adds a new handler for a specific content type.
func (c *ContentResponse) Add(content string, handler Handler) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.pipe = c.pipe.Append(content, handler)
}

// Do works with the content response pipes to provide the appropriate response
// for the giving route.
func (c *ContentResponse) Do(ctx context.Context, w http.ResponseWriter, r *http.Request, params Param) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	handler, ok := c.pipe.Get(r)
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return fmt.Errorf("Unknown Content-Type[%s]", r.Header.Get("Content-Type"))
	}

	return handler(ctx, w, r, params)
}

//==============================================================================

// route implements the Route interface, registering a route as needed.
type route struct {
	verb    string
	path    string
	handler Handler
}

// Register registers the route with the giving path mux.
func (r *route) Register(ctx context.Context, tree *httptreemux.TreeMux) {
	tree.Handle(r.verb, r.path, func(w http.ResponseWriter, rq *http.Request, param map[string]string) {
		if err := r.handler(ctx, w, rq, Param(param)); err != nil {
			RenderError(err, rq, w)
		}
	})
}

//==============================================================================

// PageRoute adds a new route to a app http Server.
func PageRoute(app interface{}, verb string, path string, h Handler) {
	rm := route{
		verb:    verb,
		path:    path,
		handler: h,
	}

	switch app.(type) {
	case sumex.Stream:
		(app.(sumex.Stream)).Data(nil, &rm)
	case *App:
		(app.(*App)).Do(nil, nil, &rm)
	}
}

// ContentRoute adds a new route to a app http Server.
func ContentRoute(app interface{}, c *ContentResponse, verb string, path string, h Handler) {
	PageRoute(app, verb, path, c.Do)
}

//==============================================================================
