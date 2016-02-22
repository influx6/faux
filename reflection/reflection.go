package reflection

import (
	"errors"
	"reflect"
)

// ErrNotFunction is returned when the type is not a reflect.Func.
var ErrNotFunction = errors.New("Not A Function Type")

// IsFuncType returns true/false if the interface provided is a func type.
func IsFuncType(elem interface{}) bool {
	_, err := FuncType(elem)
	if err != nil {
		return false
	}
	return true
}

// FuncValue return the Function reflect.Value of the interface provided else
// returns a non-nil error.
func FuncValue(elem interface{}) (reflect.Value, error) {
	tl := reflect.ValueOf(elem)

	if tl.Kind() == reflect.Ptr {
		tl = tl.Elem()
	}

	if tl.Kind() != reflect.Func {
		return tl, ErrNotFunction
	}

	return tl, nil
}

// FuncType return the Function reflect.Type of the interface provided else
// returns a non-nil error.
func FuncType(elem interface{}) (reflect.Type, error) {
	tl := reflect.TypeOf(elem)

	if tl.Kind() == reflect.Ptr {
		tl = tl.Elem()
	}

	if tl.Kind() != reflect.Func {
		return nil, ErrNotFunction
	}

	return tl, nil
}

// HasArgumentSize return true/false to indicate if the function type has the
// size of arguments. It will return false if the interface is not a function
// type.
func HasArgumentSize(elem interface{}, len int) bool {
	tl := reflect.TypeOf(elem)

	if tl.Kind() == reflect.Ptr {
		tl = tl.Elem()
	}

	if tl.Kind() != reflect.Func {
		return false
	}

	if tl.NumIn() != len {
		return false
	}

	return true
}

// GetFuncArgumentsType returns the arguments type of function which should be
// a function type,else returns a non-nil error.
func GetFuncArgumentsType(elem interface{}) ([]reflect.Type, error) {
	tl := reflect.TypeOf(elem)

	if tl.Kind() == reflect.Ptr {
		tl = tl.Elem()
	}

	if tl.Kind() != reflect.Func {
		return nil, ErrNotFunction
	}

	totalFields := tl.NumIn()

	var input []reflect.Type

	for i := 0; i < totalFields; i++ {
		indElem := tl.In(i)

		// if indElem.Kind() == reflect.Ptr {
		// 	indElem = indElem.Elem()
		// }

		input = append(input, indElem)
	}

	return input, nil
}

// MakeValueFor makes a new reflect.Value for the reflect.Type.
func MakeValueFor(t reflect.Type) reflect.Value {
	var input reflect.Value

	mtl := reflect.New(t)

	if mtl.Kind() == reflect.Ptr {
		mtl = mtl.Elem()
	}

	return input
}

// MakeArgumentsValues takes a list of reflect.Types and returns a new version of
// those types, ensuring to dereference if it receives a pointer reflect.Type.
func MakeArgumentsValues(args []reflect.Type) []reflect.Value {
	var inputs []reflect.Value

	for _, tl := range args {
		inputs = append(inputs, MakeValueFor(tl))
	}

	return inputs
}

// InterfaceFromValues returns a list of interfaces representing the concrete
// values within the lists of reflect.Value types.
func InterfaceFromValues(vals []reflect.Value) []interface{} {
	var data []interface{}

	for _, val := range vals {
		if val.Kind() == reflect.Ptr {
			val = val.Elem()
		}

		data = append(data, val.Interface())
	}

	return data
}
