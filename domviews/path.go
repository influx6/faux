package domviews

import (
	"errors"
	"fmt"
	"strings"

	"github.com/go-humble/detect"
	"github.com/gopherjs/gopherjs/js"
	"github.com/influx6/faux/domevents"
	"github.com/influx6/faux/pub"
)

// PathSpec represent the current path and hash values
type PathSpec struct {
	Hash     string
	Path     string
	Sequence string
}

// String returns the hash and path
func (p *PathSpec) String() string {
	return fmt.Sprintf("%s%s", p.Path, p.Hash)
}

// PathSequencer provides a function to convert either the path/hash into a
// dot seperated sequence string for use with States.
type PathSequencer func(path string, hash string) string

// HashSequencer provides a PathSequencer that returns the hash part of a url,
// as the path sequence.
func HashSequencer(path, hash string) string {
	cleanHash := strings.Replace(hash, "#", ".", -1)
	return strings.Replace(cleanHash, "/", ".", -1)
}

// URLPathSequencer provides a PathSequencer that returns the path part of a url,
// as the path sequence.
func URLPathSequencer(path, hash string) string {
	return strings.Replace(path, "/", ".", -1)
}

// PathObserver represent any continouse changing route path by the browser
type PathObserver struct {
	pub.Publisher
	usingHash bool
	sequencer PathSequencer
}

// Path returns a new PathObserver instance
func Path(ps PathSequencer) *PathObserver {
	if ps == nil {
		ps = HashSequencer
	}
	return &PathObserver{
		Publisher: pub.Identity(),
		sequencer: ps,
	}
}

// Follow creates a Pathspec from the hash and path and sends it
func (p *PathObserver) Follow(path, hash string) {
	p.FollowSpec(PathSpec{Hash: hash, Path: path, Sequence: p.sequencer(path, hash)})
}

// FollowSpec just passes down the Pathspec
func (p *PathObserver) FollowSpec(ps PathSpec) {
	p.Send(ps)
}

// NotifyPage is used to notify a page of
func (p *PathObserver) NotifyPage(pg *Pages) {
	p.React(func(r pub.Publisher, _ error, d interface{}) {
		if ps, ok := d.(PathSpec); ok {
			pg.All(ps.Sequence)
		}
	}, true)
}

// NotifyPartialPage is used to notify a Page using the page engine's Partial() activator
func (p *PathObserver) NotifyPartialPage(pg *Pages) {
	p.React(func(r pub.Publisher, _ error, d interface{}) {
		// if err != nil { r.SendError(err) }
		if ps, ok := d.(PathSpec); ok {
			pg.Partial(ps.Sequence)
		}
	}, true)
}

// HashChangePath returns a path observer path changes
func HashChangePath(ps PathSequencer) *PathObserver {
	panicBrowserDetect()
	path := Path(ps)
	path.usingHash = true

	js.Global.Set("onhashchange", func() {
		path.Follow(GetLocation())
	})

	return path
}

//ErrNotSupported is returned when a feature requested is not supported by the environment
var ErrNotSupported = errors.New("Feature not supported")

// PopStatePath returns a path observer path changes
func PopStatePath(ps PathSequencer) (*PathObserver, error) {
	panicBrowserDetect()

	if !BrowserSupportsPushState() {
		return nil, ErrNotSupported
	}

	path := Path(ps)

	js.Global.Set("onpopstate", func() {
		path.Follow(GetLocation())
	})

	return path, nil
}

// HistoryProvider wraps the PathObserver with methods that allow easy control of
// client location
type HistoryProvider struct {
	*PathObserver
}

// History returns a new PathObserver and depending on browser support will either use the
// popState or HashChange
func History(ps PathSequencer) *HistoryProvider {
	pop, err := PopStatePath(ps)

	if err != nil {
		pop = HashChangePath(ps)
	}

	return &HistoryProvider{pop}
}

// Go changes the path of the current browser location depending on wether its underline
// observe is hashed based or pushState based,it will use SetDOMHash or PushDOMState appropriately
func (h *HistoryProvider) Go(path string) {
	if h.usingHash {
		SetDOMHash(path)
		return
	}
	PushDOMState(path)
}

// ErrBadSelector is used to indicate if the selector returned no result
var ErrBadSelector = errors.New("Selector returned nil")

// Pages provides the concrete provider for managing a whole website or View
// you dont need two,just one is enough to manage the total web view of your app / site
// It ties directly into the page hash or popstate location to provide consistent updates
type Pages struct {
	*StateEngine
	*HistoryProvider
	// views []Views
}

// Page returns the new state engine powered page
func Page(ps PathSequencer) *Pages {
	return NewPage(History(ps))
}

// NewPage returns the new state engine powered page
func NewPage(p *HistoryProvider) *Pages {
	pg := &Pages{
		StateEngine:     NewStateEngine(),
		HistoryProvider: p,
	}

	p.NotifyPage(pg)
	pg.All(".")
	return pg
}

// Mount adds a component into the page for handling/managing of visiblity and
// gets the dom referenced by the selector using QuerySelector and returns an error if selector gave no result
func (p *Pages) Mount(selector, addr string, v Views) error {
	n := domevents.GetDocument().Call("querySelector", selector)

	if n == nil || n == js.Undefined {
		return ErrBadSelector
	}

	p.AddView(addr, v)
	v.Mount(n)
	return nil
}

// AddView adds a view to the page
func (p *Pages) AddView(addr string, v Views) {
	v.UseHistory(p.HistoryProvider)
	p.UseState(addr, v)
}

// Address returns the current path and hash of the location api
func (p *Pages) Address() (string, string) {
	return GetLocation()
}

// GetLocation returns the path and hash of the browsers location api else panics if not in a browser
func GetLocation() (string, string) {
	panicBrowserDetect()
	loc := js.Global.Get("location")
	path := loc.Get("pathname").String()
	hash := loc.Get("hash").String()
	return path, hash
}

// PushDOMState adds a new state the dom push history
func PushDOMState(path string) {
	panicBrowserDetect()
	js.Global.Get("history").Call("pushState", nil, "", path)
}

// SetDOMHash sets the dom location hash
func SetDOMHash(hash string) {
	panicBrowserDetect()
	js.Global.Get("location").Set("hash", hash)
}

func panicBrowserDetect() {
	if !detect.IsBrowser() {
		panic("expected to be used in a dom/browser env")
	}
}

// BrowserSupportsPushState checks if browser supports pushState
func BrowserSupportsPushState() bool {
	if !detect.IsBrowser() {
		return false
	}

	return (js.Global.Get("onpopstate") != js.Undefined) &&
		(js.Global.Get("history") != js.Undefined) &&
		(js.Global.Get("history").Get("pushState") != js.Undefined)
}
