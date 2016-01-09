// The generation of this package was inspired by Neelance work on DOM (https://github.com/neelance/dom)

//go:generate go run generate.go

// Documentation source: "Event reference" by Mozilla Contributors, https://developer.mozilla.org/en-US/docs/Web/Events, licensed under CC-BY-SA 2.5.

//Package events defines the event binding system for Haiku(https://github.com/influx6/Haiku)
package events

import (
	"github.com/influx6/faux/domtrees"
)

// Abort Documentation is as below:
// A transaction has been aborted.
// https://developer.mozilla.org/docs/Web/Reference/Events/abort_indexedDB
/* This event provides options() to be called when the events is triggered and an optional selector which will override the internal selector mechanism of the domtrees.Element i.e if the selectorOverride argument is an empty string then domtrees.Element will create an appropriate selector matching its type and uid value in this format  (ElementType[uid='UID_VALUE']) but if the selector value is not empty then that becomes the default selector used
match the event with. */
func Abort(fx domtrees.EventHandler, selectorOverride string) *domtrees.Event {
	return domtrees.NewEvent("abort", selectorOverride, fx)
}

// AfterPrint Documentation is as below:
// The associated document has started printing or the print preview has been closed.
// https://developer.mozilla.org/docs/Web/Events/afterprint
/* This event provides options() to be called when the events is triggered and an optional selector which will override the internal selector mechanism of the domtrees.Element i.e if the selectorOverride argument is an empty string then domtrees.Element will create an appropriate selector matching its type and uid value in this format  (ElementType[uid='UID_VALUE']) but if the selector value is not empty then that becomes the default selector used
match the event with. */
func AfterPrint(fx domtrees.EventHandler, selectorOverride string) *domtrees.Event {
	return domtrees.NewEvent("afterprint", selectorOverride, fx)
}

// AnimationEnd Documentation is as below:
// A CSS animation has completed.
// https://developer.mozilla.org/docs/Web/Events/animationend
/* This event provides options() to be called when the events is triggered and an optional selector which will override the internal selector mechanism of the domtrees.Element i.e if the selectorOverride argument is an empty string then domtrees.Element will create an appropriate selector matching its type and uid value in this format  (ElementType[uid='UID_VALUE']) but if the selector value is not empty then that becomes the default selector used
match the event with. */
func AnimationEnd(fx domtrees.EventHandler, selectorOverride string) *domtrees.Event {
	return domtrees.NewEvent("animationend", selectorOverride, fx)
}

// AnimationIteration Documentation is as below:
// A CSS animation is repeated.
// https://developer.mozilla.org/docs/Web/Events/animationiteration
/* This event provides options() to be called when the events is triggered and an optional selector which will override the internal selector mechanism of the domtrees.Element i.e if the selectorOverride argument is an empty string then domtrees.Element will create an appropriate selector matching its type and uid value in this format  (ElementType[uid='UID_VALUE']) but if the selector value is not empty then that becomes the default selector used
match the event with. */
func AnimationIteration(fx domtrees.EventHandler, selectorOverride string) *domtrees.Event {
	return domtrees.NewEvent("animationiteration", selectorOverride, fx)
}

// AnimationStart Documentation is as below:
// A CSS animation has started.
// https://developer.mozilla.org/docs/Web/Events/animationstart
/* This event provides options() to be called when the events is triggered and an optional selector which will override the internal selector mechanism of the domtrees.Element i.e if the selectorOverride argument is an empty string then domtrees.Element will create an appropriate selector matching its type and uid value in this format  (ElementType[uid='UID_VALUE']) but if the selector value is not empty then that becomes the default selector used
match the event with. */
func AnimationStart(fx domtrees.EventHandler, selectorOverride string) *domtrees.Event {
	return domtrees.NewEvent("animationstart", selectorOverride, fx)
}

// AudioProcess Documentation is as below:
// The input buffer of a ScriptProcessorNode is ready to be processed.
// https://developer.mozilla.org/docs/Web/Events/audioprocess
/* This event provides options() to be called when the events is triggered and an optional selector which will override the internal selector mechanism of the domtrees.Element i.e if the selectorOverride argument is an empty string then domtrees.Element will create an appropriate selector matching its type and uid value in this format  (ElementType[uid='UID_VALUE']) but if the selector value is not empty then that becomes the default selector used
match the event with. */
func AudioProcess(fx domtrees.EventHandler, selectorOverride string) *domtrees.Event {
	return domtrees.NewEvent("audioprocess", selectorOverride, fx)
}

// Audioend Documentation is as below:
// The user agent has finished capturing audio for speech recognition.
// https://developer.mozilla.org/docs/Web/Events/audioend
/* This event provides options() to be called when the events is triggered and an optional selector which will override the internal selector mechanism of the domtrees.Element i.e if the selectorOverride argument is an empty string then domtrees.Element will create an appropriate selector matching its type and uid value in this format  (ElementType[uid='UID_VALUE']) but if the selector value is not empty then that becomes the default selector used
match the event with. */
func Audioend(fx domtrees.EventHandler, selectorOverride string) *domtrees.Event {
	return domtrees.NewEvent("audioend", selectorOverride, fx)
}

// Audiostart Documentation is as below:
// The user agent has started to capture audio for speech recognition.
// https://developer.mozilla.org/docs/Web/Events/audiostart
/* This event provides options() to be called when the events is triggered and an optional selector which will override the internal selector mechanism of the domtrees.Element i.e if the selectorOverride argument is an empty string then domtrees.Element will create an appropriate selector matching its type and uid value in this format  (ElementType[uid='UID_VALUE']) but if the selector value is not empty then that becomes the default selector used
match the event with. */
func Audiostart(fx domtrees.EventHandler, selectorOverride string) *domtrees.Event {
	return domtrees.NewEvent("audiostart", selectorOverride, fx)
}

// BeforePrint Documentation is as below:
// The associated document is about to be printed or previewed for printing.
// https://developer.mozilla.org/docs/Web/Events/beforeprint
/* This event provides options() to be called when the events is triggered and an optional selector which will override the internal selector mechanism of the domtrees.Element i.e if the selectorOverride argument is an empty string then domtrees.Element will create an appropriate selector matching its type and uid value in this format  (ElementType[uid='UID_VALUE']) but if the selector value is not empty then that becomes the default selector used
match the event with. */
func BeforePrint(fx domtrees.EventHandler, selectorOverride string) *domtrees.Event {
	return domtrees.NewEvent("beforeprint", selectorOverride, fx)
}

// BeforeUnload Documentation is as below:
// (no documentation)
// https://developer.mozilla.org/docs/Web/Events/beforeunload
/* This event provides options() to be called when the events is triggered and an optional selector which will override the internal selector mechanism of the domtrees.Element i.e if the selectorOverride argument is an empty string then domtrees.Element will create an appropriate selector matching its type and uid value in this format  (ElementType[uid='UID_VALUE']) but if the selector value is not empty then that becomes the default selector used
match the event with. */
func BeforeUnload(fx domtrees.EventHandler, selectorOverride string) *domtrees.Event {
	return domtrees.NewEvent("beforeunload", selectorOverride, fx)
}

// BeginEvent Documentation is as below:
// A SMIL animation element begins.
// https://developer.mozilla.org/docs/Web/Events/beginEvent
/* This event provides options() to be called when the events is triggered and an optional selector which will override the internal selector mechanism of the domtrees.Element i.e if the selectorOverride argument is an empty string then domtrees.Element will create an appropriate selector matching its type and uid value in this format  (ElementType[uid='UID_VALUE']) but if the selector value is not empty then that becomes the default selector used
match the event with. */
func BeginEvent(fx domtrees.EventHandler, selectorOverride string) *domtrees.Event {
	return domtrees.NewEvent("beginEvent", selectorOverride, fx)
}

// Blocked Documentation is as below:
// An open connection to a database is blocking a versionchange transaction on the same database.
// https://developer.mozilla.org/docs/Web/Reference/Events/blocked_indexedDB
/* This event provides options() to be called when the events is triggered and an optional selector which will override the internal selector mechanism of the domtrees.Element i.e if the selectorOverride argument is an empty string then domtrees.Element will create an appropriate selector matching its type and uid value in this format  (ElementType[uid='UID_VALUE']) but if the selector value is not empty then that becomes the default selector used
match the event with. */
func Blocked(fx domtrees.EventHandler, selectorOverride string) *domtrees.Event {
	return domtrees.NewEvent("blocked", selectorOverride, fx)
}

// Blur Documentation is as below:
// An element has lost focus (does not bubble).
// https://developer.mozilla.org/docs/Web/Events/blur
/* This event provides options() to be called when the events is triggered and an optional selector which will override the internal selector mechanism of the domtrees.Element i.e if the selectorOverride argument is an empty string then domtrees.Element will create an appropriate selector matching its type and uid value in this format  (ElementType[uid='UID_VALUE']) but if the selector value is not empty then that becomes the default selector used
match the event with. */
func Blur(fx domtrees.EventHandler, selectorOverride string) *domtrees.Event {
	return domtrees.NewEvent("blur", selectorOverride, fx)
}

// Boundary Documentation is as below:
// The spoken utterance reaches a word or sentence boundary
// https://developer.mozilla.org/docs/Web/Events/boundary
/* This event provides options() to be called when the events is triggered and an optional selector which will override the internal selector mechanism of the domtrees.Element i.e if the selectorOverride argument is an empty string then domtrees.Element will create an appropriate selector matching its type and uid value in this format  (ElementType[uid='UID_VALUE']) but if the selector value is not empty then that becomes the default selector used
match the event with. */
func Boundary(fx domtrees.EventHandler, selectorOverride string) *domtrees.Event {
	return domtrees.NewEvent("boundary", selectorOverride, fx)
}

// Cached Documentation is as below:
// The resources listed in the manifest have been downloaded, and the application is now cached.
// https://developer.mozilla.org/docs/Web/Events/cached
/* This event provides options() to be called when the events is triggered and an optional selector which will override the internal selector mechanism of the domtrees.Element i.e if the selectorOverride argument is an empty string then domtrees.Element will create an appropriate selector matching its type and uid value in this format  (ElementType[uid='UID_VALUE']) but if the selector value is not empty then that becomes the default selector used
match the event with. */
func Cached(fx domtrees.EventHandler, selectorOverride string) *domtrees.Event {
	return domtrees.NewEvent("cached", selectorOverride, fx)
}

// CanPlay Documentation is as below:
// The user agent can play the media, but estimates that not enough data has been loaded to play the media up to its end without having to stop for further buffering of content.
// https://developer.mozilla.org/docs/Web/Events/canplay
/* This event provides options() to be called when the events is triggered and an optional selector which will override the internal selector mechanism of the domtrees.Element i.e if the selectorOverride argument is an empty string then domtrees.Element will create an appropriate selector matching its type and uid value in this format  (ElementType[uid='UID_VALUE']) but if the selector value is not empty then that becomes the default selector used
match the event with. */
func CanPlay(fx domtrees.EventHandler, selectorOverride string) *domtrees.Event {
	return domtrees.NewEvent("canplay", selectorOverride, fx)
}

// CanPlayThrough Documentation is as below:
// The user agent can play the media, and estimates that enough data has been loaded to play the media up to its end without having to stop for further buffering of content.
// https://developer.mozilla.org/docs/Web/Events/canplaythrough
/* This event provides options() to be called when the events is triggered and an optional selector which will override the internal selector mechanism of the domtrees.Element i.e if the selectorOverride argument is an empty string then domtrees.Element will create an appropriate selector matching its type and uid value in this format  (ElementType[uid='UID_VALUE']) but if the selector value is not empty then that becomes the default selector used
match the event with. */
func CanPlayThrough(fx domtrees.EventHandler, selectorOverride string) *domtrees.Event {
	return domtrees.NewEvent("canplaythrough", selectorOverride, fx)
}

// Change Documentation is as below:
// An element loses focus and its value changed since gaining focus.
// https://developer.mozilla.org/docs/Web/Events/change
/* This event provides options() to be called when the events is triggered and an optional selector which will override the internal selector mechanism of the domtrees.Element i.e if the selectorOverride argument is an empty string then domtrees.Element will create an appropriate selector matching its type and uid value in this format  (ElementType[uid='UID_VALUE']) but if the selector value is not empty then that becomes the default selector used
match the event with. */
func Change(fx domtrees.EventHandler, selectorOverride string) *domtrees.Event {
	return domtrees.NewEvent("change", selectorOverride, fx)
}

// ChargingChange Documentation is as below:
// The battery begins or stops charging.
// https://developer.mozilla.org/docs/Web/Events/chargingchange
/* This event provides options() to be called when the events is triggered and an optional selector which will override the internal selector mechanism of the domtrees.Element i.e if the selectorOverride argument is an empty string then domtrees.Element will create an appropriate selector matching its type and uid value in this format  (ElementType[uid='UID_VALUE']) but if the selector value is not empty then that becomes the default selector used
match the event with. */
func ChargingChange(fx domtrees.EventHandler, selectorOverride string) *domtrees.Event {
	return domtrees.NewEvent("chargingchange", selectorOverride, fx)
}

