package vfx

//==============================================================================

// EaseIn provides a struct for 'easing-in' based animation.
type EaseIn struct{}

// Ease returns a new value base on the EaseConfig received.
func (e EaseIn) Ease(c EaseConfig) float64 {
	ms := float64(c.CurrentStep) / float64(c.TotalSteps)
	return (c.DeltaValue * ms) + c.CurrentValue
}

func init() {
	RegisterEasing("ease-in", EaseIn{})
}

//==============================================================================

// EaseInQuad provides a struct for 'easing-in-quad' based animation.
type EaseInQuad struct{}

// Ease returns a new value base on the EaseConfig received.
func (e EaseInQuad) Ease(c EaseConfig) float64 {
	ms := float64(c.CurrentStep) / float64(c.TotalSteps)
	return (c.DeltaValue * ms * ms) + c.CurrentValue
}

func init() {
	RegisterEasing("ease-in-quad", EaseInQuad{})
}

//==============================================================================

// EaseOutQuad provides a struct for 'easing-out-quad' based animation.
type EaseOutQuad struct{}

// Ease returns a new value base on the EaseConfig received.
func (e EaseOutQuad) Ease(c EaseConfig) float64 {
	ms := (float64(c.CurrentStep) / float64(c.TotalSteps)) * float64(c.CurrentStep-2)
	return ((c.DeltaValue * -1) * ms) + c.CurrentValue
}

func init() {
	RegisterEasing("ease-out-quad", EaseOutQuad{})
}

//==============================================================================

// EaseInOutQuad provides a struct for 'easing-in-out-quad' based animation.
type EaseInOutQuad struct{}

// Ease returns a new value base on the EaseConfig received.
func (e EaseInOutQuad) Ease(c EaseConfig) float64 {
	diff := (float64(c.CurrentStep) / float64(c.TotalSteps))

	if diff < 1 {
		return (c.DeltaValue / 2) * diff * diff
	}

	diff--

	return (-1*c.DeltaValue)*((diff)*(diff-2)-1) + c.CurrentValue
}

func init() {
	RegisterEasing("ease-in-out-quad", EaseInOutQuad{})
}

//==============================================================================
