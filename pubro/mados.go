package pubro

import (
	"errors"
	"fmt"

	"github.com/influx6/faux/pub"
)

// Do provides a instruction struct for building a publisher from pub.
type Do struct {
	Tag  string      // a custom tag to use for result.
	Name string      // Name of registered publisher builder.
	Use  interface{} //argument to use for Inject
}

// Result represents the built publisher from the make command.
type Result struct {
	ID  string
	By  string
	Pub pub.Publisher
}

// Results define a map of
type Results map[string]Result

// Has returns true/false if the tag exists for a built result.
func (r Results) Has(tag string) bool {
	_, ok := r[tag]
	return ok
}

// Get returns the Result for the specified build tag.
func (r Results) Get(tag string) Result {
	return r[tag]
}

// Pub returns the publisher from a specified result using its tag.
func (r Results) Pub(tag string) pub.Publisher {
	rs := r.Get(tag)
	return rs.Pub
}

// Dos defines a list of Do instruction for building publishers.
type Dos []Do

// NewDos returns a new Dos lists for building pubs.
func NewDos() *Dos {
	return new(Dos)
}

// MustAdd adds a new Do instruction into the list. If it fails to add it panics.
func (d *Dos) MustAdd(do Do) {
	if err := d.Add(do); err != nil {
		panic(err)
	}
}

// Add adds a new Do instruction into the list.
func (d *Dos) Add(do Do) error {
	if do.Tag == "" {
		return errors.New("Tag can not be empty")
	}

	if do.Name == "" {
		return errors.New("Name can not be empty")
	}

	*d = append(*d, do)
	return nil
}

// Make builds the Do instruction using the Make builder.
func (d *Dos) Make() (Results, error) {
	var err error
	res := make(Results)

	for _, do := range *d {
		if res.Has(do.Tag) {
			err = fmt.Errorf("Build Instruction for %s using reserved tag %s", do.Name, do.Tag)
			break
		}

		if err != nil {
			break
		}

		func() {
			defer func() {
				if ex := recover(); ex != nil {
					err = fmt.Errorf("Pubro failed to build Pub[%s] with Tag[%s]: [%s]", do.Name, do.Tag, ex)
				}
			}()

			pb := New(do.Name, do.Use)

			res[do.Tag] = Result{
				ID:  pb.UUID(),
				By:  do.Name,
				Pub: pb,
			}

		}()
	}

	if err != nil {
		return nil, err
	}

	return res, nil
}
