package bag

import (
	"context"
	"sync"
	"time"
)

// Fields defines a map of key:value pairs.
type Fields map[interface{}]interface{}

// Getter defines a series of Get methods for which values will be retrieved with.
type Getter interface {
	GetInt(interface{}) int
	GetBool(interface{}) bool
	GetInt8(interface{}) int8
	GetInt16(interface{}) int16
	GetInt32(interface{}) int32
	GetInt64(interface{}) int64
	Get(interface{}) interface{}
	GetString(interface{}) string
	GetFloat32(interface{}) float32
	GetFloat64(interface{}) float64
	GetDuration(interface{}) time.Duration
}

// ValueBag defines a context for holding values to be shared across processes..
type ValueBag interface {
	Getter

	// Set adds a key-value pair into the bag.
	Set(key, value interface{})

	// WithValue returns a new context then adds the key and value pair into the
	// context's store.
	WithValue(key interface{}, value interface{}) ValueBag
}

// vbag defines a struct for bundling a context against specific
// use cases with a explicitly set duration which clears all its internal
// data after the giving period.
type vbag struct {
	ml     sync.RWMutex
	fields map[interface{}]interface{}
}

// ValueBagFrom adds giving key-value pairs into the bag.
func ValueBagFrom(fields map[interface{}]interface{}) ValueBag {
	return &vbag{fields: fields}
}

// NewValueBag returns a new context object that meets the Context interface.
func NewValueBag() ValueBag {
	return &vbag{
		fields: map[interface{}]interface{}{},
	}
}

// Set adds given value into context.
func (c *vbag) Set(key, value interface{}) {
	c.ml.Lock()
	defer c.ml.Unlock()
	c.fields[key] = value
}

// WithValue returns a new context based on the previos one.
func (c *vbag) WithValue(key, value interface{}) ValueBag {
	c.ml.RLock()
	defer c.ml.RUnlock()

	fields := make(map[interface{}]interface{})
	for k, v := range c.fields {
		fields[k] = v
	}

	fields[key] = value
	return ValueBagFrom(fields)
}

// Deadline returns giving time when context is expected to be canceled.
func (c *vbag) Deadline() time.Time {
	return time.Time{}
}

// GetDuration returns the duration value of a key if it exists.
func (c *vbag) GetDuration(key interface{}) time.Duration {
	val, found := c.get(key)
	if !found {
		return 0
	}

	if dval, ok := val.(time.Duration); ok {
		return dval
	}

	if dval, ok := val.(int64); ok {
		return time.Duration(dval)
	}

	if sval, ok := val.(string); ok {
		if dur, err := time.ParseDuration(sval); err == nil {
			return dur
		}
	}

	return 0
}

// GetBool returns the bool value of a key if it exists.
func (c *vbag) GetBool(key interface{}) bool {
	val, found := c.get(key)
	if !found {
		return false
	}

	return val.(bool)

}

// GetFloat64 returns the float64 value of a key if it exists.
func (c *vbag) GetFloat64(key interface{}) float64 {
	val, found := c.get(key)
	if !found {
		return 0
	}

	return val.(float64)

}

// GetFloat32 returns the float32 value of a key if it exists.
func (c *vbag) GetFloat32(key interface{}) float32 {
	val, found := c.get(key)
	if !found {
		return 0
	}

	return val.(float32)

}

// GetInt8 returns the int8 value of a key if it exists.
func (c *vbag) GetInt8(key interface{}) int8 {
	val, found := c.get(key)
	if !found {
		return 0
	}

	return val.(int8)

}

// GetInt16 returns the int16 value of a key if it exists.
func (c *vbag) GetInt16(key interface{}) int16 {
	val, found := c.get(key)
	if !found {
		return 0
	}

	return val.(int16)

}

// GetInt64 returns the value type value of a key if it exists.
func (c *vbag) GetInt64(key interface{}) int64 {
	val, found := c.get(key)
	if !found {
		return 0
	}

	return val.(int64)

}

// GetInt32 returns the value type value of a key if it exists.
func (c *vbag) GetInt32(key interface{}) int32 {
	val, found := c.get(key)
	if !found {
		return 0
	}

	return val.(int32)

}

// GetInt returns the value type value of a key if it exists.
func (c *vbag) GetInt(key interface{}) int {
	val, found := c.get(key)
	if !found {
		return 0
	}

	return val.(int)

}

// GetString returns the value type value of a key if it exists.
func (c *vbag) GetString(key interface{}) string {
	val, found := c.get(key)
	if !found {
		return ""
	}

	return val.(string)

}

// Get returns the value of a key if it exists.
func (c *vbag) Get(key interface{}) (value interface{}) {
	item, _ := c.get(key)
	return item
}

// Get returns the value of a key if it exists.
func (c *vbag) get(key interface{}) (value interface{}, found bool) {
	c.ml.RLock()
	defer c.ml.RUnlock()

	item, ok := c.fields[key]
	return item, ok
}

//==============================================================================

// googleContext implements a decorator for googles context package.
type googleContext struct {
	context.Context
}

// FromContext returns a new context object that meets the Context interface.
func FromContext(ctx context.Context) *googleContext {
	var gc googleContext
	gc.Context = ctx
	return &gc
}

