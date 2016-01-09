package domtrees

import (
	"fmt"
	"strings"

	"github.com/influx6/faux/domevents"
)

// Mutation defines the capability of an element to state its
// state if mutation occured
type Mutation interface {
	UID() string
	Hash() string
	swapHash(string)
	swapUID(string)
	Removed() bool
	Remove()
	UpdateHash()
}

// Eventers provide an interface type for elements able to register and load
// event managers.
type Eventers interface {
	LoadEvents()
	UseEventManager(domevents.EventManagers) bool
}

// Markup provide a basic specification type of how a element resolves its content
type Markup interface {
	Appliable
	Events
	Mutation
	Styles
	Attributes
	Clonable
	MarkupChildren
	Reconcilable
	Eventers

	AutoClosed() bool
	TextContent() string

	Name() string
	EventID() string
	Empty()

	CleanRemoved()
}

// Mutable is a base implementation of the Mutation interface{}
type Mutable struct {
	uid     string
	hash    string
	removed bool
}

// MutableBackdoor grants access to a mutable swapUID and swapHash methods.
type MutableBackdoor struct {
	M Mutation
}

// SwapUID swaps the uid of the internal mutable.
func (m MutableBackdoor) SwapUID(uid string) {
	m.M.swapUID(uid)
}

// SwapHash swaps the hash of the internal mutable.
func (m MutableBackdoor) SwapHash(hash string) {
	m.M.swapHash(hash)
}

// NewMutable returns a new mutable instance
func NewMutable() *Mutable {
	m := &Mutable{
		uid:  RandString(8),
		hash: RandString(10),
	}
	return m
}

// UpdateHash updates the mutable hash value
func (m *Mutable) UpdateHash() {
	m.hash = RandString(10)
}

// Remove marks this mutable as removed making this irreversible
func (m *Mutable) Remove() {
	m.removed = true
}

// Removed returns true/false if the mutable is marked removed
func (m *Mutable) Removed() bool {
	return !!m.removed
}

// Hash returns the current hash of the mutable
func (m *Mutable) Hash() string {
	return m.hash
}

// UID returns the current uid of the mutable
func (m *Mutable) UID() string {
	return m.uid
}

// swapHash provides a backdoor to swap the hash
func (m *Mutable) swapHash(h string) {
	m.hash = h
}

// swapUID provides a backdoor to swap uid
func (m *Mutable) swapUID(uid string) {
	m.uid = uid
}

// Element represent a concrete implementation of a element node
type Element struct {
	Mutation
	tagname         string
	events          []*Event
	styles          []*Style
	attrs           []*Attribute
	children        []Markup
	textContent     string
	autoclose       bool
	allowEvents     bool
	allowChildren   bool
	allowStyles     bool
	allowAttributes bool
	eventManager    domevents.EventManagers
}

// NewText returns a new Text instance element
func NewText(txt string) *Element {
	em := NewElement("text", false)
	em.allowChildren = false
	em.allowAttributes = false
	em.allowStyles = false
	em.allowEvents = false
	em.textContent = txt
	return em
}

// NewElement returns a new element instance giving the specificed name
func NewElement(tag string, hasNoEndingTag bool) *Element {
	return &Element{
		Mutation:        NewMutable(),
		tagname:         strings.ToLower(strings.TrimSpace(tag)),
		children:        make([]Markup, 0),
		styles:          make([]*Style, 0),
		attrs:           make([]*Attribute, 0),
		autoclose:       hasNoEndingTag,
		allowChildren:   true,
		allowStyles:     true,
		allowAttributes: true,
		allowEvents:     true,
	}
}

// UseEventManager adds a eventmanager into the markup and if not available before automatically registers
// the events with it,once an event manager is registered to it,it will and can not be changed
func (e *Element) UseEventManager(man domevents.EventManagers) bool {
	if man == nil {
		return true
	}

	if e.eventManager != nil {
		// e.eventManager.
		man.AttachManager(e.eventManager)
		return false
	}

	e.eventManager = man
	e.LoadEvents()
	return true
}

