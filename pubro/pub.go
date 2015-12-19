package pubro

import (
	"errors"
	"fmt"
	"reflect"
	"runtime"
	"strings"
	"sync"

	"github.com/influx6/faux/pub"
	"github.com/influx6/faux/reflection"
)

// Injector provides a function type matching the output for creating publishers.
// type Injector func(interface{}) pub.Publisher

// Meta provides a registry structure for registering publishers.
type Meta struct {
	Name      string      `json:"name"`
	Desc      string      `json:"desc"`
	Package   string      `json:"package"`
	Inject    interface{} // This is supposed to be a function
	injectArg []reflect.Type
	injectVal reflect.Value
}

// Build creates a new Publisher using the received config value.
func (p Meta) Build(config interface{}) pub.Publisher {
	if p.injectArg == nil {
		args, _ := reflection.GetFuncArgumentsType(p.Inject)
		p.injectArg = args

		tu, _ := reflection.FuncValue(p.Inject)
		p.injectVal = tu
	}

	var vak []reflect.Value

	if config == nil {
		vak = p.injectVal.Call(nil)
	} else {

		ctype := reflect.TypeOf(config)

		if !ctype.AssignableTo(p.injectArg[0]) {
			panic("UnAssignable Value for Inject")
		}

		vak = p.injectVal.Call([]reflect.Value{reflect.ValueOf(config)})
	}

	if len(vak) == 0 {
		panic(fmt.Sprintf("Meta[%s] in Pkg[%s] returns no value", p.Name, p.Package))
	}

	if len(vak) > 1 {
		panic(fmt.Sprintf("Meta[%s] in Pkg[%s] returns values greater than 1", p.Name, p.Package))
	}

	return (vak[0].Interface()).(pub.Publisher)
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

// pubos provides a singleton-registry for publishers.
var pubos struct {
	sync.RWMutex
	pubs map[string]Meta
}

func init() {
	// initialize the pubos map.
	pubos.pubs = make(map[string]Meta)
}

// Register adds a new Publisher constructor into the registery, it will
// panic if there exists a similar registered Publisher with the provided
// Meta.Name.
func Register(meta Meta) {
	if Has(meta.Name) {
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

	pubos.Lock()
	pubos.pubs[meta.Name] = meta
	pubos.Unlock()
}

// Get returns the Meta associated with a name if it exists.
func Get(name string) (Meta, error) {
	var m Meta
	if !Has(name) {
		return m, errors.New("Not Found")
	}

	pubos.RLock()
	m = pubos.pubs[name]
	pubos.RUnlock()

	return m, nil
}

// Has returns true/false if a publisher exists with the giving name.
func Has(name string) bool {
	pubos.RLock()
	defer pubos.RUnlock()

	_, ok := pubos.pubs[name]
	return ok
}

// New returns a new Publisher and applies the config value to the function
// builder, if the publisher injector is found else it will panic.
func New(name string, config interface{}) pub.Publisher {
	if !Has(name) {
		panic(fmt.Sprintf("Pub.Meta for Publisher[%s] does not exists", name))
	}

	var meta Meta
	pubos.RLock()
	meta = pubos.pubs[name]
	pubos.RUnlock()

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
