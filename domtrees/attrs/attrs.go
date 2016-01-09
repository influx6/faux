//Package attr provides attributes for html base elements

package attrs

import (
	"github.com/influx6/faux/domtrees"
)

// InputType defines the set type of input values for the input elements
type InputType string

// Types of input for attribute input types.
const (
	TypeButton InputType = "button"

	TypeCheckbox InputType = "checkbox"

	TypeColor InputType = "color"

	TypeDate InputType = "date"

	TypeDatetime InputType = "datetime"

	TypeDatetimelocal InputType = "datetime-local"

	TypeEmail InputType = "email"

	TypeFile InputType = "file"

	TypeHidden InputType = "hidden"

	TypeImage InputType = "image"

	TypeMonth InputType = "month"

	TypeNumber InputType = "number"

	TypePassword InputType = "password"

	TypeRadio InputType = "radio"

	TypeRange InputType = "range"

	TypeMin InputType = "min"

	TypeMax InputType = "max"

	TypeValue InputType = "value"

	TypeStep InputType = "step"

	TypeReset InputType = "reset"

	TypeSearch InputType = "search"

	TypeSubmit InputType = "submit"

	TypeTel InputType = "tel"

	TypeText InputType = "text"

	TypeTime InputType = "time"

	TypeURL InputType = "url"

	TypeWeek InputType = "week"
)

// Name defines attributes of type "Name" for html element types
func Name(val string) *domtrees.Attribute {
	return &domtrees.Attribute{Name: "name", Value: val}
}

// Checked defines attributes of type "Checked" for html element types
func Checked(val string) *domtrees.Attribute {
	return &domtrees.Attribute{Name: "checked", Value: val}
}

// ClassName defines attributes of type "ClassName" for html element types
func ClassName(val string) *domtrees.Attribute {
	return &domtrees.Attribute{Name: "className", Value: val}
}

// Autofocus defines attributes of type "Autofocus" for html element types
func Autofocus(val string) *domtrees.Attribute {
	return &domtrees.Attribute{Name: "autofocus", Value: val}
}

// ID defines attributes of type "Id" for html element types
func ID(val string) *domtrees.Attribute {
	return &domtrees.Attribute{Name: "id", Value: val}
}

// HTMLFor defines attributes of type "HtmlFor" for html element types
func HTMLFor(val string) *domtrees.Attribute {
	return &domtrees.Attribute{Name: "htmlFor", Value: val}
}

// Class defines attributes of type "Class" for html element types
func Class(val string) *domtrees.Attribute {
	return &domtrees.Attribute{Name: "class", Value: val}
}

// Src defines attributes of type "Src" for html element types
func Src(val string) *domtrees.Attribute {
	return &domtrees.Attribute{Name: "src", Value: val}
}

// Href defines attributes of type "Href" for html element types
func Href(val string) *domtrees.Attribute {
	return &domtrees.Attribute{Name: "href", Value: val}
}

// Rel defines attributes of type "Rel" for html element types
func Rel(val string) *domtrees.Attribute {
	return &domtrees.Attribute{Name: "rel", Value: val}
}

// IType defines attributes of type "Type" for html element types
func IType(val InputType) *domtrees.Attribute {
	return Type(string(val))
}

// Type defines attributes of type "Type" for html element types
func Type(val string) *domtrees.Attribute {
	return &domtrees.Attribute{Name: "type", Value: val}
}

// Placeholder defines attributes of type "Placeholder" for html element types
func Placeholder(val string) *domtrees.Attribute {
	return &domtrees.Attribute{Name: "placeholder", Value: val}
}

// Value defines attributes of type "Value" for html element types
func Value(val string) *domtrees.Attribute {
	return &domtrees.Attribute{Name: "value", Value: val}
}