// LoadEvents loads up the events registered by this and by its children into each respective
// available events managers
func (e *Element) LoadEvents() {
	if e.eventManager != nil {
		e.eventManager.DisconnectRemoved()
		// log.Printf("will load events: %s %+s", e.Name(), e.events)

		for _, ev := range e.events {
			if es, _ := e.eventManager.NewEventMeta(ev.Meta); es != nil {
				es.Bind(ev.Fx)
			}
		}

	}

	//load up the children events also
	for _, ech := range e.children {
		if !ech.UseEventManager(e.eventManager) {
			ech.LoadEvents()
		}
	}
}

// EventID returns the selector used for tagging events for a markup.
func (e *Element) EventID() string {
	return fmt.Sprintf("%s[uid='%s']", strings.ToLower(e.Name()), e.UID())
}

// Remove sets the markup as removable and adds a 'haikuRemoved' attribute to it
func (e *Element) Remove() {
	if !e.Removed() {
		e.attrs = append(e.attrs, &Attribute{"haikuRemoved", ""})
		e.Mutation.Remove()
	}
}

// Empty resets the elements children list as 0 length
func (e *Element) Empty() {
	e.children = e.children[:0]
}

// Name returns the tag name of the element
func (e *Element) Name() string {
	return e.tagname
}

// TextContent returns the elements text value if its a text type else an empty string
func (e *Element) TextContent() string {
	return e.textContent
}

// CleanRemoved removes all the chilren marked as removed
func (e *Element) CleanRemoved() {
	for n, em := range e.children {
		if em.Removed() {
			copy(e.children[n:], e.children[n+1:])
			e.children = e.children[:len(e.children)-1]
		} else {
			em.CleanRemoved()
		}
	}
}

// AutoClosed returns true/false if this element uses a </> or a <></> tag convention
func (e *Element) AutoClosed() bool {
	return e.autoclose
}

// Reconcilable defines the interface of markups that can reconcile their content against another
type Reconcilable interface {
	Reconcile(Markup) bool
}

// Reconcile checks each item against the given lists
// and replaces/add any missing item.
func (c *ClassList) Reconcile(m *ClassList) bool {
	var added bool
	maxlen := len(c.list)
	for ind, val := range m.list {
		if ind >= maxlen {
			added = true
			c.list = append(c.list, val)
			continue
		}

		if c.list[ind] == val {
			continue
		}

		added = true
		c.list[ind] = val
	}

	return added
}

// Reconcile checks if the attribute matches then upgrades its value.
func (a *Attribute) Reconcile(m *Attribute) bool {
	if strings.TrimSpace(a.Name) == strings.TrimSpace(m.Name) {
		a.Value = m.Value
		return true
	}
	return false
}

// Reconcile checks if the style matches then upgrades its value.
func (s *Style) Reconcile(m *Style) bool {
	if strings.TrimSpace(s.Name) == strings.TrimSpace(m.Name) {
		s.Value = m.Value
		return true
	}
	return false
}

