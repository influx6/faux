package domviews

import (
	"errors"
	"sync"
)

// ErrExists is returned when the ViewLists has the tag already asigned
var ErrExists = errors.New("Tag already assigned")

// ViewMap represent a set of Views element tagged by a string key
type ViewMap map[string]int

// ViewList represent a list of Views elements
type ViewList []Views

// ViewLists is a race-free struct providing a map for storing and retrieving Views which is used by *View
type ViewLists struct {
	//a map of keys and index values
	mo   sync.RWMutex
	keys ViewMap
	//the slice of Views
	ro    sync.RWMutex
	lists ViewList
}

// NewViewLists returns a new instance of ViewLists
func NewViewLists() *ViewLists {
	vm := ViewLists{
		keys: make(ViewMap),
	}

	return &vm
}

// Has returns true/false if the tag has a Views
func (vm *ViewLists) Has(tag string) bool {
	vm.mo.RLock()
	_, ok := vm.keys[tag]
	vm.mo.RUnlock()
	return ok
}

// Add adds a Views to the view lists
func (vm *ViewLists) Add(tag string, v Views) error {
	if vm.Has(tag) {
		return ErrExists
	}

	var ind int

	vm.ro.RLock()
	ind = len(vm.lists)
	vm.ro.RUnlock()

	vm.ro.RLock()
	vm.lists = append(vm.lists, v)
	vm.ro.RUnlock()

	vm.mo.Lock()
	vm.keys[tag] = ind
	vm.mo.Unlock()

	return nil
}

// Remove removes the Views if the given tag exists
func (vm *ViewLists) Remove(tag string) Views {
	ind, ok := vm.index(tag)

	if !ok {
		return nil
	}

	v := vm.Get(tag)

	vm.mo.Lock()
	delete(vm.keys, tag)
	vm.mo.Unlock()

	vm.ro.Lock()
	vm.lists = append(vm.lists[0:ind], vm.lists[ind:]...)
	vm.ro.Unlock()

	return v
}

// Get returns the Views with the giving tag if found else nil
func (vm *ViewLists) Get(tag string) Views {
	ind, ok := vm.index(tag)

	if !ok {
		return nil
	}

	var v Views

	vm.ro.RLock()
	v = vm.lists[ind]
	vm.ro.RUnlock()

	return v

}

// Views returns the slice within the viewlists
func (vm *ViewLists) Views() ViewList {
	vm.ro.RLock()
	defer vm.ro.RUnlock()
	return vm.lists
}

// index returns the index of the tag in the keys map if it exists
func (vm *ViewLists) index(tag string) (int, bool) {
	var ind int
	var ok bool

	vm.mo.RLock()
	ind, ok = vm.keys[tag]
	vm.mo.RUnlock()
	return ind, ok
}
