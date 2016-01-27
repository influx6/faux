package vfx

// Stats defines a interface which holds stat information
// regarding the current frame and configuration for a sequence.
type Stats interface {
	CurrentIteration() int
	TotalIterations() int
	Stamp() float64
	Loop() bool
	Reversible() bool
	IsDone() bool
	Clone() Stats
}

// Sequence defines a series of animation step which will be runned
// through by calling its .Next() method continousely until the
// sequence is done(if its not a repetitive sequence).
type Sequence interface {
	Next(Stats)
	IsDone() bool
}

// SequenceList defines a lists of animatable sequence.
type SequenceList []Sequence

// Frame defines the interface for a animation sequence generator,
// it defines the sequence of a organized step for animation.
type Frame interface {
	Sequence() SequenceList
	Stat() Stats
}

// Frames defines a type for a building animation Frame.
type Frames interface {
	Build() Frame
}

// NewFrameBuilder defines a builder for building a animation frame.
func NewFrameBuilder(stat Stats, s ...Sequence) Frames {
	return nil
}
