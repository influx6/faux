package domevents

import (
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/gopherjs/gopherjs/js"
)

// EventHandler provides the function type for event callbacks when subscribing
// to events in Haiku.
type EventHandler func(Event)

// EventActions defines standard actions for event objects.
type EventActions interface {
	StopPropagation() bool
	StopImmediatePropagation() bool
	PreventDefault() bool
}

// EventMeta provides a basic information about the events and what its targets
type EventMeta interface {
	EventActions
	Type() string
	Target() string
	Removed() bool
	Remove()
}

// EventMetable provides a concrete structural
// implementation of the EventMeta interface.
type EventMetable struct {
	EventType                      string
	EventTarget                    string
	ShouldStopPropagation          bool
	ShouldStopImmediatePropagation bool
	ShouldPreventDefault           bool
	ShouldRemove                   bool
}

// StopImmediatePropagation returns the true/false if event should
// stop immediate propagation.
func (e *EventMetable) StopImmediatePropagation() bool {
	return e.ShouldStopImmediatePropagation
}

// PreventDefault returns the true/false if event should
// prevent its default behaviour.
func (e *EventMetable) PreventDefault() bool {
	return e.ShouldPreventDefault
}

// StopPropagation returns the true/false if event should
// stop propagation.
func (e *EventMetable) StopPropagation() bool {
	return e.ShouldStopPropagation
}

// Remove sets the event meta as removed.
func (e *EventMetable) Remove() {
	e.ShouldRemove = true
}

// Removed returns true/false if the event meta is marked
// removed.
func (e *EventMetable) Removed() bool {
	return e.ShouldRemove
}

// Target returns the target of event.
func (e *EventMetable) Target() string {
	return e.EventTarget
}

// Type returns the type of event.
func (e *EventMetable) Type() string {
	return e.EventType
}

// EventSubs provides a interface for event registeration subscribers.
type EventSubs interface {
	Qu
	EventMeta
	Offload()
	DOM(*js.Object)
}

// EventSubHandler provides a event type callback for
// EventSubs.
type EventSubHandler func(EventSubs)

// EventSub represent a single event configuration for dom.Elem objects
// instance which allows chaining of events listeners like middleware
type EventSub struct {
	EventMeta
	Qu
	jslink JSEventMux
	dom    *js.Object
}

// NewEventSub returns a new event element config
func NewEventSub(evtype, evtarget string) *EventSub {
	return &EventSub{
		EventMeta: &EventMetable{
			EventType:   evtype,
			EventTarget: evtarget,
		},
		Qu: NewQu(),
	}
}

// MetaEventSub returns a new event using the supplied EventMeta
func MetaEventSub(meta EventMeta) *EventSub {
	return &EventSub{
		EventMeta: meta,
		Qu:        NewQu(),
	}
}

// DOM sets up the event subs for listening
func (e *EventSub) DOM(dom *js.Object) {
	e.Offload()
	e.dom = dom
	e.jslink = func(o *js.Object) { e.TriggerMatch(&EventObject{o}) }
	e.dom.Call("addEventListener", e.Type(), e.jslink, true)
}

// Offload removes all event bindings from current dom element
func (e *EventSub) Offload() {
	if e.dom == nil {
		return
	}

	if e.jslink != nil {
		e.dom.Call("removeEventListener", e.Type(), e.jslink, true)
		// e.dom.RemoveEventListener(e.EventType(), true, e.jslink)
		e.jslink = nil
	}
}

// Trigger provides bypass for triggering this event sub by passing down an event
// directly without matching target or selector
func (e *EventSub) Trigger(h Event) {
	e.Run(h)
}

