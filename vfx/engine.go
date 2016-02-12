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
	engine = loop.New(gear)
}

//==============================================================================

// Animate provides the central engine for managing all animation calls,
// it returns a subscription interface that allows the animation to be stopped.
// Animate uses a writer batching to reduce layout trashing. Hence multiple
// frames assigned for each animation call, will have all their writes, batched
// into one call.
func Animate(f ...Frame) loop.Looper {
	// Return this frame subscription ender, initialized and run its writers.
	return engine.Loop(func(delta float64) {
		var writers DeferWriters

		for _, frame := range f {
			if !frame.Inited() {
				writers = append(writers, frame.Init()...)
				frame.Stats().NextIteration(delta)
				continue
			}

			if frame.IsOver() {
				continue
			}

			writers = append(writers, frame.Sequence()...)
			frame.Stats().NextIteration(delta)
		}

		// batch all the writes together as one.
		for _, w := range writers {
			w.Write()
		}
	}, 0)
}

//==============================================================================
