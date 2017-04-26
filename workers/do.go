package workers

import (
	"fmt"

	"github.com/influx6/faux/regos"
)

// Workers defines a package level workers registery.
var Workers = regos.New()

// Works define a map of
type Works map[string]Worker

// Has returns true/false if the tag exists for a built result.
func (r Works) Has(tag string) bool {
	_, ok := r[tag]
	return ok
}

// Get returns the Result for the specified build tag.
func (r Works) Get(tag string) Worker {
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
					err = fmt.Errorf("Sumex failed to build Stream[%s] with Tag[%s]: [%s]", do.Name, do.Tag, ex)
				}
			}()

			pb := Workers.NewBuild(do.Name, do.Use).(Worker)
			res[do.Tag] = pb
		}()
	}

	if err != nil {
		return nil, err
	}

	return res, nil
}