// ChargingTimeChange Documentation is as below:
// The chargingTime attribute has been updated.
// https://developer.mozilla.org/docs/Web/Events/chargingtimechange
/* This event provides options() to be called when the events is triggered and an optional selector which will override the internal selector mechanism of the domtrees.Element i.e if the selectorOverride argument is an empty string then domtrees.Element will create an appropriate selector matching its type and uid value in this format  (ElementType[uid='UID_VALUE']) but if the selector value is not empty then that becomes the default selector used
match the event with. */
func ChargingTimeChange(fx domtrees.EventHandler, selectorOverride string) *domtrees.Event {
	return domtrees.NewEvent("chargingtimechange", selectorOverride, fx)
}

// Checking Documentation is as below:
// The user agent is checking for an update, or attempting to download the cache manifest for the first time.
// https://developer.mozilla.org/docs/Web/Events/checking
/* This event provides options() to be called when the events is triggered and an optional selector which will override the internal selector mechanism of the domtrees.Element i.e if the selectorOverride argument is an empty string then domtrees.Element will create an appropriate selector matching its type and uid value in this format  (ElementType[uid='UID_VALUE']) but if the selector value is not empty then that becomes the default selector used
match the event with. */
func Checking(fx domtrees.EventHandler, selectorOverride string) *domtrees.Event {
	return domtrees.NewEvent("checking", selectorOverride, fx)
}

// Click Documentation is as below:
// A pointing device button has been pressed and released on an element.
// https://developer.mozilla.org/docs/Web/Events/click
/* This event provides options() to be called when the events is triggered and an optional selector which will override the internal selector mechanism of the domtrees.Element i.e if the selectorOverride argument is an empty string then domtrees.Element will create an appropriate selector matching its type and uid value in this format  (ElementType[uid='UID_VALUE']) but if the selector value is not empty then that becomes the default selector used
match the event with. */
func Click(fx domtrees.EventHandler, selectorOverride string) *domtrees.Event {
	return domtrees.NewEvent("click", selectorOverride, fx)
}

// Close Documentation is as below:
// A WebSocket connection has been closed.
// https://developer.mozilla.org/docs/Web/Reference/Events/close_websocket
/* This event provides options() to be called when the events is triggered and an optional selector which will override the internal selector mechanism of the domtrees.Element i.e if the selectorOverride argument is an empty string then domtrees.Element will create an appropriate selector matching its type and uid value in this format  (ElementType[uid='UID_VALUE']) but if the selector value is not empty then that becomes the default selector used
match the event with. */
func Close(fx domtrees.EventHandler, selectorOverride string) *domtrees.Event {
	return domtrees.NewEvent("close", selectorOverride, fx)
}

// Complete Documentation is as below:
// The rendering of an OfflineAudioContext is terminated.
// https://developer.mozilla.org/docs/Web/Events/complete
/* This event provides options() to be called when the events is triggered and an optional selector which will override the internal selector mechanism of the domtrees.Element i.e if the selectorOverride argument is an empty string then domtrees.Element will create an appropriate selector matching its type and uid value in this format  (ElementType[uid='UID_VALUE']) but if the selector value is not empty then that becomes the default selector used
match the event with. */
func Complete(fx domtrees.EventHandler, selectorOverride string) *domtrees.Event {
	return domtrees.NewEvent("complete", selectorOverride, fx)
}

// CompositionEnd Documentation is as below:
// The composition of a passage of text has been completed or canceled.
// https://developer.mozilla.org/docs/Web/Events/compositionend
/* This event provides options() to be called when the events is triggered and an optional selector which will override the internal selector mechanism of the domtrees.Element i.e if the selectorOverride argument is an empty string then domtrees.Element will create an appropriate selector matching its type and uid value in this format  (ElementType[uid='UID_VALUE']) but if the selector value is not empty then that becomes the default selector used
match the event with. */
func CompositionEnd(fx domtrees.EventHandler, selectorOverride string) *domtrees.Event {
	return domtrees.NewEvent("compositionend", selectorOverride, fx)
}

// CompositionStart Documentation is as below:
// The composition of a passage of text is prepared (similar to keydown for a keyboard input, but works with other inputs such as speech recognition).
// https://developer.mozilla.org/docs/Web/Events/compositionstart
/* This event provides options() to be called when the events is triggered and an optional selector which will override the internal selector mechanism of the domtrees.Element i.e if the selectorOverride argument is an empty string then domtrees.Element will create an appropriate selector matching its type and uid value in this format  (ElementType[uid='UID_VALUE']) but if the selector value is not empty then that becomes the default selector used
match the event with. */
func CompositionStart(fx domtrees.EventHandler, selectorOverride string) *domtrees.Event {
	return domtrees.NewEvent("compositionstart", selectorOverride, fx)
}

// CompositionUpdate Documentation is as below:
// A character is added to a passage of text being composed.
// https://developer.mozilla.org/docs/Web/Events/compositionupdate
/* This event provides options() to be called when the events is triggered and an optional selector which will override the internal selector mechanism of the domtrees.Element i.e if the selectorOverride argument is an empty string then domtrees.Element will create an appropriate selector matching its type and uid value in this format  (ElementType[uid='UID_VALUE']) but if the selector value is not empty then that becomes the default selector used
match the event with. */
func CompositionUpdate(fx domtrees.EventHandler, selectorOverride string) *domtrees.Event {
	return domtrees.NewEvent("compositionupdate", selectorOverride, fx)
}

// ContextMenu Documentation is as below:
// The right button of the mouse is clicked (before the context menu is displayed).
// https://developer.mozilla.org/docs/Web/Events/contextmenu
/* This event provides options() to be called when the events is triggered and an optional selector which will override the internal selector mechanism of the domtrees.Element i.e if the selectorOverride argument is an empty string then domtrees.Element will create an appropriate selector matching its type and uid value in this format  (ElementType[uid='UID_VALUE']) but if the selector value is not empty then that becomes the default selector used
match the event with. */
func ContextMenu(fx domtrees.EventHandler, selectorOverride string) *domtrees.Event {
	return domtrees.NewEvent("contextmenu", selectorOverride, fx)
}

// Copy Documentation is as below:
// The text selection has been added to the clipboard.
// https://developer.mozilla.org/docs/Web/Events/copy
/* This event provides options() to be called when the events is triggered and an optional selector which will override the internal selector mechanism of the domtrees.Element i.e if the selectorOverride argument is an empty string then domtrees.Element will create an appropriate selector matching its type and uid value in this format  (ElementType[uid='UID_VALUE']) but if the selector value is not empty then that becomes the default selector used
match the event with. */
func Copy(fx domtrees.EventHandler, selectorOverride string) *domtrees.Event {
	return domtrees.NewEvent("copy", selectorOverride, fx)
}

// Cut Documentation is as below:
// The text selection has been removed from the document and added to the clipboard.
// https://developer.mozilla.org/docs/Web/Events/cut
/* This event provides options() to be called when the events is triggered and an optional selector which will override the internal selector mechanism of the domtrees.Element i.e if the selectorOverride argument is an empty string then domtrees.Element will create an appropriate selector matching its type and uid value in this format  (ElementType[uid='UID_VALUE']) but if the selector value is not empty then that becomes the default selector used
match the event with. */
func Cut(fx domtrees.EventHandler, selectorOverride string) *domtrees.Event {
	return domtrees.NewEvent("cut", selectorOverride, fx)
}

// DOMContentLoaded Documentation is as below:
// The document has finished loading (but not its dependent resources).
// https://developer.mozilla.org/docs/Web/Events/DOMContentLoaded
/* This event provides options() to be called when the events is triggered and an optional selector which will override the internal selector mechanism of the domtrees.Element i.e if the selectorOverride argument is an empty string then domtrees.Element will create an appropriate selector matching its type and uid value in this format  (ElementType[uid='UID_VALUE']) but if the selector value is not empty then that becomes the default selector used
match the event with. */
func DOMContentLoaded(fx domtrees.EventHandler, selectorOverride string) *domtrees.Event {
	return domtrees.NewEvent("DOMContentLoaded", selectorOverride, fx)
}

// DblClick Documentation is as below:
// A pointing device button is clicked twice on an element.
// https://developer.mozilla.org/docs/Web/Events/dblclick
/* This event provides options() to be called when the events is triggered and an optional selector which will override the internal selector mechanism of the domtrees.Element i.e if the selectorOverride argument is an empty string then domtrees.Element will create an appropriate selector matching its type and uid value in this format  (ElementType[uid='UID_VALUE']) but if the selector value is not empty then that becomes the default selector used
match the event with. */
func DblClick(fx domtrees.EventHandler, selectorOverride string) *domtrees.Event {
	return domtrees.NewEvent("dblclick", selectorOverride, fx)
}

// DeviceLight Documentation is as below:
// Fresh data is available from a light sensor.
// https://developer.mozilla.org/docs/Web/Events/devicelight
/* This event provides options() to be called when the events is triggered and an optional selector which will override the internal selector mechanism of the domtrees.Element i.e if the selectorOverride argument is an empty string then domtrees.Element will create an appropriate selector matching its type and uid value in this format  (ElementType[uid='UID_VALUE']) but if the selector value is not empty then that becomes the default selector used
match the event with. */
func DeviceLight(fx domtrees.EventHandler, selectorOverride string) *domtrees.Event {
	return domtrees.NewEvent("devicelight", selectorOverride, fx)
}

// DeviceMotion Documentation is as below:
// Fresh data is available from a motion sensor.
// https://developer.mozilla.org/docs/Web/Events/devicemotion
/* This event provides options() to be called when the events is triggered and an optional selector which will override the internal selector mechanism of the domtrees.Element i.e if the selectorOverride argument is an empty string then domtrees.Element will create an appropriate selector matching its type and uid value in this format  (ElementType[uid='UID_VALUE']) but if the selector value is not empty then that becomes the default selector used
match the event with. */
func DeviceMotion(fx domtrees.EventHandler, selectorOverride string) *domtrees.Event {
	return domtrees.NewEvent("devicemotion", selectorOverride, fx)
}

// DeviceOrientation Documentation is as below:
// Fresh data is available from an orientation sensor.
// https://developer.mozilla.org/docs/Web/Events/deviceorientation
/* This event provides options() to be called when the events is triggered and an optional selector which will override the internal selector mechanism of the domtrees.Element i.e if the selectorOverride argument is an empty string then domtrees.Element will create an appropriate selector matching its type and uid value in this format  (ElementType[uid='UID_VALUE']) but if the selector value is not empty then that becomes the default selector used
match the event with. */
func DeviceOrientation(fx domtrees.EventHandler, selectorOverride string) *domtrees.Event {
	return domtrees.NewEvent("deviceorientation", selectorOverride, fx)
}

// DeviceProximity Documentation is as below:
// Fresh data is available from a proximity sensor (indicates an approximated distance between the device and a nearby object).
// https://developer.mozilla.org/docs/Web/Events/deviceproximity
/* This event provides options() to be called when the events is triggered and an optional selector which will override the internal selector mechanism of the domtrees.Element i.e if the selectorOverride argument is an empty string then domtrees.Element will create an appropriate selector matching its type and uid value in this format  (ElementType[uid='UID_VALUE']) but if the selector value is not empty then that becomes the default selector used
match the event with. */
func DeviceProximity(fx domtrees.EventHandler, selectorOverride string) *domtrees.Event {
	return domtrees.NewEvent("deviceproximity", selectorOverride, fx)
}

// DischargingTimeChange Documentation is as below:
// The dischargingTime attribute has been updated.
// https://developer.mozilla.org/docs/Web/Events/dischargingtimechange
/* This event provides options() to be called when the events is triggered and an optional selector which will override the internal selector mechanism of the domtrees.Element i.e if the selectorOverride argument is an empty string then domtrees.Element will create an appropriate selector matching its type and uid value in this format  (ElementType[uid='UID_VALUE']) but if the selector value is not empty then that becomes the default selector used
match the event with. */
func DischargingTimeChange(fx domtrees.EventHandler, selectorOverride string) *domtrees.Event {
	return domtrees.NewEvent("dischargingtimechange", selectorOverride, fx)
}

// Downloading Documentation is as below:
// The user agent has found an update and is fetching it, or is downloading the resources listed by the cache manifest for the first time.
// https://developer.mozilla.org/docs/Web/Events/downloading
/* This event provides options() to be called when the events is triggered and an optional selector which will override the internal selector mechanism of the domtrees.Element i.e if the selectorOverride argument is an empty string then domtrees.Element will create an appropriate selector matching its type and uid value in this format  (ElementType[uid='UID_VALUE']) but if the selector value is not empty then that becomes the default selector used
match the event with. */
func Downloading(fx domtrees.EventHandler, selectorOverride string) *domtrees.Event {
	return domtrees.NewEvent("downloading", selectorOverride, fx)
}

// Drag Documentation is as below:
// An element or text selection is being dragged (every 350ms).
// https://developer.mozilla.org/docs/Web/Events/drag
/* This event provides options() to be called when the events is triggered and an optional selector which will override the internal selector mechanism of the domtrees.Element i.e if the selectorOverride argument is an empty string then domtrees.Element will create an appropriate selector matching its type and uid value in this format  (ElementType[uid='UID_VALUE']) but if the selector value is not empty then that becomes the default selector used
match the event with. */
func Drag(fx domtrees.EventHandler, selectorOverride string) *domtrees.Event {
	return domtrees.NewEvent("drag", selectorOverride, fx)
}

