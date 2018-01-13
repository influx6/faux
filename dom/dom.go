package dom

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/gopherjs/gopherjs/js"
	"github.com/influx6/faux/utils"
	"honnef.co/go/js/dom"
)

// ErrNotFound returns when an item is not found.
var ErrNotFound = errors.New("Not found")

// GetComputedStyle returns the dom.Element computed css styles.
func GetComputedStyle(elem dom.Element, ps string) (*dom.CSSStyleDeclaration, error) {
	css := dom.GetWindow().GetComputedStyle(elem, ps)
	if css == nil {
		return nil, ErrNotFound
	}

	return css, nil
}

// RemoveComputedStyleValue removes the value of the property from the computed
// style object.
func RemoveComputedStyleValue(css *dom.CSSStyleDeclaration, prop string) {
	defer func() {
		recover()
	}()

	css.Call("removeProperty", prop)
}

// GetComputedStyleValue retrieves the value of the property from the computed
// style object.
func GetComputedStyleValue(elem dom.Element, psudo string, prop string) (*js.Object, error) {
	vs, err := GetComputedStyle(elem, psudo)
	if err != nil {
		return nil, err
	}

	vcs, err := GetComputedStyleValueWith(vs, prop)
	if err != nil {
		return nil, err
	}

	return vcs, nil
}

// GetComputedStyleValueWith usings the CSSStyleDeclaration to
// retrieves the value of the property from the computed
// style object.
func GetComputedStyleValueWith(css *dom.CSSStyleDeclaration, prop string) (*js.Object, error) {
	vs := css.Call("getPropertyValue", prop)
	if vs == nil {
		return nil, ErrNotFound
	}

	return vs, nil
}

// GetComputedStylePriority retrieves the proritiy of the property from the computed
// style object.
func GetComputedStylePriority(css *dom.CSSStyleDeclaration, prop string) (int, error) {
	vs := css.Call("getPropertyPriority", prop)
	if vs == nil {
		return 0, ErrNotFound
	}

	if strings.TrimSpace(vs.String()) == "" {
		return 0, nil
	}

	return 1, nil
}

// GetProp retrieves the necessary property for this specific name.
func GetProp(o *js.Object, prop string) (*js.Object, error) {

	// Expand the property for possible period delimited sets.
	props := Expando(prop)

	var jsop *js.Object

	// Loop the property sets and get the next one
	for i := 0; i < len(props); i++ {
		jsop = o.Get(prop)

		if jsop == nil {
			return nil, ErrNotFound
		}
	}

	return jsop, nil
}

//==============================================================================

// MatchProp matches the string value from a property val against a provided
// expected value.
func MatchProp(o *js.Object, prop string, val string) bool {
	op, err := GetProp(o, prop)
	if err != nil {
		return false
	}

	if strings.ToLower(op.String()) != val {
		return false
	}

	return true
}

//==============================================================================

// ComputedStyle defines a style property item.
type ComputedStyle struct {
	Name       string
	VendorName string
	Value      string
	Values     []string
	Priority   bool // values between [0,1] to indicate use of '!important'
}

// ComputedStyleMap defines a map type of computed style properties and values.
type ComputedStyleMap map[string]ComputedStyle

// Has returns true/false if the property exists.
func (c ComputedStyleMap) Has(name string) bool {
	_, ok := c[name]
	return ok
}

// Get retrieves the specific property if it exists.
func (c ComputedStyleMap) Get(name string) (ComputedStyle, error) {
	cs, ok := c[name]
	if !ok {
		return ComputedStyle{}, ErrNotFound
	}

	return cs, nil
}

// GetComputedStyleMap returns a map of computed style properties and values.
// Also all vendored names are cleaned up to allow quick and easy access
// regardless of vendor.
func GetComputedStyleMap(elem dom.Element, ps string) (ComputedStyleMap, error) {
	css, err := GetComputedStyle(elem, ps)
	if err != nil {
		return nil, err
	}

	styleMap := make(ComputedStyleMap)

	// Get the map and pull the necessary property:value and importance facts.
	for key, val := range css.ToMap() {
		priority, _ := GetComputedStylePriority(css, key)

		unvendoredName := key

		// Clean key of any vendored name to allow easy access.
		for _, vo := range vendorTags {
			unvendoredName = strings.TrimPrefix(unvendoredName, fmt.Sprintf("%s", vo))
			unvendoredName = strings.TrimPrefix(unvendoredName, fmt.Sprintf("-%s-", vo))
		}

		var vals []string

		if strings.TrimSpace(val) != "none" {
			vals = append(vals, val)
		}

		styleMap[unvendoredName] = ComputedStyle{
			Name:       unvendoredName,
			VendorName: key,
			Value:      val,
			Values:     vals,
			Priority:   (priority > 0),
		}
	}

	return styleMap, nil
}

