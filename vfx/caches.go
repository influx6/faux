package vfx

import "sync"

// DeferWriterList defines a slice of DeferWriters slices.
type DeferWriterList []DeferWriters

// DeferWriterCache provides a concrete DeferWriter cache which catches specific
// writers by using the stats as the key.
type DeferWriterCache struct {
	wl sync.RWMutex
	w  map[Stats]DeferWriterList
}

// NewDeferWriterCache returns a new WriterCache implementing structure.
func NewDeferWriterCache() *DeferWriterCache {
	wc := DeferWriterCache{w: make(map[Stats]DeferWriterList)}
	return &wc
}

// Store stores the giving set of writers for a specific iteration step of an
// animation. These allows using this writers to produce reversal type effects.
func (d *DeferWriterCache) Store(stats Stats, step int, dws ...DeferWriter) {

	// Since we start from zeroth index, remove one from step to
	// attain correct index.
	rs := step - 1

	var writers DeferWriterList

	if !d.has(stats) {
		writers = make(DeferWriterList, stats.TotalIterations())
	} else {
		writers = d.get(stats)
	}

	var writeList DeferWriters

	d.wl.RLock()
	writeList = writers[rs]
	d.wl.RUnlock()

	d.wl.Lock()
	defer d.wl.Unlock()

	writeList = append(writeList, dws...)
	writers[rs] = writeList
	d.w[stats] = writers
}

// Writers returns the giving writers lists for a specific iteration step
// keyed by a given frames stats.
func (d *DeferWriterCache) Writers(stats Stats, step int) DeferWriters {

	if !d.has(stats) {
		return nil
	}

	writers := d.get(stats)

	// Since we start from zeroth index, remove one from step to
	// attain correct index.
	rs := step - 1

	var writeList DeferWriters

	d.wl.RLock()
	writeList = writers[rs]
	d.wl.RUnlock()

	return writeList
}

// ClearIteration clears all writers indexed cached pertaining to a specific
// stat at a specific interation step count.
func (d *DeferWriterCache) ClearIteration(stats Stats, step int) {
	if !d.has(stats) {
		return
	}

	// Since we start from zeroth index, remove one from step to
	// attain correct index.
	realStep := step - 1

	writersLists := d.get(stats)

	totalWriters := len(writersLists)

	// If out real step is over the bar then this is off base.
	if realStep >= totalWriters {
		return
	}

	d.wl.Lock()
	defer d.wl.Unlock()
	writersLists[realStep] = nil
}

// Clear clears all writers indexed cached pertaining to a specific stat.
func (d *DeferWriterCache) Clear(stats Stats) {
	d.remove(stats)
}

// has returns true/false if the stat is used as a key within the cache.
func (d *DeferWriterCache) has(stats Stats) bool {
	d.wl.RLock()
	defer d.wl.RUnlock()
	_, ok := d.w[stats]
	return ok
}

// get returns the DeferWriter lists keyed by the stats.
func (d *DeferWriterCache) get(stats Stats) DeferWriterList {
	var dw DeferWriterList

	d.wl.RLock()
	defer d.wl.RUnlock()
	dw = d.w[stats]

	return dw
}

// remove deletes the stats keyed item from the cache.
func (d *DeferWriterCache) remove(stats Stats) {
	d.wl.Lock()
	defer d.wl.Unlock()
	delete(d.w, stats)
}