// DragEnd Documentation is as below:
// A drag operation is being ended (by releasing a mouse button or hitting the escape key).
// https://developer.mozilla.org/docs/Web/Events/dragend
/* This event provides options() to be called when the events is triggered and an optional selector which will override the internal selector mechanism of the domtrees.Element i.e if the selectorOverride argument is an empty string then domtrees.Element will create an appropriate selector matching its type and uid value in this format  (ElementType[uid='UID_VALUE']) but if the selector value is not empty then that becomes the default selector used
match the event with. */
func DragEnd(fx domtrees.EventHandler, selectorOverride string) *domtrees.Event {
	return domtrees.NewEvent("dragend", selectorOverride, fx)
}

// DragEnter Documentation is as below:
// A dragged element or text selection enters a valid drop target.
// https://developer.mozilla.org/docs/Web/Events/dragenter
/* This event provides options() to be called when the events is triggered and an optional selector which will override the internal selector mechanism of the domtrees.Element i.e if the selectorOverride argument is an empty string then domtrees.Element will create an appropriate selector matching its type and uid value in this format  (ElementType[uid='UID_VALUE']) but if the selector value is not empty then that becomes the default selector used
match the event with. */
func DragEnter(fx domtrees.EventHandler, selectorOverride string) *domtrees.Event {
	return domtrees.NewEvent("dragenter", selectorOverride, fx)
}

// DragLeave Documentation is as below:
// A dragged element or text selection leaves a valid drop target.
// https://developer.mozilla.org/docs/Web/Events/dragleave
/* This event provides options() to be called when the events is triggered and an optional selector which will override the internal selector mechanism of the domtrees.Element i.e if the selectorOverride argument is an empty string then domtrees.Element will create an appropriate selector matching its type and uid value in this format  (ElementType[uid='UID_VALUE']) but if the selector value is not empty then that becomes the default selector used
match the event with. */
func DragLeave(fx domtrees.EventHandler, selectorOverride string) *domtrees.Event {
	return domtrees.NewEvent("dragleave", selectorOverride, fx)
}

// DragOver Documentation is as below:
// An element or text selection is being dragged over a valid drop target (every 350ms).
// https://developer.mozilla.org/docs/Web/Events/dragover
/* This event provides options() to be called when the events is triggered and an optional selector which will override the internal selector mechanism of the domtrees.Element i.e if the selectorOverride argument is an empty string then domtrees.Element will create an appropriate selector matching its type and uid value in this format  (ElementType[uid='UID_VALUE']) but if the selector value is not empty then that becomes the default selector used
match the event with. */
func DragOver(fx domtrees.EventHandler, selectorOverride string) *domtrees.Event {
	return domtrees.NewEvent("dragover", selectorOverride, fx)
}

// DragStart Documentation is as below:
// The user starts dragging an element or text selection.
// https://developer.mozilla.org/docs/Web/Events/dragstart
/* This event provides options() to be called when the events is triggered and an optional selector which will override the internal selector mechanism of the domtrees.Element i.e if the selectorOverride argument is an empty string then domtrees.Element will create an appropriate selector matching its type and uid value in this format  (ElementType[uid='UID_VALUE']) but if the selector value is not empty then that becomes the default selector used
match the event with. */
func DragStart(fx domtrees.EventHandler, selectorOverride string) *domtrees.Event {
	return domtrees.NewEvent("dragstart", selectorOverride, fx)
}

// Drop Documentation is as below:
// An element is dropped on a valid drop target.
// https://developer.mozilla.org/docs/Web/Events/drop
/* This event provides options() to be called when the events is triggered and an optional selector which will override the internal selector mechanism of the domtrees.Element i.e if the selectorOverride argument is an empty string then domtrees.Element will create an appropriate selector matching its type and uid value in this format  (ElementType[uid='UID_VALUE']) but if the selector value is not empty then that becomes the default selector used
match the event with. */
func Drop(fx domtrees.EventHandler, selectorOverride string) *domtrees.Event {
	return domtrees.NewEvent("drop", selectorOverride, fx)
}

// DurationChange Documentation is as below:
// The duration attribute has been updated.
// https://developer.mozilla.org/docs/Web/Events/durationchange
/* This event provides options() to be called when the events is triggered and an optional selector which will override the internal selector mechanism of the domtrees.Element i.e if the selectorOverride argument is an empty string then domtrees.Element will create an appropriate selector matching its type and uid value in this format  (ElementType[uid='UID_VALUE']) but if the selector value is not empty then that becomes the default selector used
match the event with. */
func DurationChange(fx domtrees.EventHandler, selectorOverride string) *domtrees.Event {
	return domtrees.NewEvent("durationchange", selectorOverride, fx)
}

// Emptied Documentation is as below:
// The media has become empty; for example, this event is sent if the media has already been loaded (or partially loaded), and the load() method is called to reload it.
// https://developer.mozilla.org/docs/Web/Events/emptied
/* This event provides options() to be called when the events is triggered and an optional selector which will override the internal selector mechanism of the domtrees.Element i.e if the selectorOverride argument is an empty string then domtrees.Element will create an appropriate selector matching its type and uid value in this format  (ElementType[uid='UID_VALUE']) but if the selector value is not empty then that becomes the default selector used
match the event with. */
func Emptied(fx domtrees.EventHandler, selectorOverride string) *domtrees.Event {
	return domtrees.NewEvent("emptied", selectorOverride, fx)
}

// End Documentation is as below:
// The utterance has finished being spoken.
// https://developer.mozilla.org/docs/Web/Events/end_(SpeechSynthesis)
/* This event provides options() to be called when the events is triggered and an optional selector which will override the internal selector mechanism of the domtrees.Element i.e if the selectorOverride argument is an empty string then domtrees.Element will create an appropriate selector matching its type and uid value in this format  (ElementType[uid='UID_VALUE']) but if the selector value is not empty then that becomes the default selector used
match the event with. */
func End(fx domtrees.EventHandler, selectorOverride string) *domtrees.Event {
	return domtrees.NewEvent("end", selectorOverride, fx)
}

// EndEvent Documentation is as below:
// A SMIL animation element ends.
// https://developer.mozilla.org/docs/Web/Events/endEvent
/* This event provides options() to be called when the events is triggered and an optional selector which will override the internal selector mechanism of the domtrees.Element i.e if the selectorOverride argument is an empty string then domtrees.Element will create an appropriate selector matching its type and uid value in this format  (ElementType[uid='UID_VALUE']) but if the selector value is not empty then that becomes the default selector used
match the event with. */
func EndEvent(fx domtrees.EventHandler, selectorOverride string) *domtrees.Event {
	return domtrees.NewEvent("endEvent", selectorOverride, fx)
}

// Ended Documentation is as below:
// (no documentation)
// https://developer.mozilla.org/docs/Web/Events/ended_(Web_Audio)
/* This event provides options() to be called when the events is triggered and an optional selector which will override the internal selector mechanism of the domtrees.Element i.e if the selectorOverride argument is an empty string then domtrees.Element will create an appropriate selector matching its type and uid value in this format  (ElementType[uid='UID_VALUE']) but if the selector value is not empty then that becomes the default selector used
match the event with. */
func Ended(fx domtrees.EventHandler, selectorOverride string) *domtrees.Event {
	return domtrees.NewEvent("ended", selectorOverride, fx)
}

// Error Documentation is as below:
// An error occurs that prevents the utterance from being succesfully spoken.
// https://developer.mozilla.org/docs/Web/Events/error_(SpeechSynthesisError)
/* This event provides options() to be called when the events is triggered and an optional selector which will override the internal selector mechanism of the domtrees.Element i.e if the selectorOverride argument is an empty string then domtrees.Element will create an appropriate selector matching its type and uid value in this format  (ElementType[uid='UID_VALUE']) but if the selector value is not empty then that becomes the default selector used
match the event with. */
func Error(fx domtrees.EventHandler, selectorOverride string) *domtrees.Event {
	return domtrees.NewEvent("error", selectorOverride, fx)
}

// Focus Documentation is as below:
// An element has received focus (does not bubble).
// https://developer.mozilla.org/docs/Web/Events/focus
/* This event provides options() to be called when the events is triggered and an optional selector which will override the internal selector mechanism of the domtrees.Element i.e if the selectorOverride argument is an empty string then domtrees.Element will create an appropriate selector matching its type and uid value in this format  (ElementType[uid='UID_VALUE']) but if the selector value is not empty then that becomes the default selector used
match the event with. */
func Focus(fx domtrees.EventHandler, selectorOverride string) *domtrees.Event {
	return domtrees.NewEvent("focus", selectorOverride, fx)
}

// FocusIn Documentation is as below:
// An element is about to receive focus (bubbles).
// https://developer.mozilla.org/docs/Web/Events/focusin
/* This event provides options() to be called when the events is triggered and an optional selector which will override the internal selector mechanism of the domtrees.Element i.e if the selectorOverride argument is an empty string then domtrees.Element will create an appropriate selector matching its type and uid value in this format  (ElementType[uid='UID_VALUE']) but if the selector value is not empty then that becomes the default selector used
match the event with. */
func FocusIn(fx domtrees.EventHandler, selectorOverride string) *domtrees.Event {
	return domtrees.NewEvent("focusin", selectorOverride, fx)
}

// FocusOut Documentation is as below:
// An element is about to lose focus (bubbles).
// https://developer.mozilla.org/docs/Web/Events/focusout
/* This event provides options() to be called when the events is triggered and an optional selector which will override the internal selector mechanism of the domtrees.Element i.e if the selectorOverride argument is an empty string then domtrees.Element will create an appropriate selector matching its type and uid value in this format  (ElementType[uid='UID_VALUE']) but if the selector value is not empty then that becomes the default selector used
match the event with. */
func FocusOut(fx domtrees.EventHandler, selectorOverride string) *domtrees.Event {
	return domtrees.NewEvent("focusout", selectorOverride, fx)
}

// FullScreenChange Documentation is as below:
// An element was turned to fullscreen mode or back to normal mode.
// https://developer.mozilla.org/docs/Web/Events/fullscreenchange
/* This event provides options() to be called when the events is triggered and an optional selector which will override the internal selector mechanism of the domtrees.Element i.e if the selectorOverride argument is an empty string then domtrees.Element will create an appropriate selector matching its type and uid value in this format  (ElementType[uid='UID_VALUE']) but if the selector value is not empty then that becomes the default selector used
match the event with. */
func FullScreenChange(fx domtrees.EventHandler, selectorOverride string) *domtrees.Event {
	return domtrees.NewEvent("fullscreenchange", selectorOverride, fx)
}

// FullScreenError Documentation is as below:
// It was impossible to switch to fullscreen mode for technical reasons or because the permission was denied.
// https://developer.mozilla.org/docs/Web/Events/fullscreenerror
/* This event provides options() to be called when the events is triggered and an optional selector which will override the internal selector mechanism of the domtrees.Element i.e if the selectorOverride argument is an empty string then domtrees.Element will create an appropriate selector matching its type and uid value in this format  (ElementType[uid='UID_VALUE']) but if the selector value is not empty then that becomes the default selector used
match the event with. */
func FullScreenError(fx domtrees.EventHandler, selectorOverride string) *domtrees.Event {
	return domtrees.NewEvent("fullscreenerror", selectorOverride, fx)
}

// GamepadConnected Documentation is as below:
// A gamepad has been connected.
// https://developer.mozilla.org/docs/Web/Events/gamepadconnected
/* This event provides options() to be called when the events is triggered and an optional selector which will override the internal selector mechanism of the domtrees.Element i.e if the selectorOverride argument is an empty string then domtrees.Element will create an appropriate selector matching its type and uid value in this format  (ElementType[uid='UID_VALUE']) but if the selector value is not empty then that becomes the default selector used
match the event with. */
func GamepadConnected(fx domtrees.EventHandler, selectorOverride string) *domtrees.Event {
	return domtrees.NewEvent("gamepadconnected", selectorOverride, fx)
}

// GamepadDisconnected Documentation is as below:
// A gamepad has been disconnected.
// https://developer.mozilla.org/docs/Web/Events/gamepaddisconnected
/* This event provides options() to be called when the events is triggered and an optional selector which will override the internal selector mechanism of the domtrees.Element i.e if the selectorOverride argument is an empty string then domtrees.Element will create an appropriate selector matching its type and uid value in this format  (ElementType[uid='UID_VALUE']) but if the selector value is not empty then that becomes the default selector used
match the event with. */
func GamepadDisconnected(fx domtrees.EventHandler, selectorOverride string) *domtrees.Event {
	return domtrees.NewEvent("gamepaddisconnected", selectorOverride, fx)
}

// Gotpointercapture Documentation is as below:
// Element receives pointer capture.
// https://developer.mozilla.org/docs/Web/Events/gotpointercapture
/* This event provides options() to be called when the events is triggered and an optional selector which will override the internal selector mechanism of the domtrees.Element i.e if the selectorOverride argument is an empty string then domtrees.Element will create an appropriate selector matching its type and uid value in this format  (ElementType[uid='UID_VALUE']) but if the selector value is not empty then that becomes the default selector used
match the event with. */
func Gotpointercapture(fx domtrees.EventHandler, selectorOverride string) *domtrees.Event {
	return domtrees.NewEvent("gotpointercapture", selectorOverride, fx)
}