// TriggerMatch check if the current event from a specific parent matches the
// eventarget by using the eventsub selector,if the target is within the results for
// that selector then it triggers the event subscribers
func (e *EventSub) TriggerMatch(h Event) {
	// if e.dom != nil
	if strings.ToLower(h.Type()) != strings.ToLower(e.Type()) {
		return
	}

	//get the current event target
	target := h.Target()

	// log.Printf("target -> %+s", target)

	//get the targets parent
	// parent := target.Get("parentElement")
	parent := e.dom

	var match bool

	children := parent.Call("querySelectorAll", e.Target())

	if children == nil || children == js.Undefined {
		return
	}

	// log.Printf("children -> %s  -> %t %t", children, children == nil, children == js.Undefined)

	//get all possible matches of this query
	// posis := parent.QuerySelectorAll(e.EventSelector())
	posis := DOMObjectToList(children)

	// log.Printf("Checking: %s for %s -> %+s", target, e.ID(), posis)

	//is our target part of those that match the selector
	for _, item := range posis {
		// log.Printf("taget %+s and item %+s -> %t", target, item, target == item)
		if item != target {
			continue
		}
		match = true
		break
	}

	//if we match then run the listeners registered
	if match {
		//if we dont want immediatepropagation kill it else check propagation also
		if e.StopImmediatePropagation() {
			h.StopImmediatePropagation()
		} else {
			//if we dont want propagation kill it
			if e.StopPropagation() {
				h.StopPropagation()
			}
		}

		//we want to PreventDefault then stop default action
		if e.PreventDefault() {
			h.PreventDefault()
		}

		e.Run(h)
	}
}

// ID returns the event id that EventManager use for this event
func (e *EventSub) ID() string {
	return GetEventID(e)
}

// ErrEventNotFound is returned when an event is not found
var ErrEventNotFound = errors.New("Event not found")

// EventManagerHandler provides a event type callback for
// EventManagers.
type EventManagerHandler func(EventManagers)

// EventManagers defines the Event.EventManager interface type and is used to clean up import path usage and standardize the api
type EventManagers interface {
	HasEvent(m string) bool
	GetEvent(event string) (EventSubs, error)
	NewEventMeta(EventMeta) (EventSubs, bool)
	NewEvent(evtype, evselector string) (EventSubs, bool)
	AttachManager(EventManagers)
	DetachManager(EventManagers)
	HasManager(EventManagers) bool
	RemoveEvent(event string)
	AddEvent(EventSubs) bool
	AddEvents(...EventSubs)
	EachEvent(EventSubHandler)
	EachManager(EventManagerHandler)
	DisconnectRemoved()
	OffloadDOM()
	LoadUpEvents()
	LoadDOM(dom *js.Object) bool
}

// EventManager provides a deffered event managing sytem for registery events with
type EventManager struct {
	events   map[string]EventSubs
	attaches map[EventManagers]bool
	ro       sync.RWMutex
	wo       sync.RWMutex
	dom      *js.Object
}

// NewEventManager returns a new event manager instance
func NewEventManager() EventManagers {
	em := EventManager{
		events:   make(map[string]EventSubs),
		attaches: make(map[EventManagers]bool),
	}

	return &em
}

// HasEvent returns true/false if an event target is already marked using the
// format selector#eventType
func (em *EventManager) HasEvent(m string) bool {
	var ok bool
	em.ro.RLock()
	_, ok = em.events[m]
	em.ro.RUnlock()
	return ok
}

// GetEvent returns the event if found by that id
func (em *EventManager) GetEvent(event string) (EventSubs, error) {
	if !em.HasEvent(event) {
		return nil, ErrEventNotFound
	}

	var ed EventSubs

	em.ro.RLock()
	ed = em.events[event]
	em.ro.RUnlock()

	return ed, nil
}

// NewEventMeta allows the adding of event using string values
func (em *EventManager) NewEventMeta(meta EventMeta) (EventSubs, bool) {
	if meta.Removed() {
		return nil, false
	}

	eo := BuildEventID(meta.Type(), meta.Target())

	if em.HasEvent(eo) {
		ed, _ := em.GetEvent(eo)
		return ed, false
	}

	emo := MetaEventSub(meta)

	em.AddEvent(emo)
	return emo, true
}

// NewEvent allows the adding of event using string values
func (em *EventManager) NewEvent(evtype, evselector string) (EventSubs, bool) {
	eo := BuildEventID(evtype, evselector)

	if em.HasEvent(eo) {
		ed, _ := em.GetEvent(eo)
		return ed, false
	}

	emo := NewEventSub(evtype, evselector)

	em.AddEvent(emo)
	return emo, true
}

