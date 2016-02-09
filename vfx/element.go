package vfx

import (
	"github.com/influx6/faux/ds"
	"honnef.co/go/js/dom"
)

// Elemental defines a dom.Element with read-write abilities for css properties.
// Calling th Sync() method of an elemental adjusts the changes as a batch into
// the real browser dom for that specific dom element.
type Elemental interface {
	dom.Element
	Read(string) (string, int, bool)
	Write(string, string, int)
	Sync()
}

// Element defines a structure that holds ehances the dom.Element api.
// Element provides a caching facility that helps to reduce layout checks
// and improve animation by returning last used data. Also it provides
// an appropriate method to update element properties apart from usings
// inlined styles.
type Element struct {
	dom.Element
	css        ComputedStyleMap // css holds the map of computed styles.
	cssDiff    ds.TruthTable    // contains lists of properties that have be change.
	diffStyles dom.StyleSheet
}

// NewElement returns an instancee of the Element struct.
func NewElement(elem dom.Element, pseudo string) *Element {
	css, err := GetComputedStyleMap(elem, pseudo)

	if err != nil {
		panic(err)
	}

	em := Element{
		Element: elem,
		css:     css,
		cssDiff: ds.NewTruthMap(ds.NewBoolStore()),
	}

	return &em
}

// Read reads out the elements internal css property rules and returns its
// value and priority(wether it has !important attached).
// If the property does not exists a false value is returned.
func (e *Element) Read(prop string) (string, int, bool) {
	cs, err := e.css.Get(prop)
	if err != nil {
		return "", 0, false
	}

	// Read the value, return both value and true state.
	return cs.Value, cs.Priority, true
}

// Write adds the necessary change of value to the giving property
// with the necessary adjustments. If the property is not found in
// the elements css stylesheet rules, it will be ignored.
func (e *Element) Write(prop string, value string, priority int) {
	cs, err := e.css.Get(prop)
	if err != nil {
		return
	}

	cs.Value = value
	cs.Priority = priority

	// Add the property into our diff map to ensure we deal with this
	// efficiently without re-writing untouched rules.
	e.cssDiff.Set(prop)
}

// Sync adjusts the necessary property changes of the giving element back into
// the dom. Any changes made to any properties will be diffed and added.
// Sync only re-writes change properties, all untouched onces are left alone.
func (e *Element) Sync() {

}
