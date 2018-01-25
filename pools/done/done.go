package done

import (
	"bytes"
	"errors"
	"io"
	"sync"
)

// errors ...
var (
	ErrClosed        = errors.New("writer already closed")
	ErrLimitExceeded = errors.New("writer exceeded available space")
	dWBuffer         = sync.Pool{New: func() interface{} {
		return new(doneWriter)
	}}
)

// DoneFunc defines a function type for calling a close op.
type DoneFunc func(int, io.WriterTo) error

// doneWriter embodies a giving writer which writes into it's content
// into a underline function after giving size has being reached. It
// has a underline bytes.Buffer which it uses as underline storage.
type doneWriter struct {
	max      int
	index    int
	DoneFunc DoneFunc
	src      *DonePool
	buffer   *bytes.Buffer
}

// Close calls the underline Done method for the giving doneWriter,
// it also resets the buffer, and adds it back into the appropriate
// DonePool.
func (bw *doneWriter) Close() error {
	if bw.buffer == nil {
		return ErrClosed
	}

	written := bw.buffer.Len()
	if written == 0 {
		return nil
	}

	buffer := bw.buffer

	var err error
	if bw.DoneFunc != nil {
		err = bw.DoneFunc(written, buffer)
	}

	buffer.Reset()
	bw.src.put(bw.index, buffer)

	bw.buffer = nil
	bw.DoneFunc = nil
	bw.max = 0
	bw.index = -1
	bw.src = nil

	dWBuffer.Put(bw)
	return err
}

// Write writes the provided data into the doneWriter's buffer.
func (bw *doneWriter) Write(d []byte) (int, error) {
	if bw.buffer == nil {
		return 0, ErrClosed
	}

	// if we have reached the max size, then
	// return return error.
	if bw.max == bw.buffer.Len() {
		return 0, ErrLimitExceeded
	}

	// if we are going to exceed max content, then
	// return error.
	if len(d)+bw.buffer.Len() > bw.max {
		return 0, ErrLimitExceeded
	}

	return bw.buffer.Write(d)
}

type rangePool struct {
	max  int
	pool *sync.Pool
}

// DonePool exists to contain multiple RangePool that lies within giving distance range.
// It creates a internal array of DonePool which are distanced between each other by
// provided distance. Whenever giving call to get a bytes.Buffer for a giving size is
// within existing pool distances, it calls that RangePool responsible for that size and
// retrieves giving bytes.Buffer from that pool. If no range as such exists, it creates
// a new RangePool for the size + DonePool.Distance set an instantiation, then retrieves
// a bytes.Buffer from that.
type DonePool struct {
	distance int
	pl       sync.Mutex
	pools    []rangePool
}

// NewDonePool returns a new instance of a DonePool with size distance used for new pools
// and creates as many as the initialAmount of RangePools internally to service those size
// requests.
func NewDonePool(distance int, initialAmount int) *DonePool {
	initials := make([]rangePool, 0)

	for i := 1; i <= initialAmount; i++ {
		sizeDist := distance * i
		initials = append(initials, rangePool{
			max: sizeDist,
			pool: &sync.Pool{
				New: func() interface{} {
					return bytes.NewBuffer(make([]byte, 0, sizeDist))
				},
			},
		})
	}

	return &DonePool{
		distance: distance,
		pools:    initials,
	}
}

// Put returns the bytes.Buffer by using the bu.Cap when greater than or equal to BytePool.distance,
// it either finds a suitable RangePool to keep this bytes.Buffer or else creates a new RangePool to cater
// for giving size.
func (bp *DonePool) put(index int, bu *bytes.Buffer) {
	bp.pl.Lock()
	defer bp.pl.Unlock()

	if index >= len(bp.pools) {
		return
	}

	pool := bp.pools[index]
	pool.pool.Put(bu)
}

// Get returns a new or existing bytes.Buffer from it's internal size RangePool.
// It gets a RangePool or creates one if non exists for the size + it's distance value
// then gets a bytes.Buffer from that RangePool.
func (bp *DonePool) Get(size int, doneFunc DoneFunc) io.WriteCloser {
	bp.pl.Lock()
	defer bp.pl.Unlock()

	doWriter := dWBuffer.Get().(*doneWriter)

	// loop through RangePool till we find the distance where size is no more
	// greater, which means that pool will be suitable as the size provider for
	// this size need.
	for index, pool := range bp.pools {
		if pool.max < size {
			continue
		}

		doWriter.max = size
		doWriter.src = bp
		doWriter.index = index
		doWriter.buffer = pool.pool.Get().(*bytes.Buffer)
		doWriter.DoneFunc = doneFunc

		return doWriter
	}

	// We dont have any pool within size range, so create new RangePool suited for this size.
	newDistance := size + bp.distance
	newPool := rangePool{
		max: newDistance,
		pool: &sync.Pool{
			New: func() interface{} {
				return bytes.NewBuffer(make([]byte, 0, newDistance))
			},
		},
	}

	index := len(bp.pools)
	bp.pools = append(bp.pools, newPool)

	doWriter.max = size
	doWriter.src = bp
	doWriter.index = index
	doWriter.DoneFunc = doneFunc
	doWriter.buffer = newPool.pool.Get().(*bytes.Buffer)
	return doWriter
}
