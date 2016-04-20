package flux

import (
	"fmt"
	"sync"
)

//Eachfunc defines the type of the Mappable.Each rule
type Eachfunc func(interface{}, interface{}, func())

//StringEachfunc defines the type of the Mappable.Each rule
type StringEachfunc func(interface{}, string, func())

// Maps define a set of method rules for maps of the string key types
type Maps interface {
	Clear()
	HasMatch(k string, v interface{}) bool
	Each(f StringEachfunc)
	Keys() []string
	Copy(map[string]interface{})
	Has(string) bool
	Get(string) interface{}
	Remove(string)
	Set(k string, v interface{})
}

//StringMappable defines member function rules for securemap
type StringMappable interface {
	Maps
	Clone() StringMappable
}

// Collectors defines member function rules for collector
type Collectors interface {
	Maps
	Clone() Collector
}

//Collector defines a typ of map string
type Collector map[string]interface{}

//NewCollector returns a new collector instance
func NewCollector() Collector {
	return make(Collector)
}

//Clone makes a new clone of this collector
func (c Collector) Clone() Collector {
	col := make(Collector)
	col.Copy(c)
	return col
}

//Remove deletes a key:value pair
func (c Collector) Remove(k string) {
	if c.Has(k) {
		delete(c, k)
	}
}

//Keys return the keys of the Collector
func (c Collector) Keys() []string {
	var keys []string
	c.Each(func(_ interface{}, k string, _ func()) {
		keys = append(keys, k)
	})
	return keys
}

//Get returns the value with the key
func (c Collector) Get(k string) interface{} {
	return c[k]
}

//Has returns if a key exists
func (c Collector) Has(k string) bool {
	_, ok := c[k]
	return ok
}

//HasMatch checks if key and value exists and are matching
func (c Collector) HasMatch(k string, v interface{}) bool {
	if c.Has(k) {
		return c.Get(k) == v
	}
	return false
}

//Set puts a specific key:value into the collector
func (c Collector) Set(k string, v interface{}) {
	c[k] = v
}

//Copy copies the map into the collector
func (c Collector) Copy(m map[string]interface{}) {
	for v, k := range m {
		c.Set(v, k)
	}
}

//Each iterates through all items in the collector
func (c Collector) Each(fx StringEachfunc) {
	var state bool
	for k, v := range c {
		if state {
			break
		}

		fx(v, k, func() {
			state = true
		})
	}
}

//Clear clears the collector
func (c Collector) Clear() {
	for k := range c {
		delete(c, k)
	}
}

// SyncCollectors defines member function rules for SyncCollector
type SyncCollectors interface {
	Maps
	Clone() SyncCollectors
}

// SyncCollector provides a mutex controlled map
type SyncCollector struct {
	c  Collector
	rw sync.RWMutex
}

//NewSyncCollector returns a new collector instance
func NewSyncCollector() *SyncCollector {
	so := SyncCollector{c: make(Collector)}
	return &so
}

//Clone makes a new clone of this collector
func (c *SyncCollector) Clone() SyncCollectors {
	var co Collector

	c.rw.RLock()
	co = c.c.Clone()
	c.rw.RUnlock()

	so := SyncCollector{c: co}
	return &so
}

//Remove deletes a key:value pair
func (c *SyncCollector) Remove(k string) {
	c.rw.Lock()
	c.c.Remove(k)
	c.rw.Unlock()
}

//Set puts a specific key:value into the collector
func (c *SyncCollector) Set(k string, v interface{}) {
	c.rw.Lock()
	c.c.Set(k, v)
	c.rw.Unlock()
}

//Copy copies the map into the collector
func (c *SyncCollector) Copy(m map[string]interface{}) {
	for v, k := range m {
		c.Set(v, k)
	}
}

//Each iterates through all items in the collector
func (c *SyncCollector) Each(fx StringEachfunc) {
	var state bool
	c.rw.RLock()
	for k, v := range c.c {
		if state {
			break
		}

		fx(v, k, func() {
			state = true
		})
	}
	c.rw.RUnlock()
}

//Keys return the keys of the Collector
func (c *SyncCollector) Keys() []string {
	var keys []string
	c.Each(func(_ interface{}, k string, _ func()) {
		keys = append(keys, k)
	})
	return keys
}

//Get returns the value with the key
func (c *SyncCollector) Get(k string) interface{} {
	var v interface{}
	c.rw.RLock()
	v = c.c.Get(k)
	c.rw.RUnlock()
	return v
}

//Has returns if a key exists
func (c *SyncCollector) Has(k string) bool {
	var ok bool
	c.rw.RLock()
	_, ok = c.c[k]
	c.rw.RUnlock()
	return ok
}

