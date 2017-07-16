// Package context is built out of my desire to understand the http context
// library and as an experiement in such a library works.
package context

import (
	gcontext "context"
	"errors"
	"sync"
	"time"
)

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

// NewPair returns a a key-value pair chain for setting fields.
func NewPair(key, value interface{}) *Pair {
	return &Pair{
		key:   key,
		value: value,
	}
}

// Append returns a new Pair with the giving key and with the provded Pair set as
// it's previous link.
func Append(p *Pair, key, value interface{}) *Pair {
	return p.Append(key, value)
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

	if p.key != "" {
		f[p.key] = p.value
	}

	return f
}

// Append returns a new pair with the giving key and value and its previous
// set to this pair.
func (p *Pair) Append(key, val interface{}) *Pair {
	return &Pair{
		prev:  p,
		key:   key,
		value: val,
	}
}

// Root returns the root Pair in the chain which links all pairs together.
func (p *Pair) Root() *Pair {
	if p.prev == nil {
		return p
	}

	return p.prev.Root()
}

// GetBool collects the string value of a key if it exists.
func (p *Pair) GetBool(key interface{}) (bool, bool) {
	val, found := p.Get(key)
	if !found {
		return false, false
	}

	value, ok := val.(bool)
	return value, ok
}

// GetFloat64 collects the string value of a key if it exists.
func (p *Pair) GetFloat64(key interface{}) (float64, bool) {
	val, found := p.Get(key)
	if !found {
		return 0, false
	}

	value, ok := val.(float64)
	return value, ok
}

// GetFloat32 collects the string value of a key if it exists.
func (p *Pair) GetFloat32(key interface{}) (float32, bool) {
	val, found := p.Get(key)
	if !found {
		return 0, false
	}

	value, ok := val.(float32)
	return value, ok
}

// GetInt8 collects the string value of a key if it exists.
func (p *Pair) GetInt8(key interface{}) (int8, bool) {
	val, found := p.Get(key)
	if !found {
		return 0, false
	}

	value, ok := val.(int8)
	return value, ok
}

// GetInt16 collects the string value of a key if it exists.
func (p *Pair) GetInt16(key interface{}) (int16, bool) {
	val, found := p.Get(key)
	if !found {
		return 0, false
	}

	value, ok := val.(int16)
	return value, ok
}

// GetInt64 collects the string value of a key if it exists.
func (p *Pair) GetInt64(key interface{}) (int64, bool) {
	val, found := p.Get(key)
	if !found {
		return 0, false
	}

	value, ok := val.(int64)
	return value, ok
}

// GetInt32 collects the string value of a key if it exists.
func (p *Pair) GetInt32(key interface{}) (int32, bool) {
	val, found := p.Get(key)
	if !found {
		return 0, false
	}

	value, ok := val.(int32)
	return value, ok
}

// GetInt collects the string value of a key if it exists.
func (p *Pair) GetInt(key interface{}) (int, bool) {
	val, found := p.Get(key)
	if !found {
		return 0, false
	}

	value, ok := val.(int)
	return value, ok
}

