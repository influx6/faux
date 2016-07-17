package regos

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"runtime"
	"strings"
	"sync"

	"github.com/influx6/faux/reflection"
)

// Do provides a instruction struct for building a publisher from pub.
type Do struct {
	Tag  string      // a custom tag to use for result.
	Name string      // Name of registered publisher builder.
	Use  interface{} //argument to use for Inject
}

// Publisher defines a list of Do instruction for building publishers.
type Publisher []Do

// Publishers returns a new instance of Publishers.
func Publishers() *Publisher {
	return new(Publisher)
}

// MustAdd adds a new Do instruction into the list. If it fails to add it panics.
func (d *Publisher) MustAdd(do Do) {
	if err := d.Add(do); err != nil {
		panic(err)
	}
}

// Add adds a new Do instruction into the list.
func (d *Publisher) Add(do Do) error {
	if do.Tag == "" {
		return errors.New("Tag can not be empty")
	}

	if do.Name == "" {
		return errors.New("Name can not be empty")
	}

	*d = append(*d, do)
	return nil
}

//==============================================================================

// Meta provides a registry structure for registering building structures.
type Meta struct {
	Name      string      `json:"name"`
	Desc      string      `json:"desc"`
	Package   string      `json:"package"`
	Inject    interface{} // This is supposed to be a function
	injectArg []reflect.Type
	injectVal reflect.Value
}

// Build creates a new Publisher using the received config value.
func (p Meta) Build(config interface{}) interface{} {
	if p.injectArg == nil {
		args, _ := reflection.GetFuncArgumentsType(p.Inject)
		p.injectArg = args

		tu, _ := reflection.FuncValue(p.Inject)
		p.injectVal = tu
	}

	var vak []reflect.Value
	var configVal reflect.Value

	if config == nil || len(p.injectArg) == 0 {
		vak = p.injectVal.Call(nil)
	} else {

		wanted := p.injectArg[0]
		ctype := reflect.TypeOf(config)

		if !ctype.AssignableTo(wanted) {
			if !ctype.ConvertibleTo(wanted) {
				panic(fmt.Sprintf("Unassignable value for Inject: %s -> %+s", query(config), wanted))
			}

			vum := reflect.ValueOf(config)
			configVal = vum.Convert(wanted)
		} else {
			configVal = reflect.ValueOf(config)
		}

		vak = p.injectVal.Call([]reflect.Value{configVal})
	}

	if len(vak) == 0 {
		panic(fmt.Sprintf("Meta[%s] in Pkg[%s] returns no value", p.Name, p.Package))
	}

	if len(vak) > 1 {
		panic(fmt.Sprintf("Meta[%s] in Pkg[%s] returns values greater than 1", p.Name, p.Package))
	}

	pubMade := vak[0]

	return (pubMade.Interface())
}

// Validate ensures the Meta provides the necessary information needed
// to be a valid pub meta.
func (p Meta) Validate() error {
	if p.Name == "" {
		return errors.New("Empty Name")
	}

	if p.Desc == "" {
		return errors.New("Empty Description")
	}

	if p.Package == "" {
		return errors.New("Empty Package")
	}

	if p.Inject == nil {
		return errors.New("No Injector")
	}

	if !reflection.IsFuncType(p.Inject) {
		return errors.New("Inject is not a Function")
	}

	// Must have either zero or one argument for functions.
	if !reflection.HasArgumentSize(p.Inject, 0) && !reflection.HasArgumentSize(p.Inject, 1) {
		return errors.New("Argument Size Greater Than 1")
	}

	return nil
}

// Regus provides a singleton-registry for publishers.
type Regus struct {
	sync.RWMutex
	pubs map[string]Meta
}

// New returns a new instance of a Regus.
func New() *Regus {
	var regus Regus
	regus.pubs = make(map[string]Meta)
	return &regus
}

// Register adds a new Publisher constructor into the registery, it will
// panic if there exists a similar registered buildable structures with
// the provided Meta.Name.
func (regus *Regus) Register(meta Meta) {
	if regus.Has(meta.Name) {
		panic(fmt.Sprintf("Name[%s] already assigned", meta.Name))
	}

	if meta.Package == "" {
		tm := reflect.TypeOf(meta.Inject)
		if tm.Kind() == reflect.Ptr {
			tm = tm.Elem()
		}

		pc, _, _, _ := runtime.Caller(2)
		pkg, pkgName := splitPath(runtime.FuncForPC(pc).Name())
		parts := strings.Split(pkgName, ".")
		plen := len(parts)

		if plen > 1 {
			if pkg == "" {
				pkg = parts[0]
			} else {
				pkg = fmt.Sprintf("%s/%s", pkg, parts[0])
			}
		}

		meta.Package = pkg
	}

	if err := meta.Validate(); err != nil {
		panic(fmt.Sprintf("Meta[%s] is Invalid: %s", meta.Name, err))
	}

	regus.Lock()
	regus.pubs[meta.Name] = meta
	regus.Unlock()
}

// Get returns the Meta associated with a name if it exists.
func (regus *Regus) Get(name string) (Meta, error) {
	var m Meta
	if !regus.Has(name) {
		return m, errors.New("Not Found")
	}

	regus.RLock()
	m = regus.pubs[name]
	regus.RUnlock()

	return m, nil
}

// Has returns true/false if a publisher exists with the giving name.
func (regus *Regus) Has(name string) bool {
	regus.RLock()
	defer regus.RUnlock()

	_, ok := regus.pubs[name]
	return ok
}

// NewBuild returns a new Publisher and applies the config value to the function
// builder, if the publisher injector is found else it will panic.
func (regus *Regus) NewBuild(name string, config interface{}) interface{} {
	if !regus.Has(name) {
		panic(fmt.Sprintf("Pub.Meta for Publisher[%s] does not exists", name))
	}

	var meta Meta
	regus.RLock()
	meta = regus.pubs[name]
	regus.RUnlock()

	return meta.Build(config)
}

// splitAndLastSlash splits a string by a formward slash and returns the left
// and right parts of it.
func splitPath(line string) (string, string) {
	parts := strings.Split(line, "/")
	partsLen := len(parts)
	right := parts[partsLen-1]
	left := strings.Join(parts[:partsLen-1], "/")
	return left, right
}

// query provides a string version of the value.
func query(value interface{}) string {
	json, err := json.Marshal(value)
	if err != nil {
		return ""
	}

	return string(json)
}
