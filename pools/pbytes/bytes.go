package pbytes

import (
	"bytes"
	"sync"
)

//************************************************************************
// BytesPool
//************************************************************************

type rangePool struct {
	max  int
	pool *sync.Pool
}

// BytesPool exists to contain multiple RangePool that lies within giving distance range.
// It creates a internal array of BytesPool which are distanced between each other by
// provided distance. Whenever giving call to get a []byte for a giving size is
// within existing pool distances, it calls that RangePool responsible for that size and
// retrieves giving []byte from that pool. If no range as such exists, it creates
// a new RangePool for the size + BytesPool.Distance set an instantiation, then retrieves
// a []byte from that.
type BytesPool struct {
	distance int
	pl       sync.Mutex
	pools    []rangePool
	indexes  map[int]int
}

// NewBytesPool returns a new instance of a BytesPool with size distance used for new pools
// and creates as many as the initialAmount of RangePools internally to service those size
// requests.
func NewBytesPool(distance int, initialAmount int) *BytesPool {
	initials := make([]rangePool, 0)
	indexes := make(map[int]int)

	for i := 1; i <= initialAmount; i++ {
		sizeDist := distance * i

		indexes[sizeDist] = len(initials)
		initials = append(initials, rangePool{
			max: sizeDist,
			pool: &sync.Pool{
				New: func() interface{} {
					return bytes.NewBuffer(make([]byte, 0, sizeDist))
				},
			},
		})
	}

	return &BytesPool{
		distance: distance,
		pools:    initials,
		indexes:  indexes,
	}
}

// Put returns the []byte by using the capacity of the slice to find its pool.
func (bp *BytesPool) Put(bu *bytes.Buffer) {
	bp.pl.Lock()
	defer bp.pl.Unlock()

	if index, ok := bp.indexes[bu.Cap()]; ok {
		pool := bp.pools[index]
		pool.pool.Put(bu)
	}
}

// Get returns a new or existing []byte from it's internal size RangePool.
// It gets a RangePool or creates one if non exists for the size + it's distance value
// then gets a []byte from that RangePool.
func (bp *BytesPool) Get(size int) *bytes.Buffer {
	bp.pl.Lock()
	defer bp.pl.Unlock()

	// loop through RangePool till we find the distance where size is no more
	// greater, which means that pool will be suitable as the size provider for
	// this size need.
	for _, pool := range bp.pools {
		if pool.max < size {
			continue
		}

		return pool.pool.Get().(*bytes.Buffer)
	}

	// We dont have any pool within size range, so create new RangePool suited for this size.
	newDistance := ((size / bp.distance) + 1) * bp.distance
	newPool := rangePool{
		max: newDistance,
		pool: &sync.Pool{
			New: func() interface{} {
				return bytes.NewBuffer(make([]byte, 0, newDistance))
			},
		},
	}

	bp.indexes[newDistance] = len(bp.pools)
	bp.pools = append(bp.pools, newPool)

	return newPool.pool.Get().(*bytes.Buffer)
}

//************************************************************************
// BytePool
//************************************************************************

// BytePool implements a leaky pool of []byte in the form of a bounded
// channel.
type BytePool struct {
	c chan []byte
	w int
}

// NewBytePool creates a new BytePool bounded to the given maxSize, with new
// byte arrays sized based on width.
func NewBytePool(maxSize int, width int) (bp *BytePool) {
	return &BytePool{
		c: make(chan []byte, maxSize),
		w: width,
	}
}

// Get gets a []byte from the BytePool, or creates a new one if none are
// available in the pool.
func (bp *BytePool) Get() (b []byte) {
	select {
	case b = <-bp.c:
		// reuse existing buffer
	default:
		// create new buffer
		b = make([]byte, bp.w)
	}
	return
}

// Put returns the given Buffer to the BytePool.
func (bp *BytePool) Put(b []byte) {
	select {
	case bp.c <- b:
		// buffer went back into pool
	default:
		// buffer didn't go back into pool, just discard
	}
}

// Width returns the width of the byte arrays in this pool.
func (bp *BytePool) Width() (n int) {
	return bp.w
}