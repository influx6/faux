package pub

import (
	"fmt"

	"github.com/influx6/faux/regos"
)

// Pubs defines a builder for creating and registering Nodes creators.
var Pubs = regos.New()

// Works define a map of
type Works map[string]Handler

// Has returns true/false if the tag exists for a built result.
func (r Works) Has(tag string) bool {
	_, ok := r[tag]
	return ok
}

// Get returns the Result for the specified build tag.
func (r Works) Get(tag string) Handler {
	return r[tag]
}

//==============================================================================

// Work defines a list of regos.DO actions.
type Work []regos.Do

// Make builds the Do instruction using the Make builder.
func (d Work) Make() (Works, error) {
	res := make(Works)

	var err error

	for _, do := range d {
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

			pb := Pubs.NewBuild(do.Name, do.Use).(Handler)
			res[do.Tag] = pb
		}()
	}

	if err != nil {
		return nil, err
	}

	return res, nil
}
