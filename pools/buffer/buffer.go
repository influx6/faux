package buffer

import (
	"bytes"
	"io"
	"sync"
)

// BytesPool exists to contain multiple RangePool that lies within giving distance range.
// It creates a internal array of BytesPool which are distanced between each other by
// provided distance. Whenever giving call to get a bytes.Buffer for a giving size is
// within existing pool distances, it calls that RangePool responsible for that size and
// retrieves giving bytes.Buffer from that pool. If no range as such exists, it creates
// a new RangePool for the size + BytesPool.Distance set an instantiation, then retrieves
// a bytes.Buffer from that.
type BytesPool struct {
	distance int
	pl       sync.RWMutex
	pool     []*RangePool
}

// NewBytesPool returns a new instance of a BytesPool with size distance used for new pools
// and creates as many as the initialAmount of RangePools internally to service those size
// requests.
func NewBytesPool(distance int, initialAmount int) *BytesPool {
	var initials []*RangePool

	for i := 1; i <= initialAmount; i++ {
		sizeDist := distance * i
		initials = append(initials, NewRangePool(sizeDist))
	}

	return &BytesPool{
		distance: distance,
		pool:     initials,
	}
}

// Put returns the bytes.Buffer by using the bu.Cap when greater than or equal to BytePool.distance,
// it either finds a suitable RangePool to keep this bytes.Buffer or else creates a new RangePool to cater
// for giving size.
func (bp *BytesPool) Put(bu *bytes.Buffer) {
	if bu.Cap() < bp.distance {
		return
	}

	size := bu.Cap()

	bp.pl.RLock()
	for _, pool := range bp.pool {
		if pool.Max < size {
			continue
		}

		bp.pl.RUnlock()
		pool.Put(bu)
		return
	}
	bp.pl.RUnlock()

	// We dont have any pool within size range, so create new RangePool suited for this size.
	newDistance := size + bp.distance
	newPool := NewRangePool(newDistance)
	newPool.Put(bu)

	bp.pl.Lock()
	bp.pool = append(bp.pool, newPool)
	bp.pl.Unlock()
}

// Get returns a new or existing bytes.Buffer from it's internal size RangePool.
// It gets a RangePool or creates one if non exists for the size + it's distance value
// then gets a bytes.Buffer from that RangePool.
func (bp *BytesPool) Get(size int) *bytes.Buffer {
	bp.pl.RLock()

	// loop through RangePool till we find the distance where size is no more
	// greater, which means that pool will be suitable as the size provider for
	// this size need.
	for _, pool := range bp.pool {
		if pool.Max < size {
			continue
		}

		bp.pl.RUnlock()
		return pool.Get()
	}
	bp.pl.RUnlock()

	// We dont have any pool within size range, so create new RangePool suited for this size.
	newDistance := size + bp.distance
	newPool := NewRangePool(newDistance)

	bp.pl.Lock()
	bp.pool = append(bp.pool, newPool)
	bp.pl.Unlock()

	return newPool.Get()
}

// RangePool exists to build on the internal sync.Pool. It builds on the idea of a pool of size ranges
// where a giving range pool exists, where all data within those range pool are of a giving size.
// This allows us get a buffer instance which is guaranteed to have an existing array of
// the giving size range always.
type RangePool struct {
	Max  int
	pool *sync.Pool
}

// NewRangePool returns a new RangePool instance.
func NewRangePool(max int) *RangePool {
	var rp RangePool
	rp.Max = max
	rp.pool = new(sync.Pool)
	rp.pool.New = func() interface{} {
		return bytes.NewBuffer(make([]byte, 0, max))
	}
	return &rp
}

// Get returns an existing bytes.Buffer that whoes internal byte slice is within
// underline size range of bytes. Ensure to call RangePool.Put when done with bytes.Buffer
// instance.
func (rp *RangePool) Get() (bu *bytes.Buffer) {
	return rp.pool.Get().(*bytes.Buffer)
}

// Put adds giving bytes.Buffer back into internal range pool.
func (rp *RangePool) Put(bu *bytes.Buffer) {
	bu.Reset()
	rp.pool.Put(bu)
}

// GuardedBuffer wraps a giving io.Writer with a mutex guard.
type GuardedBuffer struct {
	mu sync.Mutex
	w  *bytes.Buffer
}

// NewGuardedBuffer returns a new instance of a GuardedWriter.
func NewGuardedBuffer(b *bytes.Buffer) *GuardedBuffer {
	return &GuardedBuffer{w: b}
}

// Write passes provided data to underline writer guarded by mutex.
func (gw *GuardedBuffer) Do(action func(*bytes.Buffer)) {
	if action == nil {
		return
	}

	gw.mu.Lock()
	defer gw.mu.Unlock()
	action(gw.w)
}

// GuardedWriter wraps a giving io.Writer with a mutex guard.
type GuardedWriter struct {
	mu sync.Mutex
	w  io.Writer
}

// NewGuardedWriter returns a new instance of a GuardedWriter.
func NewGuardedWriter(w io.Writer) *GuardedWriter {
	return &GuardedWriter{w: w}
}

// Write passes provided data to underline writer guarded by mutex.
func (gw *GuardedWriter) Write(d []byte) (int, error) {
	gw.mu.Lock()
	defer gw.mu.Unlock()
	return gw.w.Write(d)
}
