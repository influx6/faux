package ops

import (
	"bytes"
	"encoding/json"
	"errors"

	"github.com/BurntSushi/toml"
	"github.com/influx6/faux/context"
	"github.com/influx6/faux/metrics"
)

// errors.
var (
	ErrFunctionNotFound = errors.New("Function with given id not found")
)

// Op defines an interface which expose an exec method.
type Op interface {
	Exec(context.CancelContext, metrics.Metrics) error
}

// Function defines a interface for a function which returns a giving Op.
type Function func() Op

// Generator defines a interface for a function which returns a giving Op.
type Generator func([]byte) (Op, error)

// GeneratorRegistry implements a structure that handles registration and ochestration of spells.
type GeneratorRegistry struct {
	functions map[string]Generator
}

// NewGeneratorRegistry returns a new instance of the GeneratorsRegistry.
func NewGeneratorRegistry() *GeneratorRegistry {
	return &GeneratorRegistry{
		functions: make(map[string]Generator),
	}
}

// Register adds the giving function into the list of available functions for instanction.
func (ng *GeneratorRegistry) Register(id string, fun Generator) bool {
	ng.functions[id] = fun
	return true
}

// RegisterTOML adds the giving function into the list of available functions for instanction,
// where it's config would be loaded using toml as the config unmarshaller.
func (ng *GeneratorRegistry) RegisterTOML(id string, fun Function) bool {
	ng.functions[id] = func(config []byte) (Op, error) {
		instance := fun()
		if len(config) == 0 {
			return instance, nil
		}

		if _, err := toml.DecodeReader(bytes.NewBuffer(config), instance); err != nil {
			return nil, err
		}

		return instance, nil
	}

	return true
}

// RegisterJSON adds the giving function into the list of available functions for instanction,
// where it's config would be loaded using json as the config unmarshaller.
func (ng *GeneratorRegistry) RegisterJSON(id string, fun Function) bool {
	ng.functions[id] = func(config []byte) (Op, error) {
		instance := fun()
		if len(config) == 0 {
			return instance, nil
		}

		if err := json.Unmarshal(config, instance); err != nil {
			return nil, err
		}

		return instance, nil
	}

	return true
}

// MustCreateFromBytes panics if giving function for id is not found.
func (ng *GeneratorRegistry) MustCreateFromBytes(id string, config []byte) Op {
	spell, err := ng.CreateFromBytes(id, config)
	if err != nil {
		panic(err)
	}

	return spell
}

// CreateFromBytes returns a new spell from the provided configuration
func (ng *GeneratorRegistry) CreateFromBytes(id string, config []byte) (Op, error) {
	fun, ok := ng.functions[id]
	if !ok {
		return nil, ErrFunctionNotFound
	}

	return fun(config)
}

// CreateWithTOML returns a new spell from the provided configuration map which is first
// converted into JSON then loaded using the CreateFromBytes function.
func (ng *GeneratorRegistry) CreateWithTOML(id string, config interface{}) (Op, error) {
	var bu bytes.Buffer

	if err := toml.NewEncoder(&bu).Encode(config); err != nil {
		return nil, err
	}

	return ng.CreateFromBytes(id, bu.Bytes())
}

// CreateWithJSON returns a new spell from the provided configuration map which is first
// converted into JSON then loaded using the CreateFromBytes function.
func (ng *GeneratorRegistry) CreateWithJSON(id string, config interface{}) (Op, error) {
	data, err := json.Marshal(config)
	if err != nil {
		return nil, err
	}

	return ng.CreateFromBytes(id, data)
}
