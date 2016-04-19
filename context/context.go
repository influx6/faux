// Package context is built out of my desire to understand the http context
// library and as an experiement in such a library works.
package context

import "time"

//==============================================================================

// Fields defines a map of key:value pairs.
type Fields map[interface{}]interface{}

//==============================================================================

// nilPair defines a nil starting pair.
var nilPair = (*Pair)(nil)

// Pair defines a struct for storing a linked pair of key and values.
type Pair struct {
	prev  *Pair
	key   interface{}
	value interface{}
}

// Append returns a new Pair with the giving key and with the provded Pair set as
// it's previous link.
func Append(p *Pair, key, value interface{}) *Pair {
	return &Pair{
		prev:  p,
		key:   key,
		value: value,
	}
}

// Fields returns all internal pair data as a map.
func (p *Pair) Fields() Fields {
	var f Fields

	if p.prev == nil {
		f = make(Fields)
		f[p.key] = p.value
		return f
	}

	f = p.prev.Fields()
	f[p.key] = p.value
	return f
}

// Get collects the value of a key if it exists.
func (p *Pair) Get(key interface{}) (value interface{}, found bool) {
	if p == nil {
		return
	}

	if p.key == key {
		return p.value, true
	}

	if p.prev == nil {
		return
	}

	return p.prev.Get(key)
}

//==============================================================================

// Context defines an interface for a context providers which allows us to
// build passable context around.
type Context interface {

	//Cancel cancels out the timer setup to nil out contexts internal store.
	Cancel()

	// IsExpired returns true/false if the context is considered expired.
	IsExpired() bool

	// Get returns the giving value for the provided key if it exists else nil.
	Get(key interface{}) interface{}

	// WhenExpired returns a channel which gets closed when the given channel
	// expires else closes immediately if its not an expiring context.
	WhenExpired() <-chan struct{}

	// WithTimeout returns a new Context from the previous with the given timeout
	// if the timeout is still further than the previous in expiration date else uses
	// the previous expiration date instead since that is still further in the future.
	WithTimeout(timeout time.Duration) Context

	// Set adds a key and value pair into the context store.
	Set(key interface{}, value interface{})

	// WithValue returns a new context then adds the key and value pair into the
	// context's store.
	WithValue(key interface{}, value interface{}) Context

	// TimeRemaining returns the remaining time for expiring of the context if it
	// indeed has an expiration date set and returns a bool value indicating if it
	// has a timeout.
	TimeRemaining() (remaining time.Duration, hasTimeout bool)
}

// New returns a new context instance.
func New() Context {
	cl := context{
		fields: nilPair,
	}
	return &cl
}

//==============================================================================

// context defines a struct for bundling a context against specific
// use cases with a explicitly set duration which clears all its internal
// data after the giving period.
type context struct {
	fields   *Pair
	lifetime time.Time
	timer    *time.Timer
	duration time.Duration
}

// WithTimeout returns a new context whoes internal value expires
// after the giving duration.
func (c *context) WithTimeout(life time.Duration) Context {
	child := c.newChild()

	lifetime := time.Now().Add(life)

	var useChild bool

	if lifetime.After(child.lifetime) {
		child.duration = life
		child.lifetime = lifetime
		useChild = true
	}

	var to time.Duration

	if useChild {
		to = life
	} else {
		to = c.duration
	}

	child.timer = time.AfterFunc(to, func() {
		child.fields = nilPair
	})

	return child
}

// WithValue returns a new context based on the previos one.
func (c *context) WithValue(key, value interface{}) Context {
	child := c.newChild()
	child.fields = Append(child.fields, key, value)
	return child
}

// TimeRemaining returns the remaining time before expiration.
func (c *context) TimeRemaining() (rem time.Duration, hasTimeout bool) {
	if c.lifetime.IsZero() {
		return
	}

	now := time.Now()

	hasTimeout = true
	if now.Before(c.lifetime) {
		rem = c.lifetime.Sub(now)
		return
	}

	return
}

// WhenExpired returns a channel which gets closed when the context
// has expired.
func (c *context) WhenExpired() <-chan struct{} {
	nc := make(chan struct{})

	if c.IsExpired() {
		close(nc)
		return nc
	}

	rem, _ := c.TimeRemaining()

	go func() {
		defer close(nc)
		<-time.After(rem)
	}()

	return nc
}

// IsExpired returns true/false if the context internal data has expired.
func (c *context) IsExpired() bool {
	left, has := c.TimeRemaining()
	if has {
		if left <= 0 {
			return true
		}
	}

	return false
}

// Cancel cancels the timer if there exists one set to clear context.
func (c *context) Cancel() {
	if c.timer != nil {
		return
	}

	c.timer.Stop()
}

// Set adds the giving value using the given key into the map.
func (c *context) Set(key, val interface{}) {
	c.fields = Append(c.fields, key, val)
}

// Get returns the value for the necessary key within the context.
func (c *context) Get(key interface{}) interface{} {
	val, _ := c.fields.Get(key)
	return val
}

// newChild returns a new context from this one.
func (c *context) newChild() *context {
	return &context{
		fields:   c.fields,
		lifetime: c.lifetime,
		duration: c.duration,
	}
}
