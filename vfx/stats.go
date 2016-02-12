package vfx

import (
	"sync/atomic"
	"time"
)

//==============================================================================

// Stat defines a the stats report strucuture for animation.
type Stat struct {
	totalIteration   int64
	currentIteration int64
	reversed         int64
	completed        int64
	delta            float64
	loop             bool
	reversible       bool
	done             bool
	easing           string
}

// AnimationStepsPerSec defines the total steps taking per second of each clock
// tick.
var AnimationStepsPerSec = 60

// TimeStat returns a new Stats instance which provide information concering
// the current animation frame, it uses the provided duration to calculate the
// total iteration for the animation.
func TimeStat(ms time.Duration, easing string, loop, reversible bool) Stats {
	st := Stat{
		loop:       loop,
		reversible: reversible,
		easing:     easing,
	}

	total := AnimationStepsPerSec * int(ms.Seconds())
	atomic.StoreInt64(&st.totalIteration, int64(total))
	return &st
}

// MaxStat returns a new Stats using the provided numbers for animation.
func MaxStat(maxIteration int, easing string, loop, reversible bool) Stats {
	st := Stat{
		totalIteration: int64(maxIteration),
		loop:           loop,
		reversible:     reversible,
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
		atomic.StoreInt64(&s.reversed, 1)
	}

	s.delta = m
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

	rs := atomic.LoadInt64(&s.reversed)

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
