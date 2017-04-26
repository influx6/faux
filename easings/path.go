package easings

import "math"

// PathCurve defines a type which exposes a function to retrieve a point over a
// curves length. Useful for retrieving easing values over a curve.
type PathCurve interface {
	GetPointAtLength(step float64) (float64, float64)
}

// SimplePath defines a struct which implements the easing function and returns
// a giving value of a certain point.
type SimplePath struct {
	Curve PathCurve
}

// PointTransition is a struct which contains set of transition points from
// a giving set of ease value.
type PointTransition struct {
	TranslateX float64
	TranslateY float64
	Rotate     float64
}

// Point uses the provided value to return a PointTransition that provides an
// adequate ease set for translation and rotation.
func (s *SimplePath) Point(value float64, progress float64) PointTransition {
	px, py := s.EaseValue(value, 0, progress)
	p0x, p0y := s.MinusEase(value, progress)
	p1x, p1y := s.PlusEase(value, progress)

	return PointTransition{
		TranslateX: px,
		TranslateY: py,
		Rotate:     math.Atan2((p1y-p0y), (p1x-p0x)) * (180 / math.Pi),
	}
}

// MinusEase returns the value of the ease on the negative/left side of the curve.
func (s *SimplePath) MinusEase(value float64, progress float64) (float64, float64) {
	return s.EaseValue(value, -1, progress)
}

// PlusEase returns the value of the ease on the right side of the curve.
func (s *SimplePath) PlusEase(value float64, progress float64) (float64, float64) {
	return s.EaseValue(value, +1, progress)
}

// EaseValue returns the giving easing value of the provided offset value.
func (s *SimplePath) EaseValue(value float64, offset, progress float64) (float64, float64) {
	var p float64

	if progress > 1 {
		p = value + offset
	} else {
		p = (value * progress) + offset
	}

	return s.Curve.GetPointAtLength(p)
}
