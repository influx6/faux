package httputil

import (
	"net/http"
	"sync"

	"github.com/dimfeld/httptreemux"
)

// Router defines a interface for a structure that exposes a method for
// adding routes to an underline structure.
type Router interface {
	Handle(string, string, Handler, ...Middleware)
}

// ServeRouter defines a interface which combines the Router and http.Handler
// interface for serving requests registered with the giving routes.
type ServeRouter interface {
	Router
	http.Handler
}

// MuxRouter defines a implementation of the Router interface which exposes
// the necessary methods for serving requests.
type MuxRouter struct {
	treemux   *httptreemux.TreeMux
	generator TreemuxHandlerMW
	pool      *sync.Pool
}

// NewMuxRouter returns a new instance of a MuxRouter.
func NewMuxRouter(eh ErrorHandler, ops ...Options) *MuxRouter {
	var mx MuxRouter
	mx.treemux = httptreemux.New()

	mx.generator, mx.pool = PoolTreemuxHandler(eh, ops...)

	return &mx
}

// Treemux returns the underline httptreemux.TreeMux router.
func (m *MuxRouter) Treemux() *httptreemux.TreeMux {
	return m.treemux
}

// ServeHTTP services the underline requests for response desired.
func (m *MuxRouter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	m.treemux.ServeHTTP(w, r)
}

// Handle registers the giving router into the underline router.
func (m *MuxRouter) Handle(method string, path string, hl Handler, mw ...Middleware) {
	m.treemux.Handle(method, path, m.generator(hl, mw...))
}