// GetString collects the string value of a key if it exists.
func (p *Pair) GetString(key interface{}) (string, bool) {
	val, found := p.Get(key)
	if !found {
		return "", false
	}

	value, ok := val.(string)
	return value, ok
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

// Canceler defines an interface for canceling an operation with a giving error.
type Canceler interface {

	// Cancel is called to cancel with a giving error.
	Cancel(error)

	// Err returns the error partaining to error if any giving when Cancelled is called.
	Err() error
}

// Context defines an interface for a context providers which allows us to
// build passable context around.
type Context interface {

	// IsExpired returns true/false if the context is considered expired.
	IsExpired() bool

	// Series of Gets returning value for the provided key if it exists else the type default value.
	Get(key interface{}) (interface{}, bool)
	GetInt(key interface{}) (int, bool)
	GetBool(key interface{}) (bool, bool)
	GetInt8(key interface{}) (int8, bool)
	GetInt16(key interface{}) (int16, bool)
	GetInt32(key interface{}) (int32, bool)
	GetInt64(key interface{}) (int64, bool)
	GetString(key interface{}) (string, bool)
	GetFloat32(key interface{}) (float32, bool)
	GetFloat64(key interface{}) (float64, bool)

	// Done returns a channel which gets closed when the given channel
	// expires else closes immediately if its not an expiring context.
	Done() <-chan struct{}

	// WithDeadline returns a new Context from the previous with the given timeout
	// if the timeout is still further than the previous in expiration date else uses
	// the previous expiration date instead since that is still further in the future.
	WithDeadline(timeout time.Duration, cancelWithParent bool) Context

	// New returns a new context based on the fileds of the context which its
	// called from, it does inherits the lifetime limits of the context its
	// called from.
	New(cancelWithParent bool) Context

	// Set adds a key and value pair into the context store.
	Set(key interface{}, value interface{})

	// WithValue returns a new context then adds the key and value pair into the
	// context's store.
	WithValue(key interface{}, value interface{}) Context

	// Deadline returns the remaining time for expiring of the context if it
	// indeed has an expiration date set and returns a bool value indicating if it
	// has a timeout.
	Deadline() (time.Duration, bool)
}

// CancelableContext defines the base outline for all context.
type CancelableContext interface {
	Context
	Canceler

	// Ctx returns a Context which exposes a basic context interface without  the
	// cancellable method.
	Ctx() Context
}

// New returns a new context object that meets the Context interface.
func New() CancelableContext {
	cl := context{
		fields:    nilPair,
		canceller: make(chan struct{}),
	}

	return &cl
}

// From returns a new context object that meets the Context interface.
func From(ctx gcontext.Context) CancelableContext {
	var gc googleContext

	ctx, canceller := gcontext.WithCancel(ctx)
	rem, _ := ctx.Deadline()

	gc.ctx = ctx
	gc.deadline = rem
	gc.canceller = canceller

	return &gc
}

//==============================================================================

// googleContext implements a decorator for googles context package.
type googleContext struct {
	ctx       gcontext.Context
	deadline  time.Time
	canceller func()
	err       error
}

// IsExpired returns true/false if the context is considered expired.
func (g *googleContext) IsExpired() bool {
	select {
	case <-g.ctx.Done():
		return true
	case <-time.After(1 * time.Second):
		return false
	}
}

// Get returns the giving value for the provided key if it exists else nil.
func (g *googleContext) Get(key interface{}) (interface{}, bool) {
	val := g.ctx.Value(key)
	if val == nil {
		return val, false
	}

	return val, true
}

// GetBool collects the string value of a key if it exists.
func (g *googleContext) GetBool(key interface{}) (bool, bool) {
	val, found := g.Get(key)
	if !found {
		return false, false
	}

	value, ok := val.(bool)
	return value, ok
}

// GetFloat64 collects the string value of a key if it exists.
func (g *googleContext) GetFloat64(key interface{}) (float64, bool) {
	val, found := g.Get(key)
	if !found {
		return 0, false
	}

	value, ok := val.(float64)
	return value, ok
}

// GetFloat32 collects the string value of a key if it exists.
func (g *googleContext) GetFloat32(key interface{}) (float32, bool) {
	val, found := g.Get(key)
	if !found {
		return 0, false
	}

	value, ok := val.(float32)
	return value, ok
}

// GetInt8 collects the string value of a key if it exists.
func (g *googleContext) GetInt8(key interface{}) (int8, bool) {
	val, found := g.Get(key)
	if !found {
		return 0, false
	}

	value, ok := val.(int8)
	return value, ok
}

// GetInt16 collects the string value of a key if it exists.
func (g *googleContext) GetInt16(key interface{}) (int16, bool) {
	val, found := g.Get(key)
	if !found {
		return 0, false
	}

	value, ok := val.(int16)
	return value, ok
}

// GetInt64 collects the string value of a key if it exists.
func (g *googleContext) GetInt64(key interface{}) (int64, bool) {
	val, found := g.Get(key)
	if !found {
		return 0, false
	}

	value, ok := val.(int64)
	return value, ok
}

// GetInt32 collects the string value of a key if it exists.
func (g *googleContext) GetInt32(key interface{}) (int32, bool) {
	val, found := g.Get(key)
	if !found {
		return 0, false
	}

	value, ok := val.(int32)
	return value, ok
}

// GetInt collects the string value of a key if it exists.
func (g *googleContext) GetInt(key interface{}) (int, bool) {
	val, found := g.Get(key)
	if !found {
		return 0, false
	}

	value, ok := val.(int)
	return value, ok
}

// GetString collects the string value of a key if it exists.
func (g *googleContext) GetString(key interface{}) (string, bool) {
	val, found := g.Get(key)
	if !found {
		return "", false
	}

	value, ok := val.(string)
	return value, ok
}

// Done returns a channel which gets closed when the given channel
// expires else closes immediately if its not an expiring context.
func (g *googleContext) Done() <-chan struct{} {
	return g.ctx.Done()
}

// Err returns the error pertaining to the context Err() method.
func (g *googleContext) Err() error {
	return g.ctx.Err()
}

// Ctx returns a Context which exposes a basic context interface without  the
// cancellable method.
func (g *googleContext) Ctx() Context {
	return g
}

// New returns a new context based on the fileds of the context which its
// called from, it does inherits the lifetime limits of the context its
// called from.
func (g *googleContext) New(cancelWithParent bool) Context {
	return From(g.ctx)
}

// WithDeadline returns a new Context from the previous with the given timeout
// if the timeout is still further than the previous in expiration date else uses
// the previous expiration date instead since that is still further in the future.
func (g *googleContext) WithDeadline(timeout time.Duration, cancelWithParent bool) Context {
	ctx, cancller := gcontext.WithTimeout(g.ctx, timeout)

	var gc googleContext
	gc.ctx = ctx
	gc.canceller = cancller

	return &gc
}

// Set adds a key and value pair into the context store.
func (g *googleContext) Set(key interface{}, value interface{}) {
	ctx := gcontext.WithValue(g.ctx, key, value)
	g.ctx = ctx
}

// WithValue returns a new context then adds the key and value pair into the
// context's store.
func (g *googleContext) WithValue(key interface{}, value interface{}) Context {
	ctx := gcontext.WithValue(g.ctx, key, value)

	nctx, cancel := gcontext.WithCancel(ctx)

	var gc googleContext
	gc.ctx = nctx
	gc.canceller = cancel

	return &gc
}

// Deadline returns the remaining time for expiring of the context if it
// indeed has an expiration date set and returns a bool value indicating if it
// has a timeout.
func (g *googleContext) Deadline() (remaining time.Duration, hasTimeout bool) {
	deadline, ok := g.ctx.Deadline()

	return time.Now().Sub(deadline), ok
}

// Cancel cancels the timer if there exists one set to clear context.
func (g *googleContext) Cancel(err error) {
	g.canceller()
}

//================================================================================

// context defines a struct for bundling a context against specific
// use cases with a explicitly set duration which clears all its internal
// data after the giving period.
type context struct {
	fields    *Pair
	lifetime  time.Time
	timer     *time.Timer
	duration  time.Duration
	parent    Context
	canceller chan struct{}
	cl        sync.Mutex
	err       error
	canceled  bool
}

// New returns a new context from with the configuration limits of this one.
func (c *context) New(cancelWithParent bool) Context {
	if c.timer != nil {
		return c.WithDeadline(c.duration, cancelWithParent)
	}

	return c.newChild(cancelWithParent)
}

// WithDeadline returns a new context whoes internal value expires
// after the giving duration.
func (c *context) WithDeadline(life time.Duration, cancelWithParent bool) Context {
	child := c.newChild(cancelWithParent)

	var useChild bool

	lifetime := time.Now().Add(life)
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
		child.Cancel(errors.New("Deadline passed"))
	})

	return child
}