// vendorTags provides a lists of different browser specific vendor names.
var vendorTags = []string{"moz", "webki", "O", "ms"}

// Vendorize returns a property name with the different versions known according
// browsers.
func Vendorize(u string) []string {
	var v []string

	for _, vn := range vendorTags {
		v = append(v, fmt.Sprintf("-%s-%s", vn, u))
	}

	return v
}

var unitRex = regexp.MustCompile("([\\+\\-]?[0-9|auto\\.]+)(%|px|pt|em|rem|in|cm|mm|ex|pc|vw|vh|deg)?")

// Unit returns a valid unit type in the browser, if the supplied unit is
// standard then it is return else 'px' is returned as default.
func Unit(u string) string {
	if vals := unitRex.FindStringSubmatch(u); len(vals) > 1 {
		return vals[1]
	}

	return ""
}

// doubleString doubles the giving string.
func doubleString(c string) string {
	return fmt.Sprintf("%s%s", c, c)
}

//==============================================================================

// simpleRotationMatch defines a matcher for the formation rotate(90deg).
var simpleRotationMatch = regexp.MustCompile("rotate\\(([\\d]+)deg\\)")

// IsSimpleRotation checks wether the giving string is a css rotation directive.
func IsSimpleRotation(data string) bool {
	return simpleRotationMatch.MatchString(data)
}

// rotationMatch defines a matcher for the formation rotate(90deg).
var rotationMatch = regexp.MustCompile("[rotate|rotateX|rotateY|rotateZ]\\(([\\d]+)deg\\)")

// Rotation defines the concrete representation of the css3 skew
// transform property.
type Rotation struct {
	Angle float64
}

// IsRotation checks wether the giving string is a css rotation directive.
func IsRotation(data string) bool {
	return rotationMatch.MatchString(data)
}

// ToRotation returns the rotation from the giving string else returns
// an error if it failed.
func ToRotation(data string) (*Rotation, error) {
	if !IsRotation(data) {
		return nil, errors.New("Invalid Data")
	}

	ts := strings.Split(rotationMatch.FindStringSubmatch(data)[1], ",")

	var t Rotation
	t.Angle = utils.ParseFloat(ts[0])

	return &t, nil
}

//==============================================================================

// skewMatch defines a matcher for the format eg skew(90px,40px).
var skewMatch = regexp.MustCompile("[skew|skewX|skewY]\\(([\\d,\\s]+)\\)")

// Skew defines the concrete representation of the css3 skew
// transform property.
type Skew struct {
	X float64
	Y float64
}

// IsSkew checks wether the giving string is a css skew directive.
func IsSkew(data string) bool {
	return skewMatch.MatchString(data)
}

// ToSkew returns the skew from the giving string else returns
// an error if it failed.
func ToSkew(data string) (*Skew, error) {
	if !IsSkew(data) {
		return nil, errors.New("Invalid Data")
	}

	ts := strings.Split(skewMatch.FindStringSubmatch(data)[1], ",")

	var t Skew

	if strings.HasSuffix(data, "Y") {
		t.Y = utils.ParseFloat(ts[0])
	} else if strings.Contains(data, "X") {
		t.X = utils.ParseFloat(ts[0])
	} else {
		t.X = utils.ParseFloat(ts[0])
		t.Y = utils.ParseFloat(ts[1])
	}

	return &t, nil
}

//==============================================================================

var scaleMatch = regexp.MustCompile("[translate|translateX|translateY]\\(([\\d,\\s]+)\\)")

// Scale defines the concrete representation of the css3 scale
// transform property.
type Scale struct {
	X float64
	Y float64
}

// IsScale checks wether the giving string is a css scale directive.
func IsScale(data string) bool {
	return scaleMatch.MatchString(data)
}