//HasMatch checks if key and value exists and are matching
func (c *SyncCollector) HasMatch(k string, v interface{}) bool {
	// c.rw.RLock()
	// defer c.rw.RUnlock()
	if c.Has(k) {
		return c.Get(k) == v
	}
	return false
}

//Clear clears the collector
func (c *SyncCollector) Clear() {
	for k := range c.c {
		c.rw.Lock()
		delete(c.c, k)
		c.rw.Unlock()
	}
}

//Mappable defines member function rules for securemap
type Mappable interface {
	Clear()
	HasMatch(k, v interface{}) bool
	Each(f Eachfunc)
	Keys() []interface{}
	Copy(map[interface{}]interface{})
	CopySecureMap(Mappable)
	Has(interface{}) bool
	Get(interface{}) interface{}
	Remove(interface{})
	Set(k, v interface{})
	Clone() Mappable
}

//SecureMap simple represents a map with a rwmutex locked in
type SecureMap struct {
	data map[interface{}]interface{}
	lock sync.RWMutex
}

//NewSecureMap returns a new securemap
func NewSecureMap() *SecureMap {
	return &SecureMap{data: make(map[interface{}]interface{})}
}

//SecureMapFrom returns a new securemap
func SecureMapFrom(core map[interface{}]interface{}) *SecureMap {
	return &SecureMap{data: core}
}

//Clear unlinks the previous map
func (m *SecureMap) Clear() {
	m.lock.Lock()
	m.data = make(map[interface{}]interface{})
	m.lock.Unlock()
}

//HasMatch checks if a key exists and if the value matches
func (m *SecureMap) HasMatch(key, value interface{}) bool {
	m.lock.RLock()
	k, ok := m.data[key]
	m.lock.RUnlock()

	if ok {
		return k == value
	}

	return false
}

//Each interates through the map
func (m *SecureMap) Each(fn Eachfunc) {
	stop := false
	m.lock.RLock()
	for k, v := range m.data {
		if stop {
			break
		}

		fn(v, k, func() { stop = true })
	}
	m.lock.RUnlock()
}

//Keys return the keys of the map
func (m *SecureMap) Keys() []interface{} {
	m.lock.RLock()
	keys := make([]interface{}, len(m.data))
	count := 0
	for k := range m.data {
		keys[count] = k
		count++
	}
	m.lock.RUnlock()

	return keys
}

//Clone makes a clone for this securemap
func (m *SecureMap) Clone() Mappable {
	sm := NewSecureMap()
	m.lock.RLock()
	for k, v := range m.data {
		sm.Set(k, v)
	}
	m.lock.RUnlock()
	return sm
}

//CopySecureMap Copies a  into the map
func (m *SecureMap) CopySecureMap(src Mappable) {
	src.Each(func(k, v interface{}, _ func()) {
		m.Set(k, v)
	})
}

//Copy Copies a map[interface{}]interface{} into the map
func (m *SecureMap) Copy(src map[interface{}]interface{}) {
	for k, v := range src {
		m.Set(k, v)
	}
}

//Has returns true/false if value exists by key
func (m *SecureMap) Has(key interface{}) bool {
	m.lock.RLock()
	_, ok := m.data[key]
	m.lock.RUnlock()
	return ok
}

//Get a key's value
func (m *SecureMap) Get(key interface{}) interface{} {
	m.lock.RLock()
	k := m.data[key]
	m.lock.RUnlock()
	return k
}

//Set a key with value
func (m *SecureMap) Set(key, value interface{}) {

	if _, ok := m.data[key]; ok {
		return
	}

	m.lock.Lock()
	m.data[key] = value
	m.lock.Unlock()
}

//Remove a value by its key
func (m *SecureMap) Remove(key interface{}) {
	m.lock.Lock()
	delete(m.data, key)
	m.lock.Unlock()
}

//SecureStack provides addition of functions into a stack
type SecureStack struct {
	listeners []interface{}
	lock      sync.RWMutex
}

//NewSecureStack returns a new concurrent safe array decorator
func NewSecureStack() *SecureStack {
	return &SecureStack{
		listeners: make([]interface{}, 0),
	}
}

//Splice returns a new unique slice from the list
func (f *SecureStack) Splice(begin, end int) []interface{} {
	size := f.Size()

	if end > size {
		end = size
	}

	f.lock.RLock()
	ms := f.listeners[begin:end]
	f.lock.RUnlock()
	var dup []interface{}
	dup = append(dup, ms...)
	return dup
}

//Set lets you retrieve an item in the list
func (f *SecureStack) Set(ind int, d interface{}) {
	sz := f.Size()

	if ind >= sz {
		return
	}

	if ind < 0 {
		ind = sz - ind
		if ind < 0 {
			return
		}
	}

	f.lock.Lock()
	f.listeners[ind] = d
	f.lock.Unlock()
}