// HashChange Documentation is as below:
// The fragment identifier of the URL has changed (the part of the URL after the #).
// https://developer.mozilla.org/docs/Web/Events/hashchange
/* This event provides options() to be called when the events is triggered and an optional selector which will override the internal selector mechanism of the domtrees.Element i.e if the selectorOverride argument is an empty string then domtrees.Element will create an appropriate selector matching its type and uid value in this format  (ElementType[uid='UID_VALUE']) but if the selector value is not empty then that becomes the default selector used
match the event with. */
func HashChange(fx domtrees.EventHandler, selectorOverride string) *domtrees.Event {
	return domtrees.NewEvent("hashchange", selectorOverride, fx)
}

// Input Documentation is as below:
// The value of an element changes or the content of an element with the attribute contenteditable is modified.
// https://developer.mozilla.org/docs/Web/Events/input
/* This event provides options() to be called when the events is triggered and an optional selector which will override the internal selector mechanism of the domtrees.Element i.e if the selectorOverride argument is an empty string then domtrees.Element will create an appropriate selector matching its type and uid value in this format  (ElementType[uid='UID_VALUE']) but if the selector value is not empty then that becomes the default selector used
match the event with. */
func Input(fx domtrees.EventHandler, selectorOverride string) *domtrees.Event {
	return domtrees.NewEvent("input", selectorOverride, fx)
}

// Invalid Documentation is as below:
// A submittable element has been checked and doesn't satisfy its constraints.
// https://developer.mozilla.org/docs/Web/Events/invalid
/* This event provides options() to be called when the events is triggered and an optional selector which will override the internal selector mechanism of the domtrees.Element i.e if the selectorOverride argument is an empty string then domtrees.Element will create an appropriate selector matching its type and uid value in this format  (ElementType[uid='UID_VALUE']) but if the selector value is not empty then that becomes the default selector used
match the event with. */
func Invalid(fx domtrees.EventHandler, selectorOverride string) *domtrees.Event {
	return domtrees.NewEvent("invalid", selectorOverride, fx)
}

// KeyDown Documentation is as below:
// A key is pressed down.
// https://developer.mozilla.org/docs/Web/Events/keydown
/* This event provides options() to be called when the events is triggered and an optional selector which will override the internal selector mechanism of the domtrees.Element i.e if the selectorOverride argument is an empty string then domtrees.Element will create an appropriate selector matching its type and uid value in this format  (ElementType[uid='UID_VALUE']) but if the selector value is not empty then that becomes the default selector used
match the event with. */
func KeyDown(fx domtrees.EventHandler, selectorOverride string) *domtrees.Event {
	return domtrees.NewEvent("keydown", selectorOverride, fx)
}

// KeyPress Documentation is as below:
// A key is pressed down and that key normally produces a character value (use input instead).
// https://developer.mozilla.org/docs/Web/Events/keypress
/* This event provides options() to be called when the events is triggered and an optional selector which will override the internal selector mechanism of the domtrees.Element i.e if the selectorOverride argument is an empty string then domtrees.Element will create an appropriate selector matching its type and uid value in this format  (ElementType[uid='UID_VALUE']) but if the selector value is not empty then that becomes the default selector used
match the event with. */
func KeyPress(fx domtrees.EventHandler, selectorOverride string) *domtrees.Event {
	return domtrees.NewEvent("keypress", selectorOverride, fx)
}

// KeyUp Documentation is as below:
// A key is released.
// https://developer.mozilla.org/docs/Web/Events/keyup
/* This event provides options() to be called when the events is triggered and an optional selector which will override the internal selector mechanism of the domtrees.Element i.e if the selectorOverride argument is an empty string then domtrees.Element will create an appropriate selector matching its type and uid value in this format  (ElementType[uid='UID_VALUE']) but if the selector value is not empty then that becomes the default selector used
match the event with. */
func KeyUp(fx domtrees.EventHandler, selectorOverride string) *domtrees.Event {
	return domtrees.NewEvent("keyup", selectorOverride, fx)
}

// LanguageChange Documentation is as below:
// (no documentation)
// https://developer.mozilla.org/docs/Web/Events/languagechange
/* This event provides options() to be called when the events is triggered and an optional selector which will override the internal selector mechanism of the domtrees.Element i.e if the selectorOverride argument is an empty string then domtrees.Element will create an appropriate selector matching its type and uid value in this format  (ElementType[uid='UID_VALUE']) but if the selector value is not empty then that becomes the default selector used
match the event with. */
func LanguageChange(fx domtrees.EventHandler, selectorOverride string) *domtrees.Event {
	return domtrees.NewEvent("languagechange", selectorOverride, fx)
}

// LevelChange Documentation is as below:
// The level attribute has been updated.
// https://developer.mozilla.org/docs/Web/Events/levelchange
/* This event provides options() to be called when the events is triggered and an optional selector which will override the internal selector mechanism of the domtrees.Element i.e if the selectorOverride argument is an empty string then domtrees.Element will create an appropriate selector matching its type and uid value in this format  (ElementType[uid='UID_VALUE']) but if the selector value is not empty then that becomes the default selector used
match the event with. */
func LevelChange(fx domtrees.EventHandler, selectorOverride string) *domtrees.Event {
	return domtrees.NewEvent("levelchange", selectorOverride, fx)
}

// Load Documentation is as below:
// Progression has been successful.
// https://developer.mozilla.org/docs/Web/Reference/Events/load_(ProgressEvent)
/* This event provides options() to be called when the events is triggered and an optional selector which will override the internal selector mechanism of the domtrees.Element i.e if the selectorOverride argument is an empty string then domtrees.Element will create an appropriate selector matching its type and uid value in this format  (ElementType[uid='UID_VALUE']) but if the selector value is not empty then that becomes the default selector used
match the event with. */
func Load(fx domtrees.EventHandler, selectorOverride string) *domtrees.Event {
	return domtrees.NewEvent("load", selectorOverride, fx)
}

// LoadEnd Documentation is as below:
// Progress has stopped (after "error", "abort" or "load" have been dispatched).
// https://developer.mozilla.org/docs/Web/Events/loadend
/* This event provides options() to be called when the events is triggered and an optional selector which will override the internal selector mechanism of the domtrees.Element i.e if the selectorOverride argument is an empty string then domtrees.Element will create an appropriate selector matching its type and uid value in this format  (ElementType[uid='UID_VALUE']) but if the selector value is not empty then that becomes the default selector used
match the event with. */
func LoadEnd(fx domtrees.EventHandler, selectorOverride string) *domtrees.Event {
	return domtrees.NewEvent("loadend", selectorOverride, fx)
}

// LoadStart Documentation is as below:
// Progress has begun.
// https://developer.mozilla.org/docs/Web/Events/loadstart
/* This event provides options() to be called when the events is triggered and an optional selector which will override the internal selector mechanism of the domtrees.Element i.e if the selectorOverride argument is an empty string then domtrees.Element will create an appropriate selector matching its type and uid value in this format  (ElementType[uid='UID_VALUE']) but if the selector value is not empty then that becomes the default selector used
match the event with. */
func LoadStart(fx domtrees.EventHandler, selectorOverride string) *domtrees.Event {
	return domtrees.NewEvent("loadstart", selectorOverride, fx)
}

// LoadedData Documentation is as below:
// The first frame of the media has finished loading.
// https://developer.mozilla.org/docs/Web/Events/loadeddata
/* This event provides options() to be called when the events is triggered and an optional selector which will override the internal selector mechanism of the domtrees.Element i.e if the selectorOverride argument is an empty string then domtrees.Element will create an appropriate selector matching its type and uid value in this format  (ElementType[uid='UID_VALUE']) but if the selector value is not empty then that becomes the default selector used
match the event with. */
func LoadedData(fx domtrees.EventHandler, selectorOverride string) *domtrees.Event {
	return domtrees.NewEvent("loadeddata", selectorOverride, fx)
}

// LoadedMetadata Documentation is as below:
// The metadata has been loaded.
// https://developer.mozilla.org/docs/Web/Events/loadedmetadata
/* This event provides options() to be called when the events is triggered and an optional selector which will override the internal selector mechanism of the domtrees.Element i.e if the selectorOverride argument is an empty string then domtrees.Element will create an appropriate selector matching its type and uid value in this format  (ElementType[uid='UID_VALUE']) but if the selector value is not empty then that becomes the default selector used
match the event with. */
func LoadedMetadata(fx domtrees.EventHandler, selectorOverride string) *domtrees.Event {
	return domtrees.NewEvent("loadedmetadata", selectorOverride, fx)
}

// Lostpointercapture Documentation is as below:
// Element lost pointer capture.
// https://developer.mozilla.org/docs/Web/Events/lostpointercapture
/* This event provides options() to be called when the events is triggered and an optional selector which will override the internal selector mechanism of the domtrees.Element i.e if the selectorOverride argument is an empty string then domtrees.Element will create an appropriate selector matching its type and uid value in this format  (ElementType[uid='UID_VALUE']) but if the selector value is not empty then that becomes the default selector used
match the event with. */
func Lostpointercapture(fx domtrees.EventHandler, selectorOverride string) *domtrees.Event {
	return domtrees.NewEvent("lostpointercapture", selectorOverride, fx)
}

// Mark Documentation is as below:
// The spoken utterance reaches a named SSML "mark" tag.
// https://developer.mozilla.org/docs/Web/Events/mark
/* This event provides options() to be called when the events is triggered and an optional selector which will override the internal selector mechanism of the domtrees.Element i.e if the selectorOverride argument is an empty string then domtrees.Element will create an appropriate selector matching its type and uid value in this format  (ElementType[uid='UID_VALUE']) but if the selector value is not empty then that becomes the default selector used
match the event with. */
func Mark(fx domtrees.EventHandler, selectorOverride string) *domtrees.Event {
	return domtrees.NewEvent("mark", selectorOverride, fx)
}

// Message Documentation is as below:
// A message is received through an event source.
// https://developer.mozilla.org/docs/Web/Reference/Events/message_serversentevents
/* This event provides options() to be called when the events is triggered and an optional selector which will override the internal selector mechanism of the domtrees.Element i.e if the selectorOverride argument is an empty string then domtrees.Element will create an appropriate selector matching its type and uid value in this format  (ElementType[uid='UID_VALUE']) but if the selector value is not empty then that becomes the default selector used
match the event with. */
func Message(fx domtrees.EventHandler, selectorOverride string) *domtrees.Event {
	return domtrees.NewEvent("message", selectorOverride, fx)
}

// MouseDown Documentation is as below:
// A pointing device button (usually a mouse) is pressed on an element.
// https://developer.mozilla.org/docs/Web/Events/mousedown
/* This event provides options() to be called when the events is triggered and an optional selector which will override the internal selector mechanism of the domtrees.Element i.e if the selectorOverride argument is an empty string then domtrees.Element will create an appropriate selector matching its type and uid value in this format  (ElementType[uid='UID_VALUE']) but if the selector value is not empty then that becomes the default selector used
match the event with. */
func MouseDown(fx domtrees.EventHandler, selectorOverride string) *domtrees.Event {
	return domtrees.NewEvent("mousedown", selectorOverride, fx)
}

// MouseEnter Documentation is as below:
// A pointing device is moved onto the element that has the listener attached.
// https://developer.mozilla.org/docs/Web/Events/mouseenter
/* This event provides options() to be called when the events is triggered and an optional selector which will override the internal selector mechanism of the domtrees.Element i.e if the selectorOverride argument is an empty string then domtrees.Element will create an appropriate selector matching its type and uid value in this format  (ElementType[uid='UID_VALUE']) but if the selector value is not empty then that becomes the default selector used
match the event with. */
func MouseEnter(fx domtrees.EventHandler, selectorOverride string) *domtrees.Event {
	return domtrees.NewEvent("mouseenter", selectorOverride, fx)
}

// MouseLeave Documentation is as below:
// A pointing device is moved off the element that has the listener attached.
// https://developer.mozilla.org/docs/Web/Events/mouseleave
/* This event provides options() to be called when the events is triggered and an optional selector which will override the internal selector mechanism of the domtrees.Element i.e if the selectorOverride argument is an empty string then domtrees.Element will create an appropriate selector matching its type and uid value in this format  (ElementType[uid='UID_VALUE']) but if the selector value is not empty then that becomes the default selector used
match the event with. */
func MouseLeave(fx domtrees.EventHandler, selectorOverride string) *domtrees.Event {
	return domtrees.NewEvent("mouseleave", selectorOverride, fx)
}

// MouseMove Documentation is as below:
// A pointing device is moved over an element.
// https://developer.mozilla.org/docs/Web/Events/mousemove
/* This event provides options() to be called when the events is triggered and an optional selector which will override the internal selector mechanism of the domtrees.Element i.e if the selectorOverride argument is an empty string then domtrees.Element will create an appropriate selector matching its type and uid value in this format  (ElementType[uid='UID_VALUE']) but if the selector value is not empty then that becomes the default selector used
match the event with. */
func MouseMove(fx domtrees.EventHandler, selectorOverride string) *domtrees.Event {
	return domtrees.NewEvent("mousemove", selectorOverride, fx)
}

