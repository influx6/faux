package reflection

import (
	"errors"
	"reflect"
	"strings"
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

// MatchElement attempts to validate that both element are equal in type and value.
func MatchElement(me interface{}, other interface{}, allowFunctions bool) bool {
	meType := reflect.TypeOf(me)
	otherType := reflect.TypeOf(other)

	// if one is pointer, then both must be.
	if meType.Kind() == reflect.Ptr && otherType.Kind() != reflect.Ptr {
		return false
	}

	// if one is pointer, then both must be.
	if otherType.Kind() == reflect.Ptr && meType.Kind() != reflect.Ptr {
		return false
	}

	if meType.Kind() == reflect.Ptr {
		meType = meType.Elem()
	}

	if otherType.Kind() == reflect.Ptr {
		otherType = otherType.Elem()
	}

	if meType.Kind() == reflect.Func {
		if allowFunctions {
			return MatchFunction(me, other)
		}

		return false
	}

	if otherType.AssignableTo(meType) {
		return true
	}

	return true
}

// MatchFunction attempts to validate if giving types are functions and
// exactly match in arguments and returns.
func MatchFunction(me interface{}, other interface{}) bool {
	meType := reflect.TypeOf(me)
	otherType := reflect.TypeOf(other)

	// if one is pointer, then both must be.
	if meType.Kind() == reflect.Ptr && otherType.Kind() != reflect.Ptr {
		return false
	}

	// if one is pointer, then both must be.
	if otherType.Kind() == reflect.Ptr && meType.Kind() != reflect.Ptr {
		return false
	}

	if meType.Kind() == reflect.Ptr {
		meType = meType.Elem()
	}

	if otherType.Kind() == reflect.Ptr {
		otherType = otherType.Elem()
	}

	if meType.Kind() != reflect.Func {
		return false
	}

	if otherType.Kind() != reflect.Func {
		return false
	}

	if otherType.NumIn() != meType.NumIn() {
		return false
	}

	if otherType.NumOut() != meType.NumOut() {
		return false
	}

	for i := 0; i < meType.NumIn(); i++ {
		item := meType.In(i)
		otherItem := otherType.In(i)
		if item.Kind() != reflect.Func && !MatchElement(item, otherItem, true) {
			return false
		}

		if item.Kind() == reflect.Func && !MatchFunction(item, otherItem) {
			return false
		}
	}

	for i := 0; i < meType.NumOut(); i++ {
		item := meType.Out(i)
		otherItem := otherType.Out(i)
		if item.Kind() != reflect.Func && !MatchElement(item, otherItem, true) {
			return false
		}

		if item.Kind() == reflect.Func && !MatchFunction(item, otherItem) {
			return false
		}
	}

	return true
}

// TypeAndFields returns the type of the giving element and a slice of
// all composed types.
func TypeAndFields(elem interface{}) (reflect.Type, []reflect.Type, error) {
	tl := reflect.ValueOf(elem)

	if tl.Kind() == reflect.Ptr {
		tl = tl.Elem()
	}

	if tl.Kind() != reflect.Struct {
		return nil, nil, ErrNotStruct
	}

	var embeded []reflect.Type

	for ind := 0; ind < tl.NumField(); ind++ {
		item := tl.Field(ind).Type()
		embeded = append(embeded, item)
	}

	return tl.Type(), embeded, nil
}

var internalTypes = []string{
	"string", "int", "int64", "int32", "float32",
	"float64", "bool", "char",
	"rune", "byte",
}

// FieldType defines a struct which holds details the name and package which
// a giving field belongs to.
type FieldType struct {
	TypeName string `json:"field_type"`
	Pkg      string `json:"pkg"`
}

// ExternalTypeNames returns the name and field names of the provided
// elem which must be a struct, excluding all internal types.
func ExternalTypeNames(elem interface{}) (FieldType, []FieldType, error) {
	tl := reflect.TypeOf(elem)

	if tl.Kind() == reflect.Ptr {
		tl = tl.Elem()
	}

	if tl.Kind() == reflect.Interface {
		tl = tl.Elem()
	}

	if tl.Kind() != reflect.Struct {
		return FieldType{}, nil, ErrNotStruct
	}

	var item FieldType
	item.TypeName = tl.Name()
	item.Pkg = tl.PkgPath()

	var embeded []FieldType

	{
	fieldLoop:
		for ind := 0; ind < tl.NumField(); ind++ {
			item := tl.Field(ind).Type

			for _, except := range internalTypes {
				if except == item.Name() {
					continue fieldLoop
				}
			}

			if item.Name() == "" {
				if item.Kind() == reflect.Slice || item.Kind() == reflect.Array || item.Kind() == reflect.Interface {
					for _, except := range internalTypes {
						if except == item.Elem().Name() {
							continue fieldLoop
						}
					}

					embeded = append(embeded, FieldType{
						TypeName: item.Elem().Name(),
						Pkg:      item.Elem().PkgPath(),
					})
				}

				continue
			}

			embeded = append(embeded, FieldType{
				TypeName: item.Name(),
				Pkg:      item.PkgPath(),
			})
		}
	}

	return item, embeded, nil
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

// MatchFuncArgumentTypeWithValues validates specific values matches the elems
// function arguments.
func MatchFuncArgumentTypeWithValues(elem interface{}, vals []reflect.Value) int {
	ma, err := GetFuncArgumentsType(elem)
	if err != nil {
		return -1
	}

	if len(ma) != len(vals) {
		return -1
	}

	for index, item := range ma {
		val := vals[index]

		if ok, _ := CanSetFor(item, val); !ok {
			return index
		}
	}

	return -1
}

// StrictCanSetForType checks if a val reflect.Type can be used for the target type.
// It returns true/false if value matches the expected type and another true/false
// if the value can be converted to the expected type.
// Difference between this version and the other CanSet is that, it returns
// only true/false for the Assignability of the types and not based on the Assignability
// and convertibility.
func StrictCanSetForType(target, val reflect.Type) (canSet bool, mustConvert bool) {
	if val.AssignableTo(target) {
		canSet = true
	}

	if val.ConvertibleTo(target) {
		mustConvert = true
	}

	return
}

// CanSetForType checks if a val reflect.Type can be used for the target type.
// It returns true bool, where the first returns if the value can be used and if
// it must be converted into the type first.
func CanSetForType(target, val reflect.Type) (canSet bool, mustConvert bool) {
	if val.AssignableTo(target) {
		canSet = true
		return
	}

	if val.ConvertibleTo(target) {
		canSet = true
		mustConvert = true
		return
	}

	return
}

// CanSetFor checks if the giving val can be set in the place of the target type.
// It returns true bool, where the first returns if the value can be used and if
// it must be converted into the type first.
func CanSetFor(target reflect.Type, val reflect.Value) (canSet bool, mustConvert bool) {
	valType := val.Type()

	if valType.AssignableTo(target) {
		canSet = true
		return
	}

	if valType.ConvertibleTo(target) {
		canSet = true
		mustConvert = true
		return
	}

	return
}

// Convert takes a val and converts it into the target type provided if possible.
func Convert(target reflect.Type, val reflect.Value) (reflect.Value, error) {
	valType := val.Type()

	if !valType.ConvertibleTo(target) {
		return reflect.Value{}, errors.New("Can not convert type")
	}

	return val.Convert(target), nil
}

// MakeValueFor makes a new reflect.Value for the reflect.Type.
func MakeValueFor(t reflect.Type) reflect.Value {
	mtl := reflect.New(t)

	if t.Kind() != reflect.Ptr && mtl.Kind() == reflect.Ptr {
		mtl = mtl.Elem()
	}

	return mtl
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

//==============================================================================

// ErrNotStruct is returned when the reflect type is not a struct.
var ErrNotStruct = errors.New("Not a struct type")

// Field defines a specific tag field with its details from a giving struct.
type Field struct {
	Index int
	Name  string
	Tag   string
	NameLC string
	TypeName string
	Type  reflect.Type
	Value reflect.Value
	IsSlice bool
	IsArray bool
	IsMap bool
	IsChan bool
	IsStruct bool
}

// Fields defines a lists of Field instances.
type Fields []Field

// GetTagFields retrieves all fields of the giving elements with the giving tag
// type.
func GetTagFields(elem interface{}, tag string, allowNaturalNames bool) (Fields, error) {
	if !IsStruct(elem) {
		return nil, ErrNotStruct
	}

	tl := reflect.TypeOf(elem)
	tlVal := reflect.ValueOf(elem)

	if tl.Kind() == reflect.Ptr {
		tl = tl.Elem()
	}

	if tlVal.Kind() == reflect.Ptr {
		if tlVal.IsNil() {
			return nil, errors.New("invalid value: must be non-nil struct")
		}

		tlVal = tlVal.Elem()
	}

	var fields Fields

	for i := 0; i < tl.NumField(); i++ {
		field := tl.Field(i)

		// Get the specified tag from this field if it exists.
		tagVal := strings.TrimSpace(field.Tag.Get(tag))

		// If its a - item in the tag then skip or if its an empty string.
		if tagVal == "-" {
			continue
		}

		if !allowNaturalNames && tagVal == "" {
			continue
		}

		if tagVal == "" {
			tagVal = strings.ToLower(field.Name)
		}

		fields = append(fields, Field{
			Index:  i,
			Tag:    tagVal,
			Name:   field.Name,
			Type:   field.Type,
			Value:  tlVal.Field(i),
			TypeName: field.Type.Name(),
			NameLC: strings.ToLower(field.Name),
			IsMap: field.Type.Kind() == reflect.Map,
			IsChan: field.Type.Kind() == reflect.Chan,
			IsSlice: field.Type.Kind() == reflect.Slice,
			IsArray: field.Type.Kind() == reflect.Array,
			IsStruct: field.Type.Kind() == reflect.Struct,
		})
	}

	return fields, nil
}

// GetFieldByTagAndValue returns a giving struct field which has giving tag and value pair.
func GetFieldByTagAndValue(elem interface{}, tag string, value string) (reflect.StructField, error) {
	if !IsStruct(elem) {
		return reflect.StructField{}, ErrNotStruct
	}

	tl := reflect.TypeOf(elem)
	if tl.Kind() == reflect.Ptr {
		tl = tl.Elem()
	}

	for i := 0; i < tl.NumField(); i++ {
		field := tl.Field(i)
		if field.Tag.Get(tag) == value {
			return field, nil
		}
	}

	return reflect.StructField{}, ErrNoFieldWithTagFound
}

// GetFields retrieves all fields of the giving elements with the giving tag
// type.
func GetFields(elem interface{}) (Fields, error) {
	if !IsStruct(elem) {
		return nil, ErrNotStruct
	}

	tl := reflect.TypeOf(elem)
	tlVal := reflect.ValueOf(elem)

	if tl.Kind() == reflect.Ptr {
		tl = tl.Elem()
	}

	if tlVal.Kind() == reflect.Ptr {
		if tlVal.IsNil() {
			return nil, errors.New("invalid value: must be non-nil struct")
		}

		tlVal = tlVal.Elem()
	}

	var fields Fields

	for i := 0; i < tl.NumField(); i++ {
		field := tl.Field(i)
		fieldVal := Field{
			Index:  i,
			Name:   field.Name,
			Type:   field.Type,
			Value:  tlVal.Field(i),
			TypeName: field.Type.Name(),
			NameLC: strings.ToLower(field.Name),
			IsMap: field.Type.Kind() == reflect.Map,
			IsChan: field.Type.Kind() == reflect.Chan,
			IsSlice: field.Type.Kind() == reflect.Slice,
			IsArray: field.Type.Kind() == reflect.Array,
			IsStruct: field.Type.Kind() == reflect.Struct,
		}

		fields = append(fields, fieldVal)
	}

	return fields, nil
}

// ToMap returns a map of the giving values from a struct using a provided
// tag to capture the needed values, it extracts those tags values out into
// a map. It returns an error if the element is not a struct.
func ToMap(tag string, elem interface{}, allowNaturalNames bool) (map[string]interface{}, error) {
	// Collect the fields that match the giving tag.
	fields, err := GetTagFields(elem, tag, allowNaturalNames)
	if err != nil {
		return nil, err
	}

	// If there exists no field matching the tag skip.
	if len(fields) == 0 {
		return nil, errors.New("No Tag Matches")
	}

	data := make(map[string]interface{})

	// Loop through  the fields and set the appropriate value as needed.
	for _, field := range fields {
		if !field.Value.CanInterface() {
			continue
		}

		itemType := field.Type
		item := field.Value.Interface()

		switch itemType.Kind() {
		case reflect.Struct:
			if subItem, err := ToMap(tag, item, allowNaturalNames); err == nil {
				data[field.Tag] = subItem
			} else {
				data[field.Tag] = item
			}

		default:
			data[field.Tag] = item
		}
	}

	return data, nil
}

// MergeMap merges the key names of the provided map into the appropriate field
// place where the element has the provided tag.
func MergeMap(tag string, elem interface{}, values map[string]interface{}, allowAll bool) error {

	// Collect the fields that match the giving tag.
	fields, err := GetTagFields(elem, tag, allowAll)
	if err != nil {
		return err
	}

	// If there exists no field matching the tag skip.
	if len(fields) == 0 {
		return nil
	}

	tl := reflect.ValueOf(elem)

	if tl.Kind() == reflect.Ptr {
		tl = tl.Elem()
	}

	// Loop through  the fields and set the appropriate value as needed.
	for _, field := range fields {

		item := values[field.Tag]

		if item == nil {
			continue
		}

		fl := tl.Field(field.Index)

		// If we can't set this field, then skip.
		if !fl.CanSet() {
			continue
		}

		fl.Set(reflect.ValueOf(item))
	}

	return nil
}

// IsStruct returns true/false if the elem provided is a type of struct.
func IsStruct(elem interface{}) bool {
	mc := reflect.TypeOf(elem)

	if mc.Kind() == reflect.Ptr {
		mc = mc.Elem()
	}

	if mc.Kind() != reflect.Struct {
		return false
	}

	return true
}

// MakeNew returns a new version of the giving type, returning a nonpointer type.
// If the type is not a struct then an error is returned.
func MakeNew(elem interface{}) (interface{}, error) {
	mc := reflect.TypeOf(elem)

	if mc.Kind() != reflect.Struct {
		return nil, ErrNotStruct
	}

	return reflect.New(mc).Interface(), nil
}
