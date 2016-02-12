package vfx

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

// FramePhase defines a animation phase type.
type FramePhase int

// const contains sets of Frame phase that identify the current frame animation
// phase.
const (
	NOPHASE FramePhase = iota
	STARTPHASE
	OPTIMISEPHASE
)

// Frame defines the interface for a animation sequence generator,
// it defines the sequence of a organized step for animation.
type Frame interface {
	End()
	Sync()
	Stats() Stats
	Inited() bool
	IsOver() bool
	Init() DeferWriters
	Phase() FramePhase
	Sequence() DeferWriters
}

//==============================================================================

// WriterCache provides a interface type for writer cache structures, which catch
// animation produced writers per sequence iteration state.
type WriterCache interface {
	Store(Frame, int, ...DeferWriter)
	Writers(Frame, int) DeferWriters
	ClearIteration(Frame, int)
	Clear(Frame)
}

// DeferWriter provides an interface that allows deferring the write effects
// of a sequence, these way we can collate all effects of a set of sequences
// together to perform a batch right, reducing layout trashing.
type DeferWriter interface {
	Write()
}

// DeferWriters defines a lists of DeferWriter implementing structures.
type DeferWriters []DeferWriter

// CascadeDeferWriter provides a specifc writer that can combine multiple writers
// into a single one but also allow flatten its writer sequence into a ordered
// lists of DeferWriters.
type CascadeDeferWriter interface {
	DeferWriter
	Flatten() DeferWriters
}

// Sequence defines a series of animation step which will be runned
// through by calling its .Next() method continousely until the
// sequence is done(if its not a repetitive sequence).
// Sequence when calling their next method, all sequences must return a
// DeferWriter.
type Sequence interface {
	Init(Stats) DeferWriter
	Next(Stats) DeferWriter
	IsDone() bool
}

// SequenceList defines a lists of animatable sequence.
type SequenceList []Sequence
