package vfx

import "sync/atomic"

// NewAnimationSequence defines a builder for building a animation frame.
func NewAnimationSequence(stat Stats, s ...Sequence) Frame {
	as := AnimationSequence{sequences: s, stat: stat}
	return &as
}

// AnimationSequence defines a set of sequences that operate on the behaviour of
// a dom element or lists of dom.elements.
type AnimationSequence struct {
	sequences  SequenceList
	stat       Stats
	inited     int64
	iniWriters DeferWriters
}

// Inited returns true/false if the frame has begun.
func (f *AnimationSequence) Inited() bool {
	return atomic.LoadInt64(&f.inited) > 1
}

// Init calls the initialization writers for each sequence, returning their
// respective initialization writers if any to be runned on the first loop.
func (f *AnimationSequence) Init() DeferWriters {
	if atomic.LoadInt64(&f.inited) > 1 {
		return f.iniWriters
	}

	var writers DeferWriters

	// Collect all writers from each sequence with in the frame.
	for _, seq := range f.sequences {
		isq, ok := seq.(InitableSequence)
		if !ok {
			continue
		}

		writers = append(writers, isq.Init(f.Stats()))
	}

	atomic.StoreInt64(&f.inited, 1)
	f.iniWriters = append(f.iniWriters, writers...)
	return writers
}

// Stats return the frame internal stats.
func (f *AnimationSequence) Stats() Stats {
	return f.stat
}

// Sequence builds the lists of writers from each sequence item within
// the frame sequence lists.
func (f *AnimationSequence) Sequence() DeferWriters {
	var writers DeferWriters

	// Collect all writers from each sequence with in the frame.
	for _, seq := range f.sequences {
		// If the sequence has finished its rounds, then skip it.
		if seq.IsDone() {
			continue
		}

		writers = append(writers, seq.Next(f.Stats()))
	}

	return writers
}