// Reconcile takes a old markup and reconciles its uid and its children with the new formation,it returns a true/false
// telling the parent if the children swapped hashes
/* the reconcilation uses the order in which elements are added, if the order and element types are kept,
then the uid are swapped else it firsts checks the element type and if not the same adds the old one into
the new list as removed then continues the check. The system takes position of elements in the old and new
as very important and I cant stress this enough, "Element Positioning" in the markup are very important,
if a Anchor was the first element in the old render and the next pass returns a Div in the new render, the Anchor
will be marked as removed and will be removed from the dom and ignored by the writers. When two elements position
are same and their types are the same then a checkup process is doing using the elements attributes, this is done to
determine if the hash value of the new should be swapped with the old. We can use style properties here because they are
the most volatile of the set and will periodically be either changed and returned to normal values eg display: none to display: block and vise-versa, so only attributes are used in the check process
*/
func (e *Element) Reconcile(em Markup) bool {
	// are we reconciling the proper elements type ? if not skip (i.e different types cant reconcile eachother)]
	// TODO: decide if we should mark the markup as removed in this case as a catchall system
	if e.Name() != em.Name() {
		return false
	}

	em.CleanRemoved()

	//since the tagname are the same, swap uids
	// olduid := em.UID()
	e.swapUID(em.UID())

	//since the tagname are the same and we have swapped uid, to determine who gets or keeps
	// its hash we will check the attributes against each other, but also the hash is dependent on the
	// children also, if the children observered there was a change
	oldHash := em.Hash()
	// newHash := e.Hash()

	// if we have a special case for text element then we do things differently
	if e.Name() == "text" {
		//if the contents are equal,keep the prev hash
		if e.TextContent() == em.TextContent() {
			e.swapHash(oldHash)
			return false
		}
		return true
	}

	newChildren := e.Children()
	oldChildren := em.Children()
	maxSize := len(newChildren)
	oldMaxSize := len(oldChildren)

	if maxSize <= 0 {
		// if the element had no children too, swap hash
		if oldMaxSize <= 0 {
			e.swapHash(oldHash)
			return false
		}

		return true
	}

	var childChanged bool

	for n, och := range oldChildren {
		if maxSize > n {
			nch := newChildren[n]

			// log.Printf("checking old (%s) with new(%s)", och.Name(), nch.Name())

			if nch.Name() == och.Name() {
				if nch.Reconcile(och) {
					// log.Printf("old (%s) with new(%s) changed!", och.Name(), nch.Name())
					childChanged = true
				}
			} else {
				och.Remove()
				e.AddChild(och)
			}
			continue
		}

		och.Remove()
		e.AddChild(och)
	}

	ReconcileEvents(e, em)
	if e.eventManager != nil {
		e.eventManager.DisconnectRemoved()
	}

	//if the sizes of the new node is more than the old node then ,we definitely changed
	if maxSize > oldMaxSize {
		return true
	}

	// log.Printf("did child change: %+s -> %t", e.Name(), childChanged)
	// log.Printf("are attributes of children %s equal -> %t", e.Name(), EqualAttributes(e, em))

	if !childChanged && EqualAttributes(e, em) {
		// log.Printf("children did not change: %+s -> %+s", e.Name(), em.Name())
		e.swapHash(oldHash)
		return false
	}

	// log.Printf("\n")
	return true
}

// MarkupChildren defines the interface of an element that has children
type MarkupChildren interface {
	AddChild(...Markup)
	Children() []Markup
}

// AddChild adds a new markup as the children of this element
func (e *Element) AddChild(em ...Markup) {
	if e.allowChildren {
		for _, m := range em {
			e.children = append(e.children, m)
			//if this are free elements, then use this event manager
			m.UseEventManager(e.eventManager)
		}
	}
}

// Children returns the children list for the element
func (e *Element) Children() []Markup {
	return e.children
}

// ClassList defines the struct for a lists of classes.
type ClassList struct {
	list []string
}

// Remove removes a class from into the lists.
func (c *ClassList) Remove(class string) {
	var index = -1
	var val string

	for index, val = range c.list {
		if val == class {
			break
		}
	}

	if index < 0 {
		return
	}

	c.list = append(c.list[:index], c.list[index:]...)
}

// Add adds a class name into the lists.
func (c *ClassList) Add(class string) {
	c.list = append(c.list, class)
}

// Style define the style specification for element styles
type Style struct {
	Name  string
	Value string
}

// NewStyle returns a new style instance
func NewStyle(name, val string) *Style {
	s := Style{Name: name, Value: val}
	return &s
}

// Attribute define the struct  for attributes
type Attribute struct {
	Name  string
	Value string
}

// NewAttr returns a new attribute instance
func NewAttr(name, val string) *Attribute {
	a := Attribute{Name: name, Value: val}
	return &a
}

// Styles interface defines a type that has Styles
type Styles interface {
	Styles() []*Style
}

// Styles return the internal style list of the element
func (e *Element) Styles() []*Style {
	return e.styles
}

// Attributes interface defines a type that has Attributes
type Attributes interface {
	Attributes() []*Attribute
}