// WithValue returns a new context based on the previos one.
func (c *context) WithValue(key, value interface{}) Context {
	child := c.newChild(true)
	child.fields = Append(child.fields, key, value)
	return child
}

// Deadline returns the remaining time before expiration.
func (c *context) Deadline() (rem time.Duration, hasTimeout bool) {
	if c.lifetime.IsZero() {
		return
	}

	hasTimeout = true

	now := time.Now()
	if now.Before(c.lifetime) {
		rem = c.lifetime.Sub(now)
		return
	}

	return
}

// Done returns a channel which gets closed when the context
// has expired.
func (c *context) Done() <-chan struct{} {
	if c.IsExpired() {
		cm := make(chan struct{})

		close(cm)

		return cm
	}

	return c.canceller
}

// IsExpired returns true/false if the context internal data has expired.
func (c *context) IsExpired() bool {
	left, has := c.Deadline()
	if has {
		if left <= 0 {
			return true
		}
	}

	c.cl.Lock()
	{
		if c.canceled {
			c.cl.Unlock()
			return true
		}
	}
	c.cl.Unlock()

	return false
}

// Ctx returns the Context interface  for a giving Context.
func (c *context) Ctx() Context {
	return c
}

// Err returns the error pertaining to the context Err() method.
func (c *context) Err() error {
	c.cl.Lock()
	defer c.cl.Unlock()

	return c.err
}

