package ops

import "github.com/influx6/faux/pub"

// Operations defines sets of routine events that can be created onto a pub.Node.
type Operations interface {
	pub.Node
	Collect() Operations
}