// MouseOut Documentation is as below:
// A pointing device is moved off the element that has the listener attached or off one of its children.
// https://developer.mozilla.org/docs/Web/Events/mouseout
/* This event provides options() to be called when the events is triggered and an optional selector which will override the internal selector mechanism of the domtrees.Element i.e if the selectorOverride argument is an empty string then domtrees.Element will create an appropriate selector matching its type and uid value in this format  (ElementType[uid='UID_VALUE']) but if the selector value is not empty then that becomes the default selector used
match the event with. */
func MouseOut(fx domtrees.EventHandler, selectorOverride string) *domtrees.Event {
	return domtrees.NewEvent("mouseout", selectorOverride, fx)
}

// MouseOver Documentation is as below:
// A pointing device is moved onto the element that has the listener attached or onto one of its children.
// https://developer.mozilla.org/docs/Web/Events/mouseover
/* This event provides options() to be called when the events is triggered and an optional selector which will override the internal selector mechanism of the domtrees.Element i.e if the selectorOverride argument is an empty string then domtrees.Element will create an appropriate selector matching its type and uid value in this format  (ElementType[uid='UID_VALUE']) but if the selector value is not empty then that becomes the default selector used
match the event with. */
func MouseOver(fx domtrees.EventHandler, selectorOverride string) *domtrees.Event {
	return domtrees.NewEvent("mouseover", selectorOverride, fx)
}

// MouseUp Documentation is as below:
// A pointing device button is released over an element.
// https://developer.mozilla.org/docs/Web/Events/mouseup
/* This event provides options() to be called when the events is triggered and an optional selector which will override the internal selector mechanism of the domtrees.Element i.e if the selectorOverride argument is an empty string then domtrees.Element will create an appropriate selector matching its type and uid value in this format  (ElementType[uid='UID_VALUE']) but if the selector value is not empty then that becomes the default selector used
match the event with. */
func MouseUp(fx domtrees.EventHandler, selectorOverride string) *domtrees.Event {
	return domtrees.NewEvent("mouseup", selectorOverride, fx)
}

// NoUpdate Documentation is as below:
// The manifest hadn't changed.
// https://developer.mozilla.org/docs/Web/Events/noupdate
/* This event provides options() to be called when the events is triggered and an optional selector which will override the internal selector mechanism of the domtrees.Element i.e if the selectorOverride argument is an empty string then domtrees.Element will create an appropriate selector matching its type and uid value in this format  (ElementType[uid='UID_VALUE']) but if the selector value is not empty then that becomes the default selector used
match the event with. */
func NoUpdate(fx domtrees.EventHandler, selectorOverride string) *domtrees.Event {
	return domtrees.NewEvent("noupdate", selectorOverride, fx)
}

// Nomatch Documentation is as below:
// The speech recognition service returns a final result with no significant recognition.
// https://developer.mozilla.org/docs/Web/Events/nomatch
/* This event provides options() to be called when the events is triggered and an optional selector which will override the internal selector mechanism of the domtrees.Element i.e if the selectorOverride argument is an empty string then domtrees.Element will create an appropriate selector matching its type and uid value in this format  (ElementType[uid='UID_VALUE']) but if the selector value is not empty then that becomes the default selector used
match the event with. */
func Nomatch(fx domtrees.EventHandler, selectorOverride string) *domtrees.Event {
	return domtrees.NewEvent("nomatch", selectorOverride, fx)
}

// Notificationclick Documentation is as below:
// A system notification spawned by ServiceWorkerRegistration.showNotification() has been clicked.
// https://developer.mozilla.org/docs/Web/Events/notificationclick
/* This event provides options() to be called when the events is triggered and an optional selector which will override the internal selector mechanism of the domtrees.Element i.e if the selectorOverride argument is an empty string then domtrees.Element will create an appropriate selector matching its type and uid value in this format  (ElementType[uid='UID_VALUE']) but if the selector value is not empty then that becomes the default selector used
match the event with. */
func Notificationclick(fx domtrees.EventHandler, selectorOverride string) *domtrees.Event {
	return domtrees.NewEvent("notificationclick", selectorOverride, fx)
}

// Obsolete Documentation is as below:
// The manifest was found to have become a 404 or 410 page, so the application cache is being deleted.
// https://developer.mozilla.org/docs/Web/Events/obsolete
/* This event provides options() to be called when the events is triggered and an optional selector which will override the internal selector mechanism of the domtrees.Element i.e if the selectorOverride argument is an empty string then domtrees.Element will create an appropriate selector matching its type and uid value in this format  (ElementType[uid='UID_VALUE']) but if the selector value is not empty then that becomes the default selector used
match the event with. */
func Obsolete(fx domtrees.EventHandler, selectorOverride string) *domtrees.Event {
	return domtrees.NewEvent("obsolete", selectorOverride, fx)
}

// Offline Documentation is as below:
// The browser has lost access to the network.
// https://developer.mozilla.org/docs/Web/Events/offline
/* This event provides options() to be called when the events is triggered and an optional selector which will override the internal selector mechanism of the domtrees.Element i.e if the selectorOverride argument is an empty string then domtrees.Element will create an appropriate selector matching its type and uid value in this format  (ElementType[uid='UID_VALUE']) but if the selector value is not empty then that becomes the default selector used
match the event with. */
func Offline(fx domtrees.EventHandler, selectorOverride string) *domtrees.Event {
	return domtrees.NewEvent("offline", selectorOverride, fx)
}

// Online Documentation is as below:
// The browser has gained access to the network (but particular websites might be unreachable).
// https://developer.mozilla.org/docs/Web/Events/online
/* This event provides options() to be called when the events is triggered and an optional selector which will override the internal selector mechanism of the domtrees.Element i.e if the selectorOverride argument is an empty string then domtrees.Element will create an appropriate selector matching its type and uid value in this format  (ElementType[uid='UID_VALUE']) but if the selector value is not empty then that becomes the default selector used
match the event with. */
func Online(fx domtrees.EventHandler, selectorOverride string) *domtrees.Event {
	return domtrees.NewEvent("online", selectorOverride, fx)
}

// Open Documentation is as below:
// An event source connection has been established.
// https://developer.mozilla.org/docs/Web/Reference/Events/open_serversentevents
/* This event provides options() to be called when the events is triggered and an optional selector which will override the internal selector mechanism of the domtrees.Element i.e if the selectorOverride argument is an empty string then domtrees.Element will create an appropriate selector matching its type and uid value in this format  (ElementType[uid='UID_VALUE']) but if the selector value is not empty then that becomes the default selector used
match the event with. */
func Open(fx domtrees.EventHandler, selectorOverride string) *domtrees.Event {
	return domtrees.NewEvent("open", selectorOverride, fx)
}

// OrientationChange Documentation is as below:
// The orientation of the device (portrait/landscape) has changed
// https://developer.mozilla.org/docs/Web/Events/orientationchange
/* This event provides options() to be called when the events is triggered and an optional selector which will override the internal selector mechanism of the domtrees.Element i.e if the selectorOverride argument is an empty string then domtrees.Element will create an appropriate selector matching its type and uid value in this format  (ElementType[uid='UID_VALUE']) but if the selector value is not empty then that becomes the default selector used
match the event with. */
func OrientationChange(fx domtrees.EventHandler, selectorOverride string) *domtrees.Event {
	return domtrees.NewEvent("orientationchange", selectorOverride, fx)
}

// PageHide Documentation is as below:
// A session history entry is being traversed from.
// https://developer.mozilla.org/docs/Web/Events/pagehide
/* This event provides options() to be called when the events is triggered and an optional selector which will override the internal selector mechanism of the domtrees.Element i.e if the selectorOverride argument is an empty string then domtrees.Element will create an appropriate selector matching its type and uid value in this format  (ElementType[uid='UID_VALUE']) but if the selector value is not empty then that becomes the default selector used
match the event with. */
func PageHide(fx domtrees.EventHandler, selectorOverride string) *domtrees.Event {
	return domtrees.NewEvent("pagehide", selectorOverride, fx)
}

// PageShow Documentation is as below:
// A session history entry is being traversed to.
// https://developer.mozilla.org/docs/Web/Events/pageshow
/* This event provides options() to be called when the events is triggered and an optional selector which will override the internal selector mechanism of the domtrees.Element i.e if the selectorOverride argument is an empty string then domtrees.Element will create an appropriate selector matching its type and uid value in this format  (ElementType[uid='UID_VALUE']) but if the selector value is not empty then that becomes the default selector used
match the event with. */
func PageShow(fx domtrees.EventHandler, selectorOverride string) *domtrees.Event {
	return domtrees.NewEvent("pageshow", selectorOverride, fx)
}

// Paste Documentation is as below:
// Data has been transfered from the system clipboard to the document.
// https://developer.mozilla.org/docs/Web/Events/paste
/* This event provides options() to be called when the events is triggered and an optional selector which will override the internal selector mechanism of the domtrees.Element i.e if the selectorOverride argument is an empty string then domtrees.Element will create an appropriate selector matching its type and uid value in this format  (ElementType[uid='UID_VALUE']) but if the selector value is not empty then that becomes the default selector used
match the event with. */
func Paste(fx domtrees.EventHandler, selectorOverride string) *domtrees.Event {
	return domtrees.NewEvent("paste", selectorOverride, fx)
}

// Pause Documentation is as below:
// The utterance is paused part way through.
// https://developer.mozilla.org/docs/Web/Events/pause_(SpeechSynthesis)
/* This event provides options() to be called when the events is triggered and an optional selector which will override the internal selector mechanism of the domtrees.Element i.e if the selectorOverride argument is an empty string then domtrees.Element will create an appropriate selector matching its type and uid value in this format  (ElementType[uid='UID_VALUE']) but if the selector value is not empty then that becomes the default selector used
match the event with. */
func Pause(fx domtrees.EventHandler, selectorOverride string) *domtrees.Event {
	return domtrees.NewEvent("pause", selectorOverride, fx)
}

// Play Documentation is as below:
// Playback has begun.
// https://developer.mozilla.org/docs/Web/Events/play
/* This event provides options() to be called when the events is triggered and an optional selector which will override the internal selector mechanism of the domtrees.Element i.e if the selectorOverride argument is an empty string then domtrees.Element will create an appropriate selector matching its type and uid value in this format  (ElementType[uid='UID_VALUE']) but if the selector value is not empty then that becomes the default selector used
match the event with. */
func Play(fx domtrees.EventHandler, selectorOverride string) *domtrees.Event {
	return domtrees.NewEvent("play", selectorOverride, fx)
}

// Playing Documentation is as below:
// Playback is ready to start after having been paused or delayed due to lack of data.
// https://developer.mozilla.org/docs/Web/Events/playing
/* This event provides options() to be called when the events is triggered and an optional selector which will override the internal selector mechanism of the domtrees.Element i.e if the selectorOverride argument is an empty string then domtrees.Element will create an appropriate selector matching its type and uid value in this format  (ElementType[uid='UID_VALUE']) but if the selector value is not empty then that becomes the default selector used
match the event with. */
func Playing(fx domtrees.EventHandler, selectorOverride string) *domtrees.Event {
	return domtrees.NewEvent("playing", selectorOverride, fx)
}

// PointerLockChange Documentation is as below:
// The pointer was locked or released.
// https://developer.mozilla.org/docs/Web/Events/pointerlockchange
/* This event provides options() to be called when the events is triggered and an optional selector which will override the internal selector mechanism of the domtrees.Element i.e if the selectorOverride argument is an empty string then domtrees.Element will create an appropriate selector matching its type and uid value in this format  (ElementType[uid='UID_VALUE']) but if the selector value is not empty then that becomes the default selector used
match the event with. */
func PointerLockChange(fx domtrees.EventHandler, selectorOverride string) *domtrees.Event {
	return domtrees.NewEvent("pointerlockchange", selectorOverride, fx)
}

// PointerLockError Documentation is as below:
// It was impossible to lock the pointer for technical reasons or because the permission was denied.
// https://developer.mozilla.org/docs/Web/Events/pointerlockerror
/* This event provides options() to be called when the events is triggered and an optional selector which will override the internal selector mechanism of the domtrees.Element i.e if the selectorOverride argument is an empty string then domtrees.Element will create an appropriate selector matching its type and uid value in this format  (ElementType[uid='UID_VALUE']) but if the selector value is not empty then that becomes the default selector used
match the event with. */
func PointerLockError(fx domtrees.EventHandler, selectorOverride string) *domtrees.Event {
	return domtrees.NewEvent("pointerlockerror", selectorOverride, fx)
}

// Pointercancel Documentation is as below:
// The pointer is unlikely to produce any more events.
// https://developer.mozilla.org/docs/Web/Events/pointercancel
/* This event provides options() to be called when the events is triggered and an optional selector which will override the internal selector mechanism of the domtrees.Element i.e if the selectorOverride argument is an empty string then domtrees.Element will create an appropriate selector matching its type and uid value in this format  (ElementType[uid='UID_VALUE']) but if the selector value is not empty then that becomes the default selector used
match the event with. */
func Pointercancel(fx domtrees.EventHandler, selectorOverride string) *domtrees.Event {
	return domtrees.NewEvent("pointercancel", selectorOverride, fx)
}

// Pointerdown Documentation is as below:
// The pointer enters the active buttons state.
// https://developer.mozilla.org/docs/Web/Events/pointerdown
/* This event provides options() to be called when the events is triggered and an optional selector which will override the internal selector mechanism of the domtrees.Element i.e if the selectorOverride argument is an empty string then domtrees.Element will create an appropriate selector matching its type and uid value in this format  (ElementType[uid='UID_VALUE']) but if the selector value is not empty then that becomes the default selector used
match the event with. */
func Pointerdown(fx domtrees.EventHandler, selectorOverride string) *domtrees.Event {
	return domtrees.NewEvent("pointerdown", selectorOverride, fx)
}