// ToScale returns the translation from the giving string else returns
// an error if it failed.
func ToScale(data string) (*Scale, error) {
	if !IsScale(data) {
		return nil, errors.New("Invalid Data")
	}

	ts := strings.Split(scaleMatch.FindStringSubmatch(data)[1], ",")

	var t Scale

	if strings.HasSuffix(data, "Y") {
		t.Y = utils.ParseFloat(ts[0])
	} else if strings.Contains(data, "X") {
		t.X = utils.ParseFloat(ts[0])
	} else {
		t.X = utils.ParseFloat(ts[0])
		t.Y = utils.ParseFloat(ts[1])
	}

	return &t, nil
}

//==============================================================================

var perspectiveMatch = regexp.MustCompile("perspective\\(([\\d,\\s]+)\\)")

// IsPerspective checks wether the giving string is a css perspective directive.
func IsPerspective(data string) bool {
	return perspectiveMatch.MatchString(data)
}

// Perspective provides a structure for storing current perspective data.
type Perspective struct {
	Range float64
}

// ToPerspective returns the translation from the giving string else returns
// an error if it failed.
func ToPerspective(data string) (*Perspective, error) {
	if !IsPerspective(data) {
		return nil, errors.New("Invalid Data")
	}

	ts := strings.Split(perspectiveMatch.FindStringSubmatch(data)[1], ",")

	var t Perspective
	t.Range = utils.ParseFloat(ts[0])

	return &t, nil
}

//==============================================================================

var tranlateMatch = regexp.MustCompile("translate\\(([\\d,\\s]+)\\)")

// Translation defines the concrete representation of the css3 translation
// transform property.
type Translation struct {
	X float64
	Y float64
}

// IsTranslation checks wether the giving string is a css translation directive.
func IsTranslation(data string) bool {
	return tranlateMatch.MatchString(data)
}

// ToTranslation returns the translation from the giving string else returns
// an error if it failed.
func ToTranslation(data string) (*Translation, error) {
	if !IsTranslation(data) {
		return nil, errors.New("Invalid Data")
	}

	ts := strings.Split(tranlateMatch.FindStringSubmatch(data)[1], ",")

	var t Translation

	if strings.HasSuffix(data, "Y") {
		t.Y = utils.ParseFloat(ts[0])
	} else if strings.Contains(data, "X") {
		t.X = utils.ParseFloat(ts[0])
	} else {
		t.X = utils.ParseFloat(ts[0])
		t.Y = utils.ParseFloat(ts[1])
	}

	return &t, nil
}

//==============================================================================

var matrixMatch = regexp.MustCompile("matrix(3[d|D])?\\(([,\\d\\s]+)\\)")

// IsMatrix returns true/false if the giving string is a matrix declaration.
func IsMatrix(data string) bool {
	if !matrixMatch.MatchString(data) {
		return false
	}

	ms := strings.Split(matrixMatch.FindStringSubmatch(data)[1], ",")

	if len(ms) < 6 {
		return false
	}

	return true
}

// Matrix defines a transformation matrix generated from a transform directive.
type Matrix struct {
	ScaleX    float64
	RotationX float64
	ScaleY    float64
	RotationY float64
	PositionX float64
	PositionY float64
	PositionZ float64
	ScaleZ    float64
	RotationZ float64
}

// ToMatrix2D returns a matrix from the provided data (eg matrix(0,1,0,1,3,4))
// else returns an error.
func ToMatrix2D(data string) (*Matrix, error) {
	if !IsMatrix(data) {
		return nil, errors.New("Invalid Matrix data")
	}

	ms := strings.Split(matrixMatch.FindStringSubmatch(data)[1], ",")

	if len(ms) < 6 {
		return nil, errors.New("Invalid Matrix data")
	}

	m := Matrix{
		ScaleX:    utils.ParseFloat(ms[0]),
		RotationX: utils.ParseFloat(ms[1]),
		ScaleY:    utils.ParseFloat(ms[2]),
		RotationY: utils.ParseFloat(ms[3]),
		PositionX: utils.ParseFloat(ms[4]),
		PositionY: utils.ParseFloat(ms[5]),
	}

	return &m, nil
}

//==============================================================================

// expandable defines a regexp for matching period delimited strings.
var expandable = regexp.MustCompile("([\\w\\d_-]+\\.[\\w\\d_-]+)+")

// Expando expands a property period delimited string into its component parts.
func Expando(prop string) []string {
	if !expandable.MatchString(prop) {
		return []string{prop}
	}

	return strings.Split(prop, ".")
}
