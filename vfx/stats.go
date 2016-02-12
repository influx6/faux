package vfx

import (
	"sync/atomic"
	"time"
)

//==============================================================================

// AnimationStepsPerSec defines the total steps taking per second of each clock
// tick.
var AnimationStepsPerSec = 60

//==============================================================================

// Stats defines a interface which holds stat information
// regarding the current frame and configuration for a sequence.
type Stats interface {
	Loop() bool
	Next(float64)
	Delta() float64
	Clone() Stats
	Easing() string
	IsDone() bool
	Reversed() bool
	Optimized() bool
	Reversible() bool
	CurrentIteration() int
	TotalIterations() int
	CompletedFirstTransition() bool
}

//==============================================================================

// Stat defines a the stats report strucuture for animation.
type Stat struct {
	totalIteration   int64
	currentIteration int64
	reversed         int64
	completed        int64
	completedReverse int64
	delta            float64
	loop             bool
	optimize         bool
	reversible       bool
	done             bool
	easing           string
}

// TimeStat returns a new Stats instance which provide information concering
// the current animation frame, it uses the provided duration to calculate the
// total iteration for the animation.
func TimeStat(ms time.Duration, easing string, loop, reversible, optimize bool) Stats {
	total := AnimationStepsPerSec * int(ms.Seconds())

	st := Stat{
		totalIteration: int64(total),
		loop:           loop,
		reversible:     reversible,
		easing:         easing,
		optimize:       optimize,
	}

	return &st
}

// MaxStat returns a new Stats using the provided numbers for animation.
func MaxStat(maxIteration int, easing string, loop, reversible, optimize bool) Stats {
	st := Stat{
		totalIteration: int64(maxIteration),
		loop:           loop,
		reversible:     reversible,
		optimize:       optimize,
		easing:         easing,
	}

	return &st
}

// Clone returns a clone for the stats.
func (s *Stat) Clone() Stats {
	st := Stat{
		totalIteration: int64(s.totalIteration),
		loop:           s.loop,
		reversible:     s.reversible,
		optimize:       s.optimize,
		easing:         s.easing,
	}

	return &st
}

// Easing returns the easing value for this specifc stat.
func (s *Stat) Easing() string {
	return s.easing
}

// Delta returns the current time delta from the last update.
func (s *Stat) Delta() float64 {
	return s.delta
}

// Next calls the appropriate iteration step for the stat.
func (s *Stat) Next(m float64) {
	if !s.CompletedFirstTransition() {
		s.NextIteration(m)
		return
	}

	if s.Reversible() {
		atomic.StoreInt64(&s.reversed, 1)
		s.PreviousIteration(m)
		return
	}
}

// NextIteration increments the iteration count.
func (s *Stat) NextIteration(m float64) {
	atomic.AddInt64(&s.currentIteration, 1)

	it := atomic.LoadInt64(&s.totalIteration)
	ct := atomic.LoadInt64(&s.currentIteration)

	if ct >= it {
		atomic.StoreInt64(&s.completed, 1)
	}

	s.delta = m
}

// PreviousIteration increments the iteration count.
func (s *Stat) PreviousIteration(m float64) {
	atomic.AddInt64(&s.currentIteration, -1)

	ct := atomic.LoadInt64(&s.currentIteration)

	if ct <= 0 {
		atomic.StoreInt64(&s.completedReverse, 1)
	}

	s.delta = m
}

// Optimized returns true/false if the stat is said to use optimization
// strategies.
func (s *Stat) Optimized() bool {
	return s.optimize
}

// IsDone returns true/false if the stat is done.
func (s *Stat) IsDone() bool {
	ct := atomic.LoadInt64(&s.completed)

	if !s.Reversible() {
		if ct <= 0 {
			return false
		}

		return true
	}

	rs := atomic.LoadInt64(&s.completedReverse)

	if ct > 0 && rs <= 0 {
		return false
	}

	return true
}

// Reversed returns true/false if the stats has entered a reversed state.
func (s *Stat) Reversed() bool {
	return atomic.LoadInt64(&s.reversed) > 0
}

// Reversible returns true/false if the stat animation is set to loop.
func (s *Stat) Reversible() bool {
	return s.reversible
}

// CompletedFirstTransition returns true/false if the stat has completed a full
// iteration of its total iteration, this is useful to know when loop or
// reversal is turned on, to check if this stat has entered its looping or
// reversal state. It only reports for the first completion of total iterations.
func (s *Stat) CompletedFirstTransition() bool {
	return atomic.LoadInt64(&s.completed) > 0
}

// Loop returns true/false if the stat animation is set to loop.
func (s *Stat) Loop() bool {
	return s.loop
}

// TotalIterations returns the total iteration for this specific stat.
func (s *Stat) TotalIterations() int {
	return int(atomic.LoadInt64(&s.totalIteration))
}

// CurrentIteration returns the current iteration for this specific stat.
func (s *Stat) CurrentIteration() int {
	return int(atomic.LoadInt64(&s.currentIteration))
}

//==============================================================================