// Pointerenter Documentation is as below:
// Pointing device is moved inside the hit-testing boundary.
// https://developer.mozilla.org/docs/Web/Events/pointerenter
/* This event provides options() to be called when the events is triggered and an optional selector which will override the internal selector mechanism of the domtrees.Element i.e if the selectorOverride argument is an empty string then domtrees.Element will create an appropriate selector matching its type and uid value in this format  (ElementType[uid='UID_VALUE']) but if the selector value is not empty then that becomes the default selector used
match the event with. */
func Pointerenter(fx domtrees.EventHandler, selectorOverride string) *domtrees.Event {
	return domtrees.NewEvent("pointerenter", selectorOverride, fx)
}

// Pointerleave Documentation is as below:
// Pointing device is moved out of the hit-testing boundary.
// https://developer.mozilla.org/docs/Web/Events/pointerleave
/* This event provides options() to be called when the events is triggered and an optional selector which will override the internal selector mechanism of the domtrees.Element i.e if the selectorOverride argument is an empty string then domtrees.Element will create an appropriate selector matching its type and uid value in this format  (ElementType[uid='UID_VALUE']) but if the selector value is not empty then that becomes the default selector used
match the event with. */
func Pointerleave(fx domtrees.EventHandler, selectorOverride string) *domtrees.Event {
	return domtrees.NewEvent("pointerleave", selectorOverride, fx)
}

// Pointermove Documentation is as below:
// The pointer changed coordinates.
// https://developer.mozilla.org/docs/Web/Events/pointermove
/* This event provides options() to be called when the events is triggered and an optional selector which will override the internal selector mechanism of the domtrees.Element i.e if the selectorOverride argument is an empty string then domtrees.Element will create an appropriate selector matching its type and uid value in this format  (ElementType[uid='UID_VALUE']) but if the selector value is not empty then that becomes the default selector used
match the event with. */
func Pointermove(fx domtrees.EventHandler, selectorOverride string) *domtrees.Event {
	return domtrees.NewEvent("pointermove", selectorOverride, fx)
}

// Pointerout Documentation is as below:
// The pointing device moved out of hit-testing boundary or leaves detectable hover range.
// https://developer.mozilla.org/docs/Web/Events/pointerout
/* This event provides options() to be called when the events is triggered and an optional selector which will override the internal selector mechanism of the domtrees.Element i.e if the selectorOverride argument is an empty string then domtrees.Element will create an appropriate selector matching its type and uid value in this format  (ElementType[uid='UID_VALUE']) but if the selector value is not empty then that becomes the default selector used
match the event with. */
func Pointerout(fx domtrees.EventHandler, selectorOverride string) *domtrees.Event {
	return domtrees.NewEvent("pointerout", selectorOverride, fx)
}

// Pointerover Documentation is as below:
// The pointing device is moved into the hit-testing boundary.
// https://developer.mozilla.org/docs/Web/Events/pointerover
/* This event provides options() to be called when the events is triggered and an optional selector which will override the internal selector mechanism of the domtrees.Element i.e if the selectorOverride argument is an empty string then domtrees.Element will create an appropriate selector matching its type and uid value in this format  (ElementType[uid='UID_VALUE']) but if the selector value is not empty then that becomes the default selector used
match the event with. */
func Pointerover(fx domtrees.EventHandler, selectorOverride string) *domtrees.Event {
	return domtrees.NewEvent("pointerover", selectorOverride, fx)
}

// Pointerup Documentation is as below:
// The pointer leaves the active buttons state.
// https://developer.mozilla.org/docs/Web/Events/pointerup
/* This event provides options() to be called when the events is triggered and an optional selector which will override the internal selector mechanism of the domtrees.Element i.e if the selectorOverride argument is an empty string then domtrees.Element will create an appropriate selector matching its type and uid value in this format  (ElementType[uid='UID_VALUE']) but if the selector value is not empty then that becomes the default selector used
match the event with. */
func Pointerup(fx domtrees.EventHandler, selectorOverride string) *domtrees.Event {
	return domtrees.NewEvent("pointerup", selectorOverride, fx)
}

// PopState Documentation is as below:
// A session history entry is being navigated to (in certain cases).
// https://developer.mozilla.org/docs/Web/Events/popstate
/* This event provides options() to be called when the events is triggered and an optional selector which will override the internal selector mechanism of the domtrees.Element i.e if the selectorOverride argument is an empty string then domtrees.Element will create an appropriate selector matching its type and uid value in this format  (ElementType[uid='UID_VALUE']) but if the selector value is not empty then that becomes the default selector used
match the event with. */
func PopState(fx domtrees.EventHandler, selectorOverride string) *domtrees.Event {
	return domtrees.NewEvent("popstate", selectorOverride, fx)
}

// Progress Documentation is as below:
// The user agent is downloading resources listed by the manifest.
// https://developer.mozilla.org/docs/Web/Reference/Events/progress_(appcache_event)
/* This event provides options() to be called when the events is triggered and an optional selector which will override the internal selector mechanism of the domtrees.Element i.e if the selectorOverride argument is an empty string then domtrees.Element will create an appropriate selector matching its type and uid value in this format  (ElementType[uid='UID_VALUE']) but if the selector value is not empty then that becomes the default selector used
match the event with. */
func Progress(fx domtrees.EventHandler, selectorOverride string) *domtrees.Event {
	return domtrees.NewEvent("progress", selectorOverride, fx)
}

// Push Documentation is as below:
// A Service Worker has received a push message.
// https://developer.mozilla.org/docs/Web/Events/push
/* This event provides options() to be called when the events is triggered and an optional selector which will override the internal selector mechanism of the domtrees.Element i.e if the selectorOverride argument is an empty string then domtrees.Element will create an appropriate selector matching its type and uid value in this format  (ElementType[uid='UID_VALUE']) but if the selector value is not empty then that becomes the default selector used
match the event with. */
func Push(fx domtrees.EventHandler, selectorOverride string) *domtrees.Event {
	return domtrees.NewEvent("push", selectorOverride, fx)
}

// Pushsubscriptionchange Documentation is as below:
// A PushSubscription has expired.
// https://developer.mozilla.org/docs/Web/Events/pushsubscriptionchange
/* This event provides options() to be called when the events is triggered and an optional selector which will override the internal selector mechanism of the domtrees.Element i.e if the selectorOverride argument is an empty string then domtrees.Element will create an appropriate selector matching its type and uid value in this format  (ElementType[uid='UID_VALUE']) but if the selector value is not empty then that becomes the default selector used
match the event with. */
func Pushsubscriptionchange(fx domtrees.EventHandler, selectorOverride string) *domtrees.Event {
	return domtrees.NewEvent("pushsubscriptionchange", selectorOverride, fx)
}

// RateChange Documentation is as below:
// The playback rate has changed.
// https://developer.mozilla.org/docs/Web/Events/ratechange
/* This event provides options() to be called when the events is triggered and an optional selector which will override the internal selector mechanism of the domtrees.Element i.e if the selectorOverride argument is an empty string then domtrees.Element will create an appropriate selector matching its type and uid value in this format  (ElementType[uid='UID_VALUE']) but if the selector value is not empty then that becomes the default selector used
match the event with. */
func RateChange(fx domtrees.EventHandler, selectorOverride string) *domtrees.Event {
	return domtrees.NewEvent("ratechange", selectorOverride, fx)
}

// ReadystateChange Documentation is as below:
// The readyState attribute of a document has changed.
// https://developer.mozilla.org/docs/Web/Events/readystatechange
/* This event provides options() to be called when the events is triggered and an optional selector which will override the internal selector mechanism of the domtrees.Element i.e if the selectorOverride argument is an empty string then domtrees.Element will create an appropriate selector matching its type and uid value in this format  (ElementType[uid='UID_VALUE']) but if the selector value is not empty then that becomes the default selector used
match the event with. */
func ReadystateChange(fx domtrees.EventHandler, selectorOverride string) *domtrees.Event {
	return domtrees.NewEvent("readystatechange", selectorOverride, fx)
}

// RepeatEvent Documentation is as below:
// A SMIL animation element is repeated.
// https://developer.mozilla.org/docs/Web/Events/repeatEvent
/* This event provides options() to be called when the events is triggered and an optional selector which will override the internal selector mechanism of the domtrees.Element i.e if the selectorOverride argument is an empty string then domtrees.Element will create an appropriate selector matching its type and uid value in this format  (ElementType[uid='UID_VALUE']) but if the selector value is not empty then that becomes the default selector used
match the event with. */
func RepeatEvent(fx domtrees.EventHandler, selectorOverride string) *domtrees.Event {
	return domtrees.NewEvent("repeatEvent", selectorOverride, fx)
}

// Reset Documentation is as below:
// A form is reset.
// https://developer.mozilla.org/docs/Web/Events/reset
/* This event provides options() to be called when the events is triggered and an optional selector which will override the internal selector mechanism of the domtrees.Element i.e if the selectorOverride argument is an empty string then domtrees.Element will create an appropriate selector matching its type and uid value in this format  (ElementType[uid='UID_VALUE']) but if the selector value is not empty then that becomes the default selector used
match the event with. */
func Reset(fx domtrees.EventHandler, selectorOverride string) *domtrees.Event {
	return domtrees.NewEvent("reset", selectorOverride, fx)
}

// Resize Documentation is as below:
// The document view has been resized.
// https://developer.mozilla.org/docs/Web/Events/resize
/* This event provides options() to be called when the events is triggered and an optional selector which will override the internal selector mechanism of the domtrees.Element i.e if the selectorOverride argument is an empty string then domtrees.Element will create an appropriate selector matching its type and uid value in this format  (ElementType[uid='UID_VALUE']) but if the selector value is not empty then that becomes the default selector used
match the event with. */
func Resize(fx domtrees.EventHandler, selectorOverride string) *domtrees.Event {
	return domtrees.NewEvent("resize", selectorOverride, fx)
}

// Result Documentation is as below:
// The speech recognition service returns a result — a word or phrase has been positively recognized and this has been communicated back to the app.
// https://developer.mozilla.org/docs/Web/Events/result
/* This event provides options() to be called when the events is triggered and an optional selector which will override the internal selector mechanism of the domtrees.Element i.e if the selectorOverride argument is an empty string then domtrees.Element will create an appropriate selector matching its type and uid value in this format  (ElementType[uid='UID_VALUE']) but if the selector value is not empty then that becomes the default selector used
match the event with. */
func Result(fx domtrees.EventHandler, selectorOverride string) *domtrees.Event {
	return domtrees.NewEvent("result", selectorOverride, fx)
}

// Resume Documentation is as below:
// A paused utterance is resumed.
// https://developer.mozilla.org/docs/Web/Events/resume
/* This event provides options() to be called when the events is triggered and an optional selector which will override the internal selector mechanism of the domtrees.Element i.e if the selectorOverride argument is an empty string then domtrees.Element will create an appropriate selector matching its type and uid value in this format  (ElementType[uid='UID_VALUE']) but if the selector value is not empty then that becomes the default selector used
match the event with. */
func Resume(fx domtrees.EventHandler, selectorOverride string) *domtrees.Event {
	return domtrees.NewEvent("resume", selectorOverride, fx)
}

// SVGAbort Documentation is as below:
// Page loading has been stopped before the SVG was loaded.
// https://developer.mozilla.org/docs/Web/Events/SVGAbort
/* This event provides options() to be called when the events is triggered and an optional selector which will override the internal selector mechanism of the domtrees.Element i.e if the selectorOverride argument is an empty string then domtrees.Element will create an appropriate selector matching its type and uid value in this format  (ElementType[uid='UID_VALUE']) but if the selector value is not empty then that becomes the default selector used
match the event with. */
func SVGAbort(fx domtrees.EventHandler, selectorOverride string) *domtrees.Event {
	return domtrees.NewEvent("SVGAbort", selectorOverride, fx)
}

// SVGError Documentation is as below:
// An error has occurred before the SVG was loaded.
// https://developer.mozilla.org/docs/Web/Events/SVGError
/* This event provides options() to be called when the events is triggered and an optional selector which will override the internal selector mechanism of the domtrees.Element i.e if the selectorOverride argument is an empty string then domtrees.Element will create an appropriate selector matching its type and uid value in this format  (ElementType[uid='UID_VALUE']) but if the selector value is not empty then that becomes the default selector used
match the event with. */
func SVGError(fx domtrees.EventHandler, selectorOverride string) *domtrees.Event {
	return domtrees.NewEvent("SVGError", selectorOverride, fx)
}

// SVGLoad Documentation is as below:
// An SVG document has been loaded and parsed.
// https://developer.mozilla.org/docs/Web/Events/SVGLoad
/* This event provides options() to be called when the events is triggered and an optional selector which will override the internal selector mechanism of the domtrees.Element i.e if the selectorOverride argument is an empty string then domtrees.Element will create an appropriate selector matching its type and uid value in this format  (ElementType[uid='UID_VALUE']) but if the selector value is not empty then that becomes the default selector used
match the event with. */
func SVGLoad(fx domtrees.EventHandler, selectorOverride string) *domtrees.Event {
	return domtrees.NewEvent("SVGLoad", selectorOverride, fx)
}

