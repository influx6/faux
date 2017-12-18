// Package errors is created to contain both a series of
// error values and error wrapper types.
package error

import "errors"

// varables of different error values.
var (
	ErrNotFound       = errors.New("not found")
	ErrRecordNotFound = errors.New("record not found")
)
