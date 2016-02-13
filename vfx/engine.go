package vfx

import "github.com/influx6/faux/loop"

//==============================================================================

// engine is the global gameloop engine used in managing animations within the
// global loop.
var engine loop.GameEngine

// wcache contains all writers cache with respect to each stats for a specific
// frame.
var wcache WriterCache

//==============================================================================

// Init initializes the animation system with the necessary loop engine,
// desired to be used in running the animation.
// Important: This should only be called once.
func Init(gear loop.EngineGear) {
	wcache = NewDeferWriterCache()
	engine = loop.New(gear)
}

//==============================================================================

// GetWriterCache returns the writer cache used by the animation library.
func GetWriterCache() WriterCache {
	return wcache
}

//==============================================================================

// Animate provides the central engine for managing all animation calls,
// it returns a subscription interface that allows the animation to be stopped.
// Animate uses a writer batching to reduce layout trashing. Hence multiple
// frames assigned for each animation call, will have all their writes, batched
// into one call.
func Animate(frame Frame) loop.Looper {

	var engineStopper loop.Looper

	// Return this frame subscription ender, initialized and run its writers.
	engineStopper = engine.Loop(func(delta float64) {
		var writers DeferWriters

		if frame.IsOver() {
			wcache.Clear(frame)

			// If we are over and the stopper is not nil then use it to stop our
			// frame looper.
			if engineStopper != nil {
				engineStopper.End()
			}

			return
		}

		// stats := frame.Stats()

		if !frame.Inited() {
			writers = frame.Init(delta)

			// stats.Next(delta)
			frame.Sync()
		} else {
			writers = frame.Sequence(delta)

			// stats.Next(delta)
			frame.Sync()
		}

		// Incase we end up using delays with our sequence, GopherJS can
		// not block and should not block, other processes, so lunch the
		// writers in a Goroutine. Frames have built in reconciliation system
		// to manage the variances when dealing with delays.
		go func() {
			for _, w := range writers {
				w.Write()
			}
		}()

	}, 0)

	return engineStopper
}

//==============================================================================
