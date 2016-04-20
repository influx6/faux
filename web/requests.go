package web

import (
	"github.com/influx6/faux/context"
)

// Route defines a route interface for the web package, this allows us to
// register this routes.
type Route interface {
	Register(context.Context)
}
