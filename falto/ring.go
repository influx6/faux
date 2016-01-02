package falto

import (
	"io"
	"sync"
	"time"
)

// Entries defines a lists of entry records each providing their individual
// sequence.
type Entries []*Entry

// Entry defines a record section within a buffer set.
type Entry struct {
	Sequence int64
	Buffer   []byte
}

// RingDirective provides a configuration for the RingBuffer and dictates how
// the buffer functions and what consistent level to manage the buffer with.
type RingDirective struct {
	Capacity    int
	SleepResort time.Duration
}

// Sector defines a specific record of a section of the buffer.
type Sector struct {
	Sequence int64
	Index    int
}

// RingBuffer provides a writer using a modified ring buffer to avoid the use of
// locks and mutex. RingBuffer contains a queue channel which is used when the
// the ring is considered full and unable to free space since the writer has not
// moved from its last position either due to a large write, the buffer quickly
// adds the item onto the chan which is always twice the size of the ring buffer
// internal capacity.
// Ring Buffer Operation:
// Falto RingBuffer is a modified ring buffer in that it uses a least necessary
// expansion algorithm to ensure fast and quick insertion of data into a append-only
// file.
type RingBuffer struct {
	directive RingDirective
	ring      Entries
	buffer    Entries
	imQueue   chan []byte // imQueue == immediateQueue
	lrQueue   chan []byte // lrQueue == LastResortQueue
	w         io.Writer
	lSeq      *Sector
	cSeq      *Sector
	wl        sync.Mutex
}

// NewRingBuffer returns a new instance of the Falto RingBuffer.
func NewRingBuffer(rd RingDirective, w io.Writer) *RingBuffer {
	// Creates a buffer at a maximum capacity and allocate a zerod slice to the
	// ring.
	ringBuffer := make(Entries, rd.Capacity)

	rb := RingBuffer{
		directive: rd,
		w:         w,
		buffer:    ringBuffer,
		ring:      ringBuffer[:0],
		imQueue:   make(chan []byte),
		lrQueue:   make(chan []byte, rd.Capacity),
		lSeq:      &Sector{},
		cSeq:      &Sector{},
	}

	return &rb
}

// Write writes into the ring sequence else if full, appends into the sequence
// LastResortQueue.
func (r *RingBuffer) Write(context interface{}, data []byte) {
	r.imQueue <- data
}