// SVGResize Documentation is as below:
// An SVG document is being resized.
// https://developer.mozilla.org/docs/Web/Events/SVGResize
/* This event provides options() to be called when the events is triggered and an optional selector which will override the internal selector mechanism of the domtrees.Element i.e if the selectorOverride argument is an empty string then domtrees.Element will create an appropriate selector matching its type and uid value in this format  (ElementType[uid='UID_VALUE']) but if the selector value is not empty then that becomes the default selector used
match the event with. */
func SVGResize(fx domtrees.EventHandler, selectorOverride string) *domtrees.Event {
	return domtrees.NewEvent("SVGResize", selectorOverride, fx)
}

// SVGScroll Documentation is as below:
// An SVG document is being scrolled.
// https://developer.mozilla.org/docs/Web/Events/SVGScroll
/* This event provides options() to be called when the events is triggered and an optional selector which will override the internal selector mechanism of the domtrees.Element i.e if the selectorOverride argument is an empty string then domtrees.Element will create an appropriate selector matching its type and uid value in this format  (ElementType[uid='UID_VALUE']) but if the selector value is not empty then that becomes the default selector used
match the event with. */
func SVGScroll(fx domtrees.EventHandler, selectorOverride string) *domtrees.Event {
	return domtrees.NewEvent("SVGScroll", selectorOverride, fx)
}

// SVGUnload Documentation is as below:
// An SVG document has been removed from a window or frame.
// https://developer.mozilla.org/docs/Web/Events/SVGUnload
/* This event provides options() to be called when the events is triggered and an optional selector which will override the internal selector mechanism of the domtrees.Element i.e if the selectorOverride argument is an empty string then domtrees.Element will create an appropriate selector matching its type and uid value in this format  (ElementType[uid='UID_VALUE']) but if the selector value is not empty then that becomes the default selector used
match the event with. */
func SVGUnload(fx domtrees.EventHandler, selectorOverride string) *domtrees.Event {
	return domtrees.NewEvent("SVGUnload", selectorOverride, fx)
}

// SVGZoom Documentation is as below:
// An SVG document is being zoomed.
// https://developer.mozilla.org/docs/Web/Events/SVGZoom
/* This event provides options() to be called when the events is triggered and an optional selector which will override the internal selector mechanism of the domtrees.Element i.e if the selectorOverride argument is an empty string then domtrees.Element will create an appropriate selector matching its type and uid value in this format  (ElementType[uid='UID_VALUE']) but if the selector value is not empty then that becomes the default selector used
match the event with. */
func SVGZoom(fx domtrees.EventHandler, selectorOverride string) *domtrees.Event {
	return domtrees.NewEvent("SVGZoom", selectorOverride, fx)
}

// Scroll Documentation is as below:
// The document view or an element has been scrolled.
// https://developer.mozilla.org/docs/Web/Events/scroll
/* This event provides options() to be called when the events is triggered and an optional selector which will override the internal selector mechanism of the domtrees.Element i.e if the selectorOverride argument is an empty string then domtrees.Element will create an appropriate selector matching its type and uid value in this format  (ElementType[uid='UID_VALUE']) but if the selector value is not empty then that becomes the default selector used
match the event with. */
func Scroll(fx domtrees.EventHandler, selectorOverride string) *domtrees.Event {
	return domtrees.NewEvent("scroll", selectorOverride, fx)
}

// Seeked Documentation is as below:
// A seek operation completed.
// https://developer.mozilla.org/docs/Web/Events/seeked
/* This event provides options() to be called when the events is triggered and an optional selector which will override the internal selector mechanism of the domtrees.Element i.e if the selectorOverride argument is an empty string then domtrees.Element will create an appropriate selector matching its type and uid value in this format  (ElementType[uid='UID_VALUE']) but if the selector value is not empty then that becomes the default selector used
match the event with. */
func Seeked(fx domtrees.EventHandler, selectorOverride string) *domtrees.Event {
	return domtrees.NewEvent("seeked", selectorOverride, fx)
}

// Seeking Documentation is as below:
// A seek operation began.
// https://developer.mozilla.org/docs/Web/Events/seeking
/* This event provides options() to be called when the events is triggered and an optional selector which will override the internal selector mechanism of the domtrees.Element i.e if the selectorOverride argument is an empty string then domtrees.Element will create an appropriate selector matching its type and uid value in this format  (ElementType[uid='UID_VALUE']) but if the selector value is not empty then that becomes the default selector used
match the event with. */
func Seeking(fx domtrees.EventHandler, selectorOverride string) *domtrees.Event {
	return domtrees.NewEvent("seeking", selectorOverride, fx)
}

// Select Documentation is as below:
// Some text is being selected.
// https://developer.mozilla.org/docs/Web/Events/select
/* This event provides options() to be called when the events is triggered and an optional selector which will override the internal selector mechanism of the domtrees.Element i.e if the selectorOverride argument is an empty string then domtrees.Element will create an appropriate selector matching its type and uid value in this format  (ElementType[uid='UID_VALUE']) but if the selector value is not empty then that becomes the default selector used
match the event with. */
func Select(fx domtrees.EventHandler, selectorOverride string) *domtrees.Event {
	return domtrees.NewEvent("select", selectorOverride, fx)
}

// Selectionchange Documentation is as below:
// The selection in the document has been changed.
// https://developer.mozilla.org/docs/Web/Events/selectionchange
/* This event provides options() to be called when the events is triggered and an optional selector which will override the internal selector mechanism of the domtrees.Element i.e if the selectorOverride argument is an empty string then domtrees.Element will create an appropriate selector matching its type and uid value in this format  (ElementType[uid='UID_VALUE']) but if the selector value is not empty then that becomes the default selector used
match the event with. */
func Selectionchange(fx domtrees.EventHandler, selectorOverride string) *domtrees.Event {
	return domtrees.NewEvent("selectionchange", selectorOverride, fx)
}

// Selectstart Documentation is as below:
// A selection just started.
// https://developer.mozilla.org/docs/Web/Events/selectstart
/* This event provides options() to be called when the events is triggered and an optional selector which will override the internal selector mechanism of the domtrees.Element i.e if the selectorOverride argument is an empty string then domtrees.Element will create an appropriate selector matching its type and uid value in this format  (ElementType[uid='UID_VALUE']) but if the selector value is not empty then that becomes the default selector used
match the event with. */
func Selectstart(fx domtrees.EventHandler, selectorOverride string) *domtrees.Event {
	return domtrees.NewEvent("selectstart", selectorOverride, fx)
}

// Show Documentation is as below:
// A contextmenu event was fired on/bubbled to an element that has a contextmenu attribute
// https://developer.mozilla.org/docs/Web/Events/show
/* This event provides options() to be called when the events is triggered and an optional selector which will override the internal selector mechanism of the domtrees.Element i.e if the selectorOverride argument is an empty string then domtrees.Element will create an appropriate selector matching its type and uid value in this format  (ElementType[uid='UID_VALUE']) but if the selector value is not empty then that becomes the default selector used
match the event with. */
func Show(fx domtrees.EventHandler, selectorOverride string) *domtrees.Event {
	return domtrees.NewEvent("show", selectorOverride, fx)
}

// Soundend Documentation is as below:
// Any sound — recognisable speech or not — has stopped being detected.
// https://developer.mozilla.org/docs/Web/Events/soundend
/* This event provides options() to be called when the events is triggered and an optional selector which will override the internal selector mechanism of the domtrees.Element i.e if the selectorOverride argument is an empty string then domtrees.Element will create an appropriate selector matching its type and uid value in this format  (ElementType[uid='UID_VALUE']) but if the selector value is not empty then that becomes the default selector used
match the event with. */
func Soundend(fx domtrees.EventHandler, selectorOverride string) *domtrees.Event {
	return domtrees.NewEvent("soundend", selectorOverride, fx)
}

// Soundstart Documentation is as below:
// Any sound — recognisable speech or not — has been detected.
// https://developer.mozilla.org/docs/Web/Events/soundstart
/* This event provides options() to be called when the events is triggered and an optional selector which will override the internal selector mechanism of the domtrees.Element i.e if the selectorOverride argument is an empty string then domtrees.Element will create an appropriate selector matching its type and uid value in this format  (ElementType[uid='UID_VALUE']) but if the selector value is not empty then that becomes the default selector used
match the event with. */
func Soundstart(fx domtrees.EventHandler, selectorOverride string) *domtrees.Event {
	return domtrees.NewEvent("soundstart", selectorOverride, fx)
}

// Speechend Documentation is as below:
// Speech recognised by the speech recognition service has stopped being detected.
// https://developer.mozilla.org/docs/Web/Events/speechend
/* This event provides options() to be called when the events is triggered and an optional selector which will override the internal selector mechanism of the domtrees.Element i.e if the selectorOverride argument is an empty string then domtrees.Element will create an appropriate selector matching its type and uid value in this format  (ElementType[uid='UID_VALUE']) but if the selector value is not empty then that becomes the default selector used
match the event with. */
func Speechend(fx domtrees.EventHandler, selectorOverride string) *domtrees.Event {
	return domtrees.NewEvent("speechend", selectorOverride, fx)
}

// Speechstart Documentation is as below:
// Sound that is recognised by the speech recognition service as speech has been detected.
// https://developer.mozilla.org/docs/Web/Events/speechstart
/* This event provides options() to be called when the events is triggered and an optional selector which will override the internal selector mechanism of the domtrees.Element i.e if the selectorOverride argument is an empty string then domtrees.Element will create an appropriate selector matching its type and uid value in this format  (ElementType[uid='UID_VALUE']) but if the selector value is not empty then that becomes the default selector used
match the event with. */
func Speechstart(fx domtrees.EventHandler, selectorOverride string) *domtrees.Event {
	return domtrees.NewEvent("speechstart", selectorOverride, fx)
}

// Stalled Documentation is as below:
// The user agent is trying to fetch media data, but data is unexpectedly not forthcoming.
// https://developer.mozilla.org/docs/Web/Events/stalled
/* This event provides options() to be called when the events is triggered and an optional selector which will override the internal selector mechanism of the domtrees.Element i.e if the selectorOverride argument is an empty string then domtrees.Element will create an appropriate selector matching its type and uid value in this format  (ElementType[uid='UID_VALUE']) but if the selector value is not empty then that becomes the default selector used
match the event with. */
func Stalled(fx domtrees.EventHandler, selectorOverride string) *domtrees.Event {
	return domtrees.NewEvent("stalled", selectorOverride, fx)
}

// Start Documentation is as below:
// The utterance has begun to be spoken.
// https://developer.mozilla.org/docs/Web/Events/start_(SpeechSynthesis)
/* This event provides options() to be called when the events is triggered and an optional selector which will override the internal selector mechanism of the domtrees.Element i.e if the selectorOverride argument is an empty string then domtrees.Element will create an appropriate selector matching its type and uid value in this format  (ElementType[uid='UID_VALUE']) but if the selector value is not empty then that becomes the default selector used
match the event with. */
func Start(fx domtrees.EventHandler, selectorOverride string) *domtrees.Event {
	return domtrees.NewEvent("start", selectorOverride, fx)
}

// Storage Documentation is as below:
// A storage area (localStorage or sessionStorage) has changed.
// https://developer.mozilla.org/docs/Web/Events/storage
/* This event provides options() to be called when the events is triggered and an optional selector which will override the internal selector mechanism of the domtrees.Element i.e if the selectorOverride argument is an empty string then domtrees.Element will create an appropriate selector matching its type and uid value in this format  (ElementType[uid='UID_VALUE']) but if the selector value is not empty then that becomes the default selector used
match the event with. */
func Storage(fx domtrees.EventHandler, selectorOverride string) *domtrees.Event {
	return domtrees.NewEvent("storage", selectorOverride, fx)
}

// Submit Documentation is as below:
// A form is submitted.
// https://developer.mozilla.org/docs/Web/Events/submit
/* This event provides options() to be called when the events is triggered and an optional selector which will override the internal selector mechanism of the domtrees.Element i.e if the selectorOverride argument is an empty string then domtrees.Element will create an appropriate selector matching its type and uid value in this format  (ElementType[uid='UID_VALUE']) but if the selector value is not empty then that becomes the default selector used
match the event with. */
func Submit(fx domtrees.EventHandler, selectorOverride string) *domtrees.Event {
	return domtrees.NewEvent("submit", selectorOverride, fx)
}

// Success Documentation is as below:
// A request successfully completed.
// https://developer.mozilla.org/docs/Web/Reference/Events/success_indexedDB
/* This event provides options() to be called when the events is triggered and an optional selector which will override the internal selector mechanism of the domtrees.Element i.e if the selectorOverride argument is an empty string then domtrees.Element will create an appropriate selector matching its type and uid value in this format  (ElementType[uid='UID_VALUE']) but if the selector value is not empty then that becomes the default selector used
match the event with. */
func Success(fx domtrees.EventHandler, selectorOverride string) *domtrees.Event {
	return domtrees.NewEvent("success", selectorOverride, fx)
}

