package fs

import (
	"context"
	"errors"
)

// Listing defines a type which returns the directory listings recieved from the
// passed it gets.
type Listing struct{}

// Do implements the sumex.Handler interface, and performs the requests for
// retrieving the dir address recieved.
func (d Listing) Do(ctx context.Context, fail error, dirPath interface{}) (interface{}, error) {
	dirAddr, ok := dirPath.(string)
	if !ok {
		return errors.New("Invalid data type, expected string")
	}

}
