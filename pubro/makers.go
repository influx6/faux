package pubro

import (
	"errors"
	"fmt"

	"github.com/influx6/faux/pub"
	"github.com/influx6/faux/sumex"
)

// Do provides a instruction struct for building a publisher from pub.
type Do struct {
	Tag  string      // a custom tag to use for result.
	Name string      // Name of registered publisher builder.
	Use  interface{} //argument to use for Inject
}

//==============================================================================

// Publisher defines a list of Do instruction for building publishers.
type Publisher []Do

// Publishers returns a new instance of Publishers.
func Publishers() *Publisher {
	return new(Publisher)
}

// MustAdd adds a new Do instruction into the list. If it fails to add it panics.
func (d *Publisher) MustAdd(do Do) {
	if err := d.Add(do); err != nil {
		panic(err)
	}
}

// Add adds a new Do instruction into the list.
func (d *Publisher) Add(do Do) error {
	if do.Tag == "" {
		return errors.New("Tag can not be empty")
	}

	if do.Name == "" {
		return errors.New("Name can not be empty")
	}

	*d = append(*d, do)
	return nil
}

//==============================================================================

// PBResult represents the built publisher from the make command.
type PBResult struct {
	ID  string
	By  string
	Pub pub.Publisher
}

// PBResults define a map of
type PBResults map[string]PBResult

// Has returns true/false if the tag exists for a built result.
func (r PBResults) Has(tag string) bool {
	_, ok := r[tag]
	return ok
}

// Get returns the Result for the specified build tag.
func (r PBResults) Get(tag string) PBResult {
	return r[tag]
}

// Pub returns the publisher from a specified result using its tag.
func (r PBResults) Pub(tag string) pub.Publisher {
	rs := r.Get(tag)
	return rs.Pub
}

//==============================================================================

// Make builds the Do instruction using the Make builder.
func (d *Publisher) Make() (PBResults, error) {
	var err error
	res := make(PBResults)

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

			pb := Pub(do.Name, do.Use)

			res[do.Tag] = PBResult{
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

//==============================================================================

// Streamer defines a list of Do instruction for building publishers.
type Streamer []Do

// Streamers returns a new instance of Publishers.
func Streamers() *Streamer {
	return new(Streamer)
}

// MustAdd adds a new Do instruction into the list. If it fails to add it panics.
func (d *Streamer) MustAdd(do Do) {
	if err := d.Add(do); err != nil {
		panic(err)
	}
}

// Add adds a new Do instruction into the list.
func (d *Streamer) Add(do Do) error {
	if do.Tag == "" {
		return errors.New("Tag can not be empty")
	}

	if do.Name == "" {
		return errors.New("Name can not be empty")
	}

	*d = append(*d, do)
	return nil
}

//==============================================================================

// SMResult represents the built publisher from the make command.
type SMResult struct {
	ID     string
	By     string
	Stream sumex.Streams
}

// SMResults define a map of
type SMResults map[string]SMResult

// Has returns true/false if the tag exists for a built result.
func (r SMResults) Has(tag string) bool {
	_, ok := r[tag]
	return ok
}

// Get returns the Result for the specified build tag.
func (r SMResults) Get(tag string) SMResult {
	return r[tag]
}

// SMResults returns the SMResults from a specified result using its tag.
func (r SMResults) SMResults(tag string) sumex.Streams {
	rs := r.Get(tag)
	return rs.Stream
}

//==============================================================================

// Make builds the Do instruction using the Make builder.
func (d *Streamer) Make() (SMResults, error) {
	var err error
	res := make(SMResults)

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

			sm := Stream(do.Name, do.Use)

			res[do.Tag] = SMResult{
				ID:     sm.UUID(),
				By:     do.Name,
				Stream: sm,
			}

		}()
	}

	if err != nil {
		return nil, err
	}

	return res, nil
}

//==============================================================================