// Attributes return the internal attribute list of the element
func (e *Element) Attributes() []*Attribute {
	return e.attrs
}

// Events provide an interface for markup event addition system
type Events interface {
	Events() []*Event
}

// Event provide a meta registry for helps in registering events for dom markups
// which is translated to the nodes themselves
type Event struct {
	Meta *domevents.EventMetable
	Fx   domevents.EventHandler
	tree Markup
}

// EventHandler provides a custom event handler which allows access to the
// markup producing the event.
type EventHandler func(domevents.Event, Markup)

// NewEvent returns a event object that allows registering events to eventlisteners
func NewEvent(etype, eselector string, efx EventHandler) *Event {
	ex := Event{
		Meta: &domevents.EventMetable{EventType: etype, EventTarget: eselector},
	}

	// wireup the function to get the ev and tree.
	ex.Fx = func(ev domevents.Event) {
		if efx != nil {
			efx(ev, ex.tree)
		}
	}

	return &ex
}

// StopImmediatePropagation will return itself and set StopPropagation to true
func (e *Event) StopImmediatePropagation() *Event {
	e.Meta.ShouldStopImmediatePropagation = true
	return e
}

// StopPropagation will return itself and set StopPropagation to true
func (e *Event) StopPropagation() *Event {
	e.Meta.ShouldStopPropagation = true
	return e
}

// PreventDefault will return itself and set PreventDefault to true
func (e *Event) PreventDefault() *Event {
	e.Meta.ShouldPreventDefault = true
	return e
}

// Events return the elements events
func (e *Element) Events() []*Event {
	return e.events
}

// Appliable define the interface specification for applying changes to elements elements in tree
type Appliable interface {
	Apply(*Element)
}

// Apply checks for a class attribute
func (c *ClassList) Apply(e *Element) {
	if len(c.list) == 0 {
		return
	}

	list := strings.Join(c.list, " ")

	a, err := GetAttr(e, "class")

	if err != nil {
		(&Attribute{Name: "class", Value: "list"}).Apply(e)
		return
	}

	// TODO: should we make Apply smarter?
	// // lets do some hypothesis check?
	// // do we have the first item in this list added?
	// // if its we probabl
	// first := c.list[0]
	// if strings.Contains(a.Value, first) {
	//  a.Value = list
	// }

	a.Value = fmt.Sprintf("%s %s", a.Value, list)
}

//Apply adds the giving element into the current elements children tree
func (e *Element) Apply(em *Element) {
	if em.allowChildren {
		em.AddChild(e)
	}
}

// Apply applies a set change to the giving element attributes list
func (a *Attribute) Apply(e *Element) {
	if e.allowAttributes {
		e.attrs = append(e.attrs, a)
	}
}

// Apply applies a set change to the giving element style list
func (s *Style) Apply(e *Element) {
	if e.allowStyles {
		e.styles = append(e.styles, s)
	}
}

// Apply adds the event into the elements events lists
func (e *Event) Apply(em *Element) {
	if em.allowEvents {
		if e.Meta.EventTarget == "" {
			e.Meta.EventTarget = em.EventID()
		}
		e.tree = em
		em.events = append(em.events, e)
	}
}

// Clonable defines an interface for objects that can be cloned
type Clonable interface {
	Clone() Markup
}

//Clone replicates the style into a unique instance
func (e *Event) Clone() *Event {
	return &Event{
		Meta: &domevents.EventMetable{EventType: e.Meta.EventType, EventTarget: e.Meta.EventTarget},
		Fx:   e.Fx,
	}
}

// Clone replicates the lists of classnames.
func (c *ClassList) Clone() *ClassList {
	cl := ClassList{
		list: append([]string{}, c.list...),
	}

	return &cl
}

//Clone replicates the style into a unique instance
func (s *Style) Clone() *Style {
	return &Style{Name: s.Name, Value: s.Value}
}

//Clone replicates the attribute into a unique instance
func (a *Attribute) Clone() *Attribute {
	return &Attribute{Name: a.Name, Value: a.Value}
}

