package domevents

import "github.com/gopherjs/gopherjs/js"

// JSEventMux represents a js.listener function which is returned when attached
// using AddEventListeners and is used for removals with RemoveEventListeners
type JSEventMux func(*js.Object)

// Event defines the base interface for browser events and defines the basic interface methods they must provide
type Event interface {
	Bubbles() bool
	Cancelable() bool
	CurrentTarget() *js.Object
	DefaultPrevented() bool
	EventPhase() int
	Target() *js.Object
	Timestamp() int
	Type() string
	Core() *js.Object
	StopPropagation()
	StopImmediatePropagation()
	PreventDefault()
}

// EventObject implements the Event interface and is embedded by
// concrete event types.
type EventObject struct {
	*js.Object
}

// Core returns the internal js object for the event
func (ev *EventObject) Core() *js.Object {
	return ev.Object
}

// Bubbles returns true/false if the event can bubble up
func (ev *EventObject) Bubbles() bool {
	return ev.Get("bubbles").Bool()
}

// Cancelable returns true/false if the event is cancelable
func (ev *EventObject) Cancelable() bool {
	return ev.Get("cancelable").Bool()
}

// CurrentTarget returns the current target of the event when received
func (ev *EventObject) CurrentTarget() *js.Object {
	return ev.Get("currentTarget")
}

// DefaultPrevented returns true/false if the event was prevented
func (ev *EventObject) DefaultPrevented() bool {
	return ev.Get("defaultPrevented").Bool()
}

// EventPhase returns the state of the event
func (ev *EventObject) EventPhase() int {
	return ev.Get("eventPhase").Int()
}

// Target returns the target of the event as a js.Object
func (ev *EventObject) Target() *js.Object {
	return ev.Get("target")
}

// Timestamp returns the event timestamp value as a int (in seconds)
func (ev *EventObject) Timestamp() int {
	ms := ev.Get("timeStamp").Int()
	s := ms / 1000
	// ns := (ms % 1000 * 1e6)
	// return time.Unix(int64(s), int64(ns))
	return s
}

// Type returns the event type value
func (ev *EventObject) Type() string {
	return ev.Get("type").String()
}

// PreventDefault prevents the default value of the event
func (ev *EventObject) PreventDefault() {
	ev.Call("preventDefault")
}

// StopImmediatePropagation stops the propagation of the event forcefully
func (ev *EventObject) StopImmediatePropagation() {
	ev.Call("stopImmediatePropagation")
}

// StopPropagation stops the propagation of the event
func (ev *EventObject) StopPropagation() {
	ev.Call("stopPropagation")
}