// GetDuration returns the giving value for the provided key if it exists else nil.
func (g *googleContext) GetDuration(key interface{}) time.Duration {
	val := g.Context.Value(key)
	if val == nil {
		return 0
	}

	if dval, ok := val.(time.Duration); ok {
		return dval
	}

	if dval, ok := val.(int64); ok {
		return time.Duration(dval)
	}

	if sval, ok := val.(string); ok {
		if dur, err := time.ParseDuration(sval); err == nil {
			return dur
		}
	}

	return 0
}

// Get returns the giving value for the provided key if it exists else nil.
func (g *googleContext) Get(key interface{}) interface{} {
	val := g.Context.Value(key)
	if val == nil {
		return val
	}

	return val
}

// GetBool returns the value type value of a key if it exists.
func (g *googleContext) GetBool(key interface{}) bool {
	return g.Get(key).(bool)
}

// GetFloat64 returns the value type value of a key if it exists.
func (g *googleContext) GetFloat64(key interface{}) float64 {
	return g.Get(key).(float64)

}

// GetFloat32 returns the value type value of a key if it exists.
func (g *googleContext) GetFloat32(key interface{}) float32 {
	return g.Get(key).(float32)
}

// GetInt8 returns the value type value of a key if it exists.
func (g *googleContext) GetInt8(key interface{}) int8 {
	return g.Get(key).(int8)
}

// GetInt16 returns the value type value of a key if it exists.
func (g *googleContext) GetInt16(key interface{}) int16 {
	return g.Get(key).(int16)
}

// GetInt64 returns the value type value of a key if it exists.
func (g *googleContext) GetInt64(key interface{}) int64 {
	return g.Get(key).(int64)
}

// GetInt32 returns the value type value of a key if it exists.
func (g *googleContext) GetInt32(key interface{}) int32 {
	return g.Get(key).(int32)
}

// GetInt returns the value type value of a key if it exists.
func (g *googleContext) GetInt(key interface{}) int {
	return g.Get(key).(int)
}

// GetString returns the value type value of a key if it exists.
func (g *googleContext) GetString(key interface{}) string {
	return g.Get(key).(string)
}

//==============================================================================

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

// RemoveAll sets all key-value pairs to nil for all connected pair, till it reaches
// the root.
func (p *Pair) RemoveAll() {
	p.key = nil
	p.value = nil

	if p.prev != nil {
		p.prev.RemoveAll()
	}
}

// Root returns the root Pair in the chain which links all pairs together.
func (p *Pair) Root() *Pair {
	if p.prev == nil {
		return p
	}

	return p.prev.Root()
}

// GetDuration returns the duration value of a key if it exists.
func (p *Pair) GetDuration(key interface{}) time.Duration {
	val, found := p.Get(key)
	if !found {
		return 0
	}

	if dval, ok := val.(time.Duration); ok {
		return dval
	}

	if dval, ok := val.(int64); ok {
		return time.Duration(dval)
	}

	if sval, ok := val.(string); ok {
		if dur, err := time.ParseDuration(sval); err == nil {
			return dur
		}
	}

	return 0
}

// GetBool returns the bool value of a key if it exists.
func (p *Pair) GetBool(key interface{}) bool {
	val, found := p.Get(key)
	if !found {
		return false
	}

	return val.(bool)

}

// GetFloat64 returns the float64 value of a key if it exists.
func (p *Pair) GetFloat64(key interface{}) float64 {
	val, found := p.Get(key)
	if !found {
		return 0
	}

	return val.(float64)

}

// GetFloat32 returns the float32 value of a key if it exists.
func (p *Pair) GetFloat32(key interface{}) float32 {
	val, found := p.Get(key)
	if !found {
		return 0
	}

	return val.(float32)

}

// GetInt8 returns the int8 value of a key if it exists.
func (p *Pair) GetInt8(key interface{}) int8 {
	val, found := p.Get(key)
	if !found {
		return 0
	}

	return val.(int8)

}

// GetInt16 returns the int16 value of a key if it exists.
func (p *Pair) GetInt16(key interface{}) int16 {
	val, found := p.Get(key)
	if !found {
		return 0
	}

	return val.(int16)

}

// GetInt64 returns the value type value of a key if it exists.
func (p *Pair) GetInt64(key interface{}) int64 {
	val, found := p.Get(key)
	if !found {
		return 0
	}

	return val.(int64)

}

// GetInt32 returns the value type value of a key if it exists.
func (p *Pair) GetInt32(key interface{}) int32 {
	val, found := p.Get(key)
	if !found {
		return 0
	}

	return val.(int32)

}

// GetInt returns the value type value of a key if it exists.
func (p *Pair) GetInt(key interface{}) int {
	val, found := p.Get(key)
	if !found {
		return 0
	}

	return val.(int)

}

// GetString returns the value type value of a key if it exists.
func (p *Pair) GetString(key interface{}) string {
	val, found := p.Get(key)
	if !found {
		return ""
	}

	return val.(string)

}

// Get returns the value of a key if it exists.
func (p *Pair) Get(key interface{}) (interface{}, bool) {
	if p == nil {
		return nil, false
	}

	if p.key == key {
		return p.value, true
	}

	if p.prev == nil {
		return nil, false
	}

	return p.prev.Get(key)
}
