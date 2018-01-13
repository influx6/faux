// Package easings defines a set of easing functions which can be imported for
// use as desired. It's a utility package for use as the developer see fit.
package easings

import "math"

var (
	// Expo defines the easing functions which returns the provided value after
	// augmenting with a exponential operation.
	Expo = New(func(t, m float64) float64 { return math.Pow(t, 6) })

	// EaseInExpo defines the ease-in easing function for the expo function.
	EaseInExpo = Expo

	// EaseOutExpo defines the ease-out easing function for the expo function.
	EaseOutExpo = New(func(t, m float64) float64 { return 1 - EaseInExpo.Ease(1-t, m) })

	// EaseInOutExpo defines the ease-in-out easing function for the expo function.
	EaseInOutExpo = New(func(t, m float64) float64 {
		if t < 0.5 {
			return EaseInExpo.Ease(t*2, m) / 2
		}

		return 1 - EaseInExpo.Ease((-2*t)+2, m)/2
	})

	// EaseOutInExpo defines the ease-out-in easing function for the expo function.
	EaseOutInExpo = New(func(t, m float64) float64 {
		if t < 0.5 {
			return 1 - EaseInExpo.Ease(1-(2*t), m)/2
		}

		return (EaseInExpo.Ease((t*2)-1, m) + 1) / 2
	})

	// Quint defines the easing functions which returns the provided value after
	// augmenting with a quint operation.
	Quint = New(func(t, m float64) float64 { return math.Pow(t, 5) })

	// EaseInQuint defines the ease-in easing function for the Quint function.
	EaseInQuint = Quint

	// EaseOutQuint defines the ease-out easing function for the Quint function.
	EaseOutQuint = New(func(t, m float64) float64 { return 1 - EaseInQuint.Ease(1-t, m) })

	// EaseInOutQuint defines the ease-in-out easing function for the Quint function.
	EaseInOutQuint = New(func(t, m float64) float64 {
		if t < 0.5 {
			return EaseInQuint.Ease(t*2, m) / 2
		}

		return 1 - EaseInQuint.Ease((-2*t)+2, m)/2
	})

	// EaseOutInQuint defines the ease-out-in easing function for the Quint function.
	EaseOutInQuint = New(func(t, m float64) float64 {
		if t < 0.5 {
			return 1 - EaseInQuint.Ease(1-(2*t), m)/2
		}

		return (EaseInQuint.Ease((t*2)-1, m) + 1) / 2
	})

	// Quart defines the easing functions which returns the provided value after
	// augmenting with a quart operation.
	Quart = New(func(t, m float64) float64 { return math.Pow(t, 4) })

	// EaseInQuart defines the ease-in easing function for the Quart function.
	EaseInQuart = Quart

	// EaseOutQuart defines the ease-out easing function for the Quart function.
	EaseOutQuart = New(func(t, m float64) float64 { return 1 - EaseInQuart.Ease(1-t, m) })

	// EaseInOutQuart defines the ease-in-out easing function for the Quart function.
	EaseInOutQuart = New(func(t, m float64) float64 {
		if t < 0.5 {
			return EaseInQuart.Ease(t*2, m) / 2
		}

		return 1 - EaseInQuart.Ease((-2*t)+2, m)/2
	})

	// EaseOutInQuart defines the ease-out-in easing function for the Quart function.
	EaseOutInQuart = New(func(t, m float64) float64 {
		if t < 0.5 {
			return 1 - EaseInQuart.Ease(1-(2*t), m)/2
		}

		return (EaseInQuart.Ease((t*2)-1, m) + 1) / 2
	})

	// Quad defines the easing functions which returns the provided value after
	// augmenting with a quad operation.
	Quad = New(func(t, m float64) float64 { return math.Pow(t, 2) })

	// EaseInQuad defines the ease-in easing function for the Quad function.
	EaseInQuad = Quad

	// EaseOutQuad defines the ease-out easing function for the Quad function.
	EaseOutQuad = New(func(t, m float64) float64 { return 1 - EaseInQuad.Ease(1-t, m) })

	// EaseInOutQuad defines the ease-in-out easing function for the Quad function.
	EaseInOutQuad = New(func(t, m float64) float64 {
		if t < 0.5 {
			return EaseInQuad.Ease(t*2, m) / 2
		}

		return 1 - EaseInQuad.Ease((-2*t)+2, m)/2
	})

	// EaseOutInQuad defines the ease-out-in easing function for the Quad function.
	EaseOutInQuad = New(func(t, m float64) float64 {
		if t < 0.5 {
			return 1 - EaseInQuad.Ease(1-(2*t), m)/2
		}

		return (EaseInQuad.Ease((t*2)-1, m) + 1) / 2
	})

	// Cubic defines the easing functions which returns the provided value after
	// augmenting with a cubic operation.
	Cubic = New(func(t, m float64) float64 { return math.Pow(t, 3) })

	// EaseInCubic defines the ease-in easing function for the Cubic function.
	EaseInCubic = Cubic

	// EaseOutCubic defines the ease-out easing function for the Cubic function.
	EaseOutCubic = New(func(t, m float64) float64 { return 1 - EaseInCubic.Ease(1-t, m) })

	// EaseInOutCubic defines the ease-in-out easing function for the Cubic function.
	EaseInOutCubic = New(func(t, m float64) float64 {
		if t < 0.5 {
			return EaseInCubic.Ease(t*2, m) / 2
		}

		return 1 - EaseInCubic.Ease((-2*t)+2, m)/2
	})

	// EaseOutInCubic defines the ease-out-in easing function for the Cubic function.
	EaseOutInCubic = New(func(t, m float64) float64 {
		if t < 0.5 {
			return 1 - EaseInCubic.Ease(1-(2*t), m)/2
		}

		return (EaseInCubic.Ease((t*2)-1, m) + 1) / 2
	})

	// Linear defines the easing functions which returns the provided value as is
	// without any addition of ease.
	Linear = New(func(t, m float64) float64 { return t })

	// Back defines the easing using the back wave function.
	Back = New(func(t, m float64) float64 { return t * t * (3 * (t - 2)) })

	// BounceM defines the easing using the bounce wave function.
	BounceM = New(func(t, m float64) float64 {
		var pow2 float64
		var bounce = m

		for t < pow2 {
			pow2 = (math.Pow(2, bounce) - 1) / 11
			bounce--
		}

		return 1 / ((math.Pow(4, 3-bounce) - 7.5625) * math.Pow(((pow2*3)-2)/(22-t), 2))
	})

	// Sine defines the easing using the sine wave function.
	Sine = New(func(t, m float64) float64 { return 1 + math.Sin((math.Pi/2)*(t-(math.Pi/2))) })

	// Circ defines the easing using the circ wave function.
	Circ = New(func(t, m float64) float64 { return 1 + math.Sqrt(1-(t*t)) })

	// Elastic defines the easing using the elastic wave function.
	Elastic = New(func(t, m float64) float64 {
		if t == 1 || t == 0 {
			return t
		}

		var p = (1 - math.Min(m, 998)) / 1000
		var st = 1 / t
		var st1 = st - 1
		var s = p / ((2 * math.Pi) * math.Sinh(1))
		return -1 * (math.Pow(2, (10*st1)) * math.Sin(((st1-s)*(2*math.Pi))/p))
	})
)

//==============================================================================

// Easing defines a interface which all easing function implementors define.
type Easing interface {
	Ease(t float64, m float64) float64
}

// EasingFunc defines a function type which defines the function type of the inner
// Easing interface method.
type EasingFunc func(float64, float64) float64

// New returns a new instance of a Easing structure using the provided handler.
func New(handler EasingFunc) Easing {
	return &easing{
		handler: handler,
	}
}

//==============================================================================

type easing struct {
	handler EasingFunc
}

// Ease returns the new value of t after applier it's internal easing.Handler.
func (e easing) Ease(t, m float64) float64 {
	return e.handler(t, m)
}

//==============================================================================