// Clone makes a new copy of the markup structure
func (e *Element) Clone() Markup {
	co := NewElement(e.Name(), e.AutoClosed())

	//copy over the textContent
	co.textContent = e.textContent

	//copy over the attribute lockers
	co.allowChildren = e.allowChildren
	co.allowEvents = e.allowEvents
	co.allowAttributes = e.allowAttributes
	co.eventManager = e.eventManager

	if e.Removed() {
		co.Removed()
	}

	//clone the internal styles
	for _, so := range e.styles {
		so.Clone().Apply(co)
	}

	co.allowStyles = e.allowStyles

	//clone the internal attribute
	for _, ao := range e.attrs {
		ao.Clone().Apply(co)
	}

	// co.allowAttributes = e.allowAttributes
	//clone the internal children
	for _, ch := range e.children {
		ch.Clone().Apply(co)
	}

	for _, ch := range e.events {
		ch.Clone().Apply(co)
	}

	return co
}

// Augment adds new markup to an the root if its Element
func Augment(root Markup, m ...Markup) {
	if el, ok := root.(*Element); ok {
		for _, mo := range m {
			mo.Apply(el)
		}
	}
}

// ReconcileEvents checks through two markup events against each other and if it finds any disparity marks
// event objects as Removed
func ReconcileEvents(e, em Markup) {
	oldevents := em.Events()
	newevents := e.Events()

	if len(oldevents) <= 0 && len(newevents) <= 0 {
		return
	}

	if len(newevents) <= 0 && len(oldevents) > 0 {
		for _, ev := range oldevents {
			ev.Meta.Remove()
		}
		return
	}

	checkOut := func(ev *Event) bool {
		for _, evs := range newevents {
			if evs.Meta.EventType == ev.Meta.EventType {
				return true
			}
		}
		return false

	}
	//set to equal as the logic will try to assert its falsiness

	// outerfind:
	for _, ev := range oldevents {
		if checkOut(ev) {
			continue
		}

		ev.Meta.Remove()
	}

}

// EqualAttributes returns true/false if the elements and the giving markup have equal attribute
func EqualAttributes(e, em Markup) bool {
	oldAttrs := em.Attributes()

	if len(oldAttrs) <= 0 {
		if len(e.Attributes()) <= 0 {
			return true
		}
		return false
	}

	//set to equal as the logic will try to assert its falsiness
	var equal = true

	for _, oa := range oldAttrs {
		//lets get the attribute type from the element, if it exists then check the value if its equal
		// continue the loop and check the rest, else we found a contention point, attribute of old markup
		// does not exists in new markup, so we break and mark as different,letting the new markup keep its hash
		// but if the loop finishes and all are equal then we swap the hashes
		if ta, err := GetAttr(e, oa.Name); err == nil {
			if ta.Value == oa.Value {
				continue
			}

			equal = false
			break
		} else {
			equal = false
			break
		}
	}

	return equal
}

// GetStyles returns the styles that contain the specified name and if not empty that contains the specified value also, note that strings
// NOTE: string.Contains is used when checking value parameter if present
func GetStyles(e Markup, f, val string) []*Style {
	var found []*Style
	var styles = e.Styles()

	for _, as := range styles {
		if as.Name != f {
			continue
		}

		if val != "" {
			if !strings.Contains(as.Value, val) {
				continue
			}
		}

		found = append(found, as)
	}

	return found
}

// GetStyle returns the style with the specified tag name
func GetStyle(e Markup, f string) (*Style, error) {
	styles := e.Styles()
	for _, as := range styles {
		if as.Name == f {
			return as, nil
		}
	}
	return nil, ErrNotFound
}

// StyleContains returns the styles that contain the specified name and if the val is not empty then
// that contains the specified value also, note that strings
// NOTE: string.Contains is used
func StyleContains(e Markup, f, val string) bool {
	styles := e.Styles()
	for _, as := range styles {
		if !strings.Contains(as.Name, f) {
			continue
		}

		if val != "" {
			if !strings.Contains(as.Value, val) {
				continue
			}
		}

		return true
	}

	return false
}

