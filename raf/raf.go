package raf

import (
	"math"
	"sync"
	"sync/atomic"
	"time"

	"github.com/go-humble/detect"
	"github.com/gopherjs/gopherjs/js"
)

// Mux defines a handler for using with RAF.
type Mux func(float64)

//==============================================================================

// Clock defines an interface which exposes methods which allow a timeloop run.
type Clock struct {
	mux     Mux
	paused  int64
	clockID int64
}

// New returns a new instance pointer of the Clock type.
func New(m Mux) *Clock {
	return &Clock{
		mux:     m,
		clockID: -1,
		paused:  -1,
	}
}

// Start registers the clock with the animation call loop. Calls all passed in
// functions once the clock has being successfully registered.
func (c *Clock) Start(f ...func()) {
	if atomic.LoadInt64(&c.clockID) != -1 {
		return
	}

	atomic.StoreInt64(&c.clockID, int64(RequestAnimationFrame(c.Tick, f...)))
}

// Stop deregisters the clock and stops all loop calls and calls the passed in
// functions.
func (c *Clock) Stop(f ...func()) {
	if atomic.LoadInt64(&c.clockID) == -1 {
		return
	}

	CancelAnimationFrame(int(atomic.LoadInt64(&c.clockID)), f...)
	atomic.StoreInt64(&c.clockID, -1)
}

// Tick enables the clocks ticking if it has been paused.
func (c *Clock) Tick(f float64) {
	if atomic.LoadInt64(&c.paused) != -1 {
		return
	}

	go func() {
		c.mux(f)

		atomic.StoreInt64(&c.clockID, int64(RequestAnimationFrame(c.Tick)))
	}()
}

// Toggle switches the state of the clock from paused to resume and vise versa.
func (c *Clock) Toggle() {
	if atomic.LoadInt64(&c.paused) < 0 {
		atomic.StoreInt64(&c.paused, 1)
		return
	}

	atomic.StoreInt64(&c.paused, -1)
}

// Resume enables the clocks ticking if it has been paused.
func (c *Clock) Resume() {
	atomic.StoreInt64(&c.paused, -1)
}

// Pause disabbles the clocks ticking if it has been paused.
func (c *Clock) Pause() {
	atomic.StoreInt64(&c.paused, 1)
}

//==============================================================================

var minElapse = 16 * time.Millisecond

var tickers = struct {
	m       sync.Mutex
	tickers []*ticker
}{}

type ticker struct {
	id   int
	fn   Mux
	kill bool
	last float64
}

// RequestAnimationFrame provides a cover for RAF using the js
// api for requestAnimationFrame.
func RequestAnimationFrame(r Mux, f ...func()) int {
	if !detect.IsBrowser() {
		var tLen int

		tickers.m.Lock()
		{
			tLen = len(tickers.tickers)

			newticker := &ticker{
				id:   tLen,
				fn:   r,
				last: float64(time.Now().Unix()),
			}

			go func() {
				for {
					if newticker.kill {
						break
					}

					now := float64(time.Now().Unix())
					elapsed := math.Max(0, now-newticker.last)
					newticker.last = now

					tickerElapsed := time.Duration(float64(minElapse) + elapsed)

					tk := time.NewTicker(tickerElapsed)
					<-tk.C

					newticker.fn(tickerElapsed.Seconds())
				}
			}()

			tickers.tickers = append(tickers.tickers, newticker)
		}
		tickers.m.Unlock()

		return tLen
	}

	id := js.Global.Call("requestAnimationFrame", r).Int()

	for _, fx := range f {
		fx()
	}

	return id
}

// CancelAnimationFrame provides a cover for RAF using the
// js api cancelAnimationFrame.
func CancelAnimationFrame(id int, f ...func()) {
	if !detect.IsBrowser() {
		tickers.m.Lock()
		{
			tLen := len(tickers.tickers)
			if id >= tLen {
				tickers.m.Unlock()
				return
			}

			tm := tickers.tickers[id]
			tm.kill = true
		}
		tickers.m.Unlock()

		return
	}

	js.Global.Call("cancelAnimationFrame", id)

	for _, fx := range f {
		fx()
	}
}

//==============================================================================

func init() {
	if detect.IsBrowser() {
		rafPolyfill()
	}
}

func rafPolyfill() {
	window := js.Global
	vendors := []string{"ms", "moz", "webkit", "o"}
	if window.Get("requestAnimationFrame") == nil {
		for i := 0; i < len(vendors) && window.Get("requestAnimationFrame") == nil; i++ {
			vendor := vendors[i]
			window.Set("requestAnimationFrame", window.Get(vendor+"RequestAnimationFrame"))
			window.Set("cancelAnimationFrame", window.Get(vendor+"CancelAnimationFrame"))
			if window.Get("cancelAnimationFrame") == nil {
				window.Set("cancelAnimationFrame", window.Get(vendor+"CancelRequestAnimationFrame"))
			}
		}
	}

	lastTime := 0.0
	if window.Get("requestAnimationFrame") == nil {
		window.Set("requestAnimationFrame", func(callback func(float64)) int {
			currTime := js.Global.Get("Date").New().Call("getTime").Float()
			timeToCall := math.Max(0, 16-(currTime-lastTime))
			id := window.Call("setTimeout", func() { callback(float64(currTime + timeToCall)) }, timeToCall)
			lastTime = currTime + timeToCall
			return id.Int()
		})
	}

	if window.Get("cancelAnimationFrame") == nil {
		window.Set("cancelAnimationFrame", func(id int) {
			js.Global.Get("clearTimeout").Invoke(id)
		})
	}
}

//==============================================================================
