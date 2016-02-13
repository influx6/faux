package vfx

import "sync"

// DeferWriterList defines a slice of DeferWriters slices.
type DeferWriterList []DeferWriters

// DeferWriterCache provides a concrete DeferWriter cache which catches specific
// writers by using the frame as the key.
type DeferWriterCache struct {
	wl sync.RWMutex
	w  map[Frame]DeferWriterList
}

// NewDeferWriterCache returns a new WriterCache implementing structure.
func NewDeferWriterCache() *DeferWriterCache {
	wc := DeferWriterCache{w: make(map[Frame]DeferWriterList)}
	return &wc
}

// Store stores the giving set of writers for a specific iteration step of an
// animation. These allows using this writers to produce reversal type effects.
func (d *DeferWriterCache) Store(frame Frame, rs int, dws ...DeferWriter) {
	// Since we start from zeroth index, remove one from step to
	// attain correct index.
	var writers DeferWriterList

	if !d.has(frame) {
		writers = make(DeferWriterList, frame.Stats().TotalIterations())
	} else {
		writers = d.get(frame)
	}

	if rs >= len(writers) {
		return
	}

	var writeList DeferWriters

	d.wl.RLock()
	writeList = writers[rs]
	d.wl.RUnlock()

	d.wl.Lock()
	defer d.wl.Unlock()

	writeList = append(writeList, dws...)
	writers[rs] = writeList
	d.w[frame] = writers
}

// Writers returns the giving writers lists for a specific iteration step
// keyed by a given frames frame.
func (d *DeferWriterCache) Writers(frame Frame, rs int) DeferWriters {
	if !d.has(frame) {
		return nil
	}

	if rs < 0 {
		return nil
	}

	writers := d.get(frame)
	writersLen := len(writers)

	if rs >= writersLen {
		rs = writersLen - 1
	}

	// fmt.Printf("Writers at index %d at len %d\n", rs, writersLen)

	var writeList DeferWriters

	d.wl.RLock()
	writeList = writers[rs]
	d.wl.RUnlock()

	return writeList
}

// ClearIteration clears all writers indexed cached pertaining to a specific
// stat at a specific interation step count.
func (d *DeferWriterCache) ClearIteration(frame Frame, rs int) {
	if !d.has(frame) {
		return
	}

	writersLists := d.get(frame)

	totalWriters := len(writersLists)

	// If out real step is over the bar then this is off base.
	if rs >= totalWriters {
		return
	}

	d.wl.Lock()
	defer d.wl.Unlock()
	writersLists[rs] = nil
}

// Clear clears all writers indexed cached pertaining to a specific stat.
func (d *DeferWriterCache) Clear(frame Frame) {
	d.remove(frame)
}

// has returns true/false if the stat is used as a key within the cache.
func (d *DeferWriterCache) has(frame Frame) bool {
	d.wl.RLock()
	defer d.wl.RUnlock()
	_, ok := d.w[frame]
	return ok
}

// get returns the DeferWriter lists keyed by the frame.
func (d *DeferWriterCache) get(frame Frame) DeferWriterList {
	var dw DeferWriterList

	d.wl.RLock()
	defer d.wl.RUnlock()
	dw = d.w[frame]

	return dw
}

// remove deletes the frame keyed item from the cache.
func (d *DeferWriterCache) remove(frame Frame) {
	d.wl.Lock()
	defer d.wl.Unlock()
	delete(d.w, frame)
}
