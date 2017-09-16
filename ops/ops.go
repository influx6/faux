package ops

import (
	"bytes"
	"encoding/json"
	"errors"
	"sync"

	"github.com/BurntSushi/toml"
)

// errors.
var (
	ErrFunctionNotFound = errors.New("Function with given id not found")
)

// CancelContext defines a type which provides Done signal for cancelling operations.
type CancelContext interface {
	Done() <-chan struct{}
}

// Op defines an interface which expose an exec method.
type Op interface {
	Exec(CancelContext) error
}

// Function defines a interface for a function which returns a giving Op.
type Function func() Op

// Generator defines a interface for a function which returns a giving Op.
type Generator func([]byte) (Op, error)

// GeneratorRegistry implements a structure that handles registration and ochestration of spells.
type GeneratorRegistry struct {
	fl        sync.Mutex
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
	ng.fl.Lock()
	defer ng.fl.Unlock()

	ng.functions[id] = fun
	return true
}

// RegisterTOML adds the giving function into the list of available functions for instanction,
// where it's config would be loaded using toml as the config unmarshaller.
func (ng *GeneratorRegistry) RegisterTOML(id string, fun Function) bool {
	ng.fl.Lock()
	defer ng.fl.Unlock()

	ng.functions[id] = func(config []byte) (Op, error) {
		instance := fun()
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
	ng.fl.Lock()
	defer ng.fl.Unlock()

	ng.functions[id] = func(config []byte) (Op, error) {
		instance := fun()
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
	ng.fl.Lock()
	defer ng.fl.Unlock()

	fun, ok := ng.functions[id]
	if !ok {
		return nil, ErrFunctionNotFound
	}

	return fun(config)
}

// CreateWithTOML returns a new spell from the provided configuration map which is first
// converted into JSON then loaded using the CreateFromBytes function.
func (ng *GeneratorRegistry) CreateWithTOML(id string, config map[string]interface{}) (Op, error) {
	var bu bytes.Buffer

	if err := toml.NewEncoder(&bu).Encode(config); err != nil {
		return nil, err
	}

	return ng.CreateFromBytes(id, bu.Bytes())
}

// CreateWithJSON returns a new spell from the provided configuration map which is first
// converted into JSON then loaded using the CreateFromBytes function.
func (ng *GeneratorRegistry) CreateWithJSON(id string, config map[string]interface{}) (Op, error) {
	data, err := json.Marshal(config)
	if err != nil {
		return nil, err
	}

	return ng.CreateFromBytes(id, data)
}