// GetAttrs returns the attributes that have the specified text within the naming
// convention and if it also contains the set val if not an empty "",
// NOTE: string.Contains is used
func GetAttrs(e Markup, f, val string) []*Attribute {

	var found []*Attribute

	for _, as := range e.Attributes() {
		if as.Name != f {
			continue
		}

		if val != "" {
			if !strings.Contains(as.Value, val) {
				continue
			}
		}

		found = append(found, as)
	}

	return found
}

// AttrContains returns the attributes that have the specified text within the naming
// convention and if it also contains the set val if not an empty "",
// NOTE: string.Contains is used
func AttrContains(e Markup, f, val string) bool {
	for _, as := range e.Attributes() {
		if !strings.Contains(as.Name, f) {
			continue
		}

		if val != "" {
			if !strings.Contains(as.Value, val) {
				continue
			}
		}

		return true
	}

	return false
}

// GetAttr returns the attribute with the specified tag name
func GetAttr(e Markup, f string) (*Attribute, error) {
	for _, as := range e.Attributes() {
		if as.Name == f {
			return as, nil
		}
	}
	return nil, ErrNotFound
}

// ElementsUsingStyle returns the children within the element matching the
// stlye restrictions passed.
// NOTE: is uses StyleContains
func ElementsUsingStyle(e Markup, f, val string) []Markup {
	return DeepElementsUsingStyle(e, f, val, 1)
}

// ElementsWithAttr returns the children within the element matching the
// stlye restrictions passed.
// NOTE: is uses AttrContains
func ElementsWithAttr(e Markup, f, val string) []Markup {
	return DeepElementsWithAttr(e, f, val, 1)
}

// DeepElementsUsingStyle returns the children within the element matching the
// style restrictions passed allowing control of search depth
// NOTE: is uses StyleContains
// WARNING: depth must start at 1
func DeepElementsUsingStyle(e Markup, f, val string, depth int) []Markup {
	if depth <= 0 {
		return nil
	}

	var found []Markup

	for _, ch := range e.Children() {
		// if che, ok := ch.(*Element); ok {
		if StyleContains(ch, f, val) {
			found = append(found, ch)
			cfo := DeepElementsUsingStyle(ch, f, val, depth-1)
			if len(cfo) > 0 {
				found = append(found, cfo...)
			}
		}
		// }
	}

	return found
}

// DeepElementsWithAttr returns the children within the element matching the
// attributes restrictions passed allowing control of search depth
// NOTE: is uses Element.AttrContains
// WARNING: depth must start at 1
func DeepElementsWithAttr(e Markup, f, val string, depth int) []Markup {
	if depth <= 0 {
		return nil
	}

	var found []Markup

	for _, ch := range e.Children() {
		// if che, ok := ch.(*Element); ok {
		if AttrContains(ch, f, val) {
			found = append(found, ch)
			cfo := DeepElementsWithAttr(ch, f, val, depth-1)
			if len(cfo) > 0 {
				found = append(found, cfo...)
			}
		}
		// }
	}

	return found
}

// ElementsWithTag returns elements matching the tag type in the parent markup children list
// only without going deeper into children's children lists
func ElementsWithTag(e Markup, f string) []Markup {
	return DeepElementsWithTag(e, f, 1)
}

// DeepElementsWithTag returns elements matching the tag type in the parent markup
// and depending on the depth will walk down other children within the children.
// WARNING: depth must start at 1
func DeepElementsWithTag(e Markup, f string, depth int) []Markup {
	if depth <= 0 {
		return nil
	}

	f = strings.TrimSpace(strings.ToLower(f))

	var found []Markup

	for _, ch := range e.Children() {
		// if che, ok := ch.(*Element); ok {
		if ch.Name() == f {
			found = append(found, ch)
			cfo := DeepElementsWithTag(ch, f, depth-1)
			if len(cfo) > 0 {
				found = append(found, cfo...)
			}
		}
		// }
	}

	return found
}
