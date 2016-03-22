package maker

import (
	"fmt"
	"reflect"
	"sync"

	"github.com/influx6/faux/reflection"
)

//==============================================================================

// Transform returns an interface that uses the returned value if any from a
// Meta to construct its desired result.
type Transform interface {
	Transform(result interface{}) interface{}
}

// IdentityTransform provides an indentity transformer which takes and returns
// the results. It implements the transform interface.
type IdentityTransform struct{}

// Transform returns the provided result without any modification.
func (IdentityTransform) Transform(result interface{}) interface{} {
	return result
}

// Identity provides a global handle for the identity transformer.
var Identity IdentityTransform

//==============================================================================

// Makers provides an interface for defining makers.
type Makers interface {
	Register(string, interface{}) error
	Create(string, interface{}) (interface{}, error)
}

// Maker provides a central registery and builder.
type Maker struct {
	transform Transform
	rw        sync.RWMutex
	makers    map[string]*Meta
}

// New returns a new Maker instance.
func New(t Transform) *Maker {
	if t == nil {
		t = Identity
	}

	vs := Maker{
		transform: t,
		makers:    make(map[string]*Meta),
	}

	return &vs
}

// Register registers the giving ViewMaker with the specified name.
// The second argument must be a function which must be called.
func (v *Maker) Register(name string, fx interface{}) error {
	v.rw.RLock()
	_, ok := v.makers[name]
	v.rw.RUnlock()

	if ok {
		return fmt.Errorf("%s already registered", name)
	}

	v.rw.Lock()
	defer v.rw.Unlock()
	v.makers[name] = &Meta{Name: name, Inject: fx}
	return nil
}

// Create builds the item with the giving name.
// Returns an error if the item was not found or failed to build.
func (v *Maker) Create(name string, param interface{}) (interface{}, error) {
	v.rw.RLock()
	mv, ok := v.makers[name]
	v.rw.RUnlock()

	if !ok {
		return nil, fmt.Errorf("%s not found", name)
	}

	vs, err := mv.Build(param)
	if err != nil {
		return nil, err
	}

	vss := v.transform.Transform(vs)

	return vss, nil
}

//==============================================================================

// Meta provides a registry structure for registering building structures.
type Meta struct {
	Name      string      `json:"name"`
	Inject    interface{} // This is supposed to be a function
	injectArg []reflect.Type
	injectVal reflect.Value
}

// Build creates a new Publisher using the received config value.
func (p *Meta) Build(config interface{}) (interface{}, error) {
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
				return nil, fmt.Errorf("Unassignable value for Inject: %+v -> %+v", config, wanted)
			}

			vum := reflect.ValueOf(config)
			configVal = vum.Convert(wanted)
		} else {
			configVal = reflect.ValueOf(config)
		}

		vak = p.injectVal.Call([]reflect.Value{configVal})
	}

	if len(vak) == 0 {
		return nil, fmt.Errorf("Build[%s] returns no value", p.Name)
	}

	if len(vak) > 1 {
		return nil, fmt.Errorf("Build[%s] returns more than one value", p.Name)
	}

	res := vak[0]
	return res.Interface(), nil
}

//==============================================================================