//Get lets you retrieve an item in the list
func (f *SecureStack) Get(ind int) interface{} {
	sz := f.Size()
	if ind >= sz {
		ind = sz - 1
	}

	if ind < 0 {
		ind = sz - ind
		if ind < 0 {
			return nil
		}
	}

	f.lock.RLock()
	r := f.listeners[ind]
	f.lock.RUnlock()
	return r
}

//Strings return the stringified version of the internal list
func (f *SecureStack) String() string {
	f.lock.RLock()
	sz := fmt.Sprintf("%+v", f.listeners)
	f.lock.RUnlock()
	return sz
}

//Clear flushes the stack listener
func (f *SecureStack) Clear() {
	f.lock.Lock()
	f.listeners = f.listeners[0:0]
	f.lock.Unlock()
}

//Size returns the total number of listeners
func (f *SecureStack) Size() int {
	f.lock.RLock()
	sz := len(f.listeners)
	f.lock.RUnlock()
	return sz
}

//Add adds a function into the stack
func (f *SecureStack) Add(fx interface{}) int {
	f.lock.Lock()
	ind := len(f.listeners)
	f.listeners = append(f.listeners, fx)
	f.lock.Unlock()
	return ind
}

//Delete removes the function at the provided index
func (f *SecureStack) Delete(ind int) {

	if ind <= 0 && len(f.listeners) <= 0 {
		return
	}

	f.lock.Lock()
	copy(f.listeners[ind:], f.listeners[ind+1:])
	f.lock.Unlock()

	f.lock.RLock()
	f.listeners[len(f.listeners)-1] = nil
	f.lock.RUnlock()

	f.lock.Lock()
	f.listeners = f.listeners[:len(f.listeners)-1]
	f.lock.Unlock()

}

//Each runs through the function lists and executing with args
func (f *SecureStack) Each(fx func(interface{})) {
	if f.Size() <= 0 {
		return
	}

	f.lock.RLock()
	for _, d := range f.listeners {
		fx(d)
	}
	f.lock.RUnlock()
}

//FunctionStack provides addition of functions into a stack
type FunctionStack struct {
	listeners []func(...interface{})
	lock      sync.RWMutex
}

//NewFunctionStack returns a new functionstack instance
func NewFunctionStack() *FunctionStack {
	return &FunctionStack{
		listeners: make([]func(...interface{}), 0),
	}
}

//Clear flushes the stack listener
func (f *FunctionStack) Clear() {
	f.lock.Lock()
	f.listeners = f.listeners[0:0]
	f.lock.Unlock()
}

//Size returns the total number of listeners
func (f *FunctionStack) Size() int {
	f.lock.RLock()
	sz := len(f.listeners)
	f.lock.RUnlock()
	return sz
}

//Add adds a function into the stack
func (f *FunctionStack) Add(fx func(...interface{})) int {
	var ind int

	f.lock.RLock()
	ind = len(f.listeners)
	f.lock.RUnlock()

	f.lock.Lock()
	f.listeners = append(f.listeners, fx)
	f.lock.Unlock()

	return ind
}

//Delete removes the function at the provided index
func (f *FunctionStack) Delete(ind int) {

	if ind <= 0 && len(f.listeners) <= 0 {
		return
	}

	f.lock.Lock()
	copy(f.listeners[ind:], f.listeners[ind+1:])
	f.lock.Unlock()

	f.lock.RLock()
	f.listeners[len(f.listeners)-1] = nil
	f.lock.RUnlock()

	f.lock.Lock()
	f.listeners = f.listeners[:len(f.listeners)-1]
	f.lock.Unlock()

}

//Each runs through the function lists and executing with args
func (f *FunctionStack) Each(d ...interface{}) {
	if f.Size() <= 0 {
		return
	}

	f.lock.RLock()
	for _, fx := range f.listeners {
		fx(d...)
	}
	f.lock.RUnlock()
}

//Splice returns a new unique slice from the list
// func (f *SecureStack) Splice(begin, end int) []interface{} {
// 	size := f.Size()
//
// 	if end > size {
// 		end = size
// 	}
//
// 	f.lock.RLock()
// 	ms := f.listeners[begin:end]
// 	f.lock.RUnlock()
// 	var dup []interface{}
// 	dup = append(dup, ms...)
// 	return dup
// }

//SingleStack provides a function stack fro single argument
//functions
type SingleStack struct {
	*FunctionStack
}

//NewSingleStack returns a singlestack instance
func NewSingleStack() *SingleStack {
	return &SingleStack{
		FunctionStack: NewFunctionStack(),
	}
}

//Add adds a function into the stack
func (s *SingleStack) Add(fx func(interface{})) int {
	return s.FunctionStack.Add(func(f ...interface{}) {
		if len(f) <= 0 {
			fx(nil)
			return
		}

		fx(f[0])
	})
}