// Suspend Documentation is as below:
// Media data loading has been suspended.
// https://developer.mozilla.org/docs/Web/Events/suspend
/* This event provides options() to be called when the events is triggered and an optional selector which will override the internal selector mechanism of the domtrees.Element i.e if the selectorOverride argument is an empty string then domtrees.Element will create an appropriate selector matching its type and uid value in this format  (ElementType[uid='UID_VALUE']) but if the selector value is not empty then that becomes the default selector used
match the event with. */
func Suspend(fx domtrees.EventHandler, selectorOverride string) *domtrees.Event {
	return domtrees.NewEvent("suspend", selectorOverride, fx)
}

// TimeUpdate Documentation is as below:
// The time indicated by the currentTime attribute has been updated.
// https://developer.mozilla.org/docs/Web/Events/timeupdate
/* This event provides options() to be called when the events is triggered and an optional selector which will override the internal selector mechanism of the domtrees.Element i.e if the selectorOverride argument is an empty string then domtrees.Element will create an appropriate selector matching its type and uid value in this format  (ElementType[uid='UID_VALUE']) but if the selector value is not empty then that becomes the default selector used
match the event with. */
func TimeUpdate(fx domtrees.EventHandler, selectorOverride string) *domtrees.Event {
	return domtrees.NewEvent("timeupdate", selectorOverride, fx)
}

// Timeout Documentation is as below:
// (no documentation)
// https://developer.mozilla.org/docs/Web/Events/timeout
/* This event provides options() to be called when the events is triggered and an optional selector which will override the internal selector mechanism of the domtrees.Element i.e if the selectorOverride argument is an empty string then domtrees.Element will create an appropriate selector matching its type and uid value in this format  (ElementType[uid='UID_VALUE']) but if the selector value is not empty then that becomes the default selector used
match the event with. */
func Timeout(fx domtrees.EventHandler, selectorOverride string) *domtrees.Event {
	return domtrees.NewEvent("timeout", selectorOverride, fx)
}

// TouchCancel Documentation is as below:
// A touch point has been disrupted in an implementation-specific manners (too many touch points for example).
// https://developer.mozilla.org/docs/Web/Events/touchcancel
/* This event provides options() to be called when the events is triggered and an optional selector which will override the internal selector mechanism of the domtrees.Element i.e if the selectorOverride argument is an empty string then domtrees.Element will create an appropriate selector matching its type and uid value in this format  (ElementType[uid='UID_VALUE']) but if the selector value is not empty then that becomes the default selector used
match the event with. */
func TouchCancel(fx domtrees.EventHandler, selectorOverride string) *domtrees.Event {
	return domtrees.NewEvent("touchcancel", selectorOverride, fx)
}

// TouchEnd Documentation is as below:
// A touch point is removed from the touch surface.
// https://developer.mozilla.org/docs/Web/Events/touchend
/* This event provides options() to be called when the events is triggered and an optional selector which will override the internal selector mechanism of the domtrees.Element i.e if the selectorOverride argument is an empty string then domtrees.Element will create an appropriate selector matching its type and uid value in this format  (ElementType[uid='UID_VALUE']) but if the selector value is not empty then that becomes the default selector used
match the event with. */
func TouchEnd(fx domtrees.EventHandler, selectorOverride string) *domtrees.Event {
	return domtrees.NewEvent("touchend", selectorOverride, fx)
}

// TouchEnter Documentation is as below:
// A touch point is moved onto the interactive area of an element.
// https://developer.mozilla.org/docs/Web/Events/touchenter
/* This event provides options() to be called when the events is triggered and an optional selector which will override the internal selector mechanism of the domtrees.Element i.e if the selectorOverride argument is an empty string then domtrees.Element will create an appropriate selector matching its type and uid value in this format  (ElementType[uid='UID_VALUE']) but if the selector value is not empty then that becomes the default selector used
match the event with. */
func TouchEnter(fx domtrees.EventHandler, selectorOverride string) *domtrees.Event {
	return domtrees.NewEvent("touchenter", selectorOverride, fx)
}

// TouchLeave Documentation is as below:
// A touch point is moved off the interactive area of an element.
// https://developer.mozilla.org/docs/Web/Events/touchleave
/* This event provides options() to be called when the events is triggered and an optional selector which will override the internal selector mechanism of the domtrees.Element i.e if the selectorOverride argument is an empty string then domtrees.Element will create an appropriate selector matching its type and uid value in this format  (ElementType[uid='UID_VALUE']) but if the selector value is not empty then that becomes the default selector used
match the event with. */
func TouchLeave(fx domtrees.EventHandler, selectorOverride string) *domtrees.Event {
	return domtrees.NewEvent("touchleave", selectorOverride, fx)
}

// TouchMove Documentation is as below:
// A touch point is moved along the touch surface.
// https://developer.mozilla.org/docs/Web/Events/touchmove
/* This event provides options() to be called when the events is triggered and an optional selector which will override the internal selector mechanism of the domtrees.Element i.e if the selectorOverride argument is an empty string then domtrees.Element will create an appropriate selector matching its type and uid value in this format  (ElementType[uid='UID_VALUE']) but if the selector value is not empty then that becomes the default selector used
match the event with. */
func TouchMove(fx domtrees.EventHandler, selectorOverride string) *domtrees.Event {
	return domtrees.NewEvent("touchmove", selectorOverride, fx)
}

// TouchStart Documentation is as below:
// A touch point is placed on the touch surface.
// https://developer.mozilla.org/docs/Web/Events/touchstart
/* This event provides options() to be called when the events is triggered and an optional selector which will override the internal selector mechanism of the domtrees.Element i.e if the selectorOverride argument is an empty string then domtrees.Element will create an appropriate selector matching its type and uid value in this format  (ElementType[uid='UID_VALUE']) but if the selector value is not empty then that becomes the default selector used
match the event with. */
func TouchStart(fx domtrees.EventHandler, selectorOverride string) *domtrees.Event {
	return domtrees.NewEvent("touchstart", selectorOverride, fx)
}

// TransitionEnd Documentation is as below:
// A CSS transition has completed.
// https://developer.mozilla.org/docs/Web/Events/transitionend
/* This event provides options() to be called when the events is triggered and an optional selector which will override the internal selector mechanism of the domtrees.Element i.e if the selectorOverride argument is an empty string then domtrees.Element will create an appropriate selector matching its type and uid value in this format  (ElementType[uid='UID_VALUE']) but if the selector value is not empty then that becomes the default selector used
match the event with. */
func TransitionEnd(fx domtrees.EventHandler, selectorOverride string) *domtrees.Event {
	return domtrees.NewEvent("transitionend", selectorOverride, fx)
}

// Unload Documentation is as below:
// The document or a dependent resource is being unloaded.
// https://developer.mozilla.org/docs/Web/Events/unload
/* This event provides options() to be called when the events is triggered and an optional selector which will override the internal selector mechanism of the domtrees.Element i.e if the selectorOverride argument is an empty string then domtrees.Element will create an appropriate selector matching its type and uid value in this format  (ElementType[uid='UID_VALUE']) but if the selector value is not empty then that becomes the default selector used
match the event with. */
func Unload(fx domtrees.EventHandler, selectorOverride string) *domtrees.Event {
	return domtrees.NewEvent("unload", selectorOverride, fx)
}

// UpdateReady Documentation is as below:
// The resources listed in the manifest have been newly redownloaded, and the script can use swapCache() to switch to the new cache.
// https://developer.mozilla.org/docs/Web/Events/updateready
/* This event provides options() to be called when the events is triggered and an optional selector which will override the internal selector mechanism of the domtrees.Element i.e if the selectorOverride argument is an empty string then domtrees.Element will create an appropriate selector matching its type and uid value in this format  (ElementType[uid='UID_VALUE']) but if the selector value is not empty then that becomes the default selector used
match the event with. */
func UpdateReady(fx domtrees.EventHandler, selectorOverride string) *domtrees.Event {
	return domtrees.NewEvent("updateready", selectorOverride, fx)
}

// UpgradeNeeded Documentation is as below:
// An attempt was made to open a database with a version number higher than its current version. A versionchange transaction has been created.
// https://developer.mozilla.org/docs/Web/Reference/Events/upgradeneeded_indexedDB
/* This event provides options() to be called when the events is triggered and an optional selector which will override the internal selector mechanism of the domtrees.Element i.e if the selectorOverride argument is an empty string then domtrees.Element will create an appropriate selector matching its type and uid value in this format  (ElementType[uid='UID_VALUE']) but if the selector value is not empty then that becomes the default selector used
match the event with. */
func UpgradeNeeded(fx domtrees.EventHandler, selectorOverride string) *domtrees.Event {
	return domtrees.NewEvent("upgradeneeded", selectorOverride, fx)
}

// UserProximity Documentation is as below:
// Fresh data is available from a proximity sensor (indicates whether the nearby object is near the device or not).
// https://developer.mozilla.org/docs/Web/Events/userproximity
/* This event provides options() to be called when the events is triggered and an optional selector which will override the internal selector mechanism of the domtrees.Element i.e if the selectorOverride argument is an empty string then domtrees.Element will create an appropriate selector matching its type and uid value in this format  (ElementType[uid='UID_VALUE']) but if the selector value is not empty then that becomes the default selector used
match the event with. */
func UserProximity(fx domtrees.EventHandler, selectorOverride string) *domtrees.Event {
	return domtrees.NewEvent("userproximity", selectorOverride, fx)
}

// VersionChange Documentation is as below:
// A versionchange transaction completed.
// https://developer.mozilla.org/docs/Web/Reference/Events/versionchange_indexedDB
/* This event provides options() to be called when the events is triggered and an optional selector which will override the internal selector mechanism of the domtrees.Element i.e if the selectorOverride argument is an empty string then domtrees.Element will create an appropriate selector matching its type and uid value in this format  (ElementType[uid='UID_VALUE']) but if the selector value is not empty then that becomes the default selector used
match the event with. */
func VersionChange(fx domtrees.EventHandler, selectorOverride string) *domtrees.Event {
	return domtrees.NewEvent("versionchange", selectorOverride, fx)
}

// VisibilityChange Documentation is as below:
// The content of a tab has become visible or has been hidden.
// https://developer.mozilla.org/docs/Web/Events/visibilitychange
/* This event provides options() to be called when the events is triggered and an optional selector which will override the internal selector mechanism of the domtrees.Element i.e if the selectorOverride argument is an empty string then domtrees.Element will create an appropriate selector matching its type and uid value in this format  (ElementType[uid='UID_VALUE']) but if the selector value is not empty then that becomes the default selector used
match the event with. */
func VisibilityChange(fx domtrees.EventHandler, selectorOverride string) *domtrees.Event {
	return domtrees.NewEvent("visibilitychange", selectorOverride, fx)
}

// Voiceschanged Documentation is as below:
// The list of SpeechSynthesisVoice objects that would be returned by the SpeechSynthesis.getVoices() method has changed (when the voiceschanged event fires.)
// https://developer.mozilla.org/docs/Web/Events/voiceschanged
/* This event provides options() to be called when the events is triggered and an optional selector which will override the internal selector mechanism of the domtrees.Element i.e if the selectorOverride argument is an empty string then domtrees.Element will create an appropriate selector matching its type and uid value in this format  (ElementType[uid='UID_VALUE']) but if the selector value is not empty then that becomes the default selector used
match the event with. */
func Voiceschanged(fx domtrees.EventHandler, selectorOverride string) *domtrees.Event {
	return domtrees.NewEvent("voiceschanged", selectorOverride, fx)
}

// VolumeChange Documentation is as below:
// The volume has changed.
// https://developer.mozilla.org/docs/Web/Events/volumechange
/* This event provides options() to be called when the events is triggered and an optional selector which will override the internal selector mechanism of the domtrees.Element i.e if the selectorOverride argument is an empty string then domtrees.Element will create an appropriate selector matching its type and uid value in this format  (ElementType[uid='UID_VALUE']) but if the selector value is not empty then that becomes the default selector used
match the event with. */
func VolumeChange(fx domtrees.EventHandler, selectorOverride string) *domtrees.Event {
	return domtrees.NewEvent("volumechange", selectorOverride, fx)
}

// Waiting Documentation is as below:
// Playback has stopped because of a temporary lack of data.
// https://developer.mozilla.org/docs/Web/Events/waiting
/* This event provides options() to be called when the events is triggered and an optional selector which will override the internal selector mechanism of the domtrees.Element i.e if the selectorOverride argument is an empty string then domtrees.Element will create an appropriate selector matching its type and uid value in this format  (ElementType[uid='UID_VALUE']) but if the selector value is not empty then that becomes the default selector used
match the event with. */
func Waiting(fx domtrees.EventHandler, selectorOverride string) *domtrees.Event {
	return domtrees.NewEvent("waiting", selectorOverride, fx)
}

// Wheel Documentation is as below:
// A wheel button of a pointing device is rotated in any direction.
// https://developer.mozilla.org/docs/Web/Events/wheel
/* This event provides options() to be called when the events is triggered and an optional selector which will override the internal selector mechanism of the domtrees.Element i.e if the selectorOverride argument is an empty string then domtrees.Element will create an appropriate selector matching its type and uid value in this format  (ElementType[uid='UID_VALUE']) but if the selector value is not empty then that becomes the default selector used
match the event with. */
func Wheel(fx domtrees.EventHandler, selectorOverride string) *domtrees.Event {
	return domtrees.NewEvent("wheel", selectorOverride, fx)
}
