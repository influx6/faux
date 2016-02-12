package vfx

// Stats defines a interface which holds stat information
// regarding the current frame and configuration for a sequence.
type Stats interface {
	Loop() bool
	Easing() string
	Reversible() bool
	Optimized() bool
	Reversed() bool
	IsDone() bool
	Delta() float64
	Clone() Stats
	CurrentIteration() int
	TotalIterations() int
	NextIteration(float64)
	PreviousIteration(float64)
}

// WriterCache provides a interface type for writer cache structures, which catch
// animation produced writers per sequence iteration state.
type WriterCache interface {
	Store(Stats, int, ...DeferWriter)
	Writers(Stats, int) DeferWriters
	ClearIteration(Stats, int)
	Clear(Stats)
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

// Frame defines the interface for a animation sequence generator,
// it defines the sequence of a organized step for animation.
type Frame interface {
	Init() DeferWriters
	Inited() bool
	Sequence() DeferWriters
	Stats() Stats
	IsOver() bool
	End()
}

// ResetableFrame defines a frame that can be reset to its default stat
type ResetableFrame interface {
	resetStats()
}