// Cancel cancels the timer if there exists one set to clear context.
func (c *context) Cancel(err error) {
	if c.IsExpired() {
		return
	}

	c.cl.Lock()
	c.err = err
	c.canceled = true
	c.cl.Unlock()

	close(c.canceller)

	if c.timer != nil {
		c.timer.Stop()
		return
	}
}

// Set adds the giving value using the given key into the map.
func (c *context) Set(key, val interface{}) {
	c.fields = Append(c.fields, key, val)
}

// Get returns the value for the necessary key within the context.
func (c *context) Get(key interface{}) (item interface{}, found bool) {
	item, found = c.fields.Get(key)
	return
}

// GetBool collects the string value of a key if it exists.
func (c *context) GetBool(key interface{}) (bool, bool) {
	return c.fields.GetBool(key)
}

// GetFloat64 collects the string value of a key if it exists.
func (c *context) GetFloat64(key interface{}) (float64, bool) {
	return c.fields.GetFloat64(key)
}

// GetFloat32 collects the string value of a key if it exists.
func (c *context) GetFloat32(key interface{}) (float32, bool) {
	return c.fields.GetFloat32(key)
}

// GetInt8 collects the string value of a key if it exists.
func (c *context) GetInt8(key interface{}) (int8, bool) {
	return c.fields.GetInt8(key)
}

// GetInt16 collects the string value of a key if it exists.
func (c *context) GetInt16(key interface{}) (int16, bool) {
	return c.fields.GetInt16(key)
}

// GetInt64 collects the string value of a key if it exists.
func (c *context) GetInt64(key interface{}) (int64, bool) {
	return c.fields.GetInt64(key)
}

// GetInt32 collects the string value of a key if it exists.
func (c *context) GetInt32(key interface{}) (int32, bool) {
	return c.fields.GetInt32(key)
}

// GetInt collects the string value of a key if it exists.
func (c *context) GetInt(key interface{}) (int, bool) {
	return c.fields.GetInt(key)
}

// GetString collects the string value of a key if it exists.
func (c *context) GetString(key interface{}) (string, bool) {
	return c.fields.GetString(key)
}

// newChild returns a new fresh context based on the fields of this context.
func (c *context) newChild(cancelWithParent bool) *context {
	canceller := make(chan struct{})

	if c.IsExpired() {
		close(canceller)
	}

	cm := &context{
		parent:    c,
		fields:    c.fields,
		lifetime:  c.lifetime,
		duration:  c.duration,
		canceled:  c.canceled,
		canceller: canceller,
	}

	if cancelWithParent {
		go func() {
			cancel := c.Done()

			<-cancel
			cm.Cancel(errors.New("Already canceled"))
		}()
	}

	return cm
}