// AttachManager allows a manager to get attached to another manajor to receive a dom binding
// when this receives one
func (em *EventManager) AttachManager(esm EventManagers) {
	//incase of stupid loops, are we attached to the supplied manager? if so duck this
	if esm.HasManager(em) {
		return
	}

	//if we have it already we skip
	if em.HasManager(esm) {
		return
	}

	//its not found so we add it
	em.wo.Lock()
	em.attaches[esm] = true
	em.wo.Unlock()

	//do we already have a dom attached?, then notify this manager immediately
	if em.dom != nil {
		esm.LoadDOM(em.dom)
	}
}

// DetachManager detaches the manager if attached already
func (em *EventManager) DetachManager(esm EventManagers) {
	//if we dont have it attached then skip
	if !em.HasManager(esm) {
		return
	}

	//we got one so we kill it
	em.wo.Lock()
	delete(em.attaches, esm)
	em.wo.Unlock()
}

// HasManager returns true/false if a manager was already attached
func (em *EventManager) HasManager(esm EventManagers) bool {
	//have we already attached before
	em.wo.RLock()
	ok := em.attaches[esm]
	em.wo.RUnlock()
	return ok
}

// RemoveEvent removes a event from the list
func (em *EventManager) RemoveEvent(event string) {
	if !em.HasEvent(event) {
		return
	}

	ev, _ := em.GetEvent(event)
	ev.Offload()

	em.wo.Lock()
	delete(em.events, event)
	em.wo.Unlock()
}

//AddEvent adds Event elements into the event manager if a
//event element is already matching that is the combination of selector#eventtype
// then it returns false but if added then true
func (em *EventManager) AddEvent(eo EventSubs) bool {
	id := GetEventID(eo)

	if em.HasEvent(id) {
		return false
	}

	em.ro.Lock()
	em.events[id] = eo
	em.ro.Unlock()
	return true
}

// AddEvents adds a set of ElemEvent into the EventManager
func (em *EventManager) AddEvents(ems ...EventSubs) {
	for _, eo := range ems {
		em.AddEvent(eo)
	}
}

// EachEvent runnings a function over all events
func (em *EventManager) EachEvent(fx EventSubHandler) {
	em.ro.Lock()
	for _, eo := range em.events {
		fx(eo)
	}
	em.ro.Unlock()
}

// EachManager runnings a function over all attached managers
func (em *EventManager) EachManager(fx EventManagerHandler) {
	em.wo.Lock()
	for eo := range em.attaches {
		fx(eo)
	}
	em.wo.Unlock()
}

// DisconnectRemoved disconnects all events that must be removed and removes them
func (em *EventManager) DisconnectRemoved() {
	// send the dom out to all registered event subs for loadup
	em.EachEvent(func(es EventSubs) {
		if es.Removed() {
			// es.Offload()
			em.RemoveEvent(GetEventID(es))
		}
	})
}

// OffloadDOM deregisters the dom and offloads its events to allow other dom to be attached.
// Must call this first before try to use LoadDOM if the EventManager already is loaded
func (em *EventManager) OffloadDOM() {
	if em.dom == nil {
		return
	}
	// send the dom out to all registered event subs for loadup
	em.EachEvent(func(es EventSubs) {
		es.Offload()
	})

	em.dom = nil
}

// LoadUpEvents registers the events into the dom object
func (em *EventManager) LoadUpEvents() {
	if em.dom == nil {
		return
	}

	dom := em.dom

	// send the dom out to all registered event subs for loadup
	em.EachEvent(func(es EventSubs) {
		if !es.Removed() {
			es.DOM(dom)
		}
	})

	// send out to all other attach eventmanagers for loadup
	em.EachManager(func(ems EventManagers) {
		ems.LoadDOM(dom)
	})

}

// LoadDOM passes down the dom element to all EventSub to initialize and listen for their respective events
func (em *EventManager) LoadDOM(dom *js.Object) bool {
	if em.dom != nil {
		return false
	}

	//replace the current dom node be used
	em.dom = dom
	em.LoadUpEvents()

	return true
}

// GetEventID returns the id for a ElemEvent object
func GetEventID(m EventSubs) string {
	sel := strings.TrimSpace(m.Target())
	return BuildEventID(sel, m.Type())
}

// BuildEventID returns the string represent of the values using the select#event format
func BuildEventID(etype, eselect string) string {
	return fmt.Sprintf("%s#%s", eselect, etype)
}
