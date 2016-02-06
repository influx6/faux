package vfx

import (
	"regexp"
	"strconv"
)

//==============================================================================

// nodigits defines a regexp for matching non-digits.
var nodigits = regexp.MustCompile("[^\\d]+")

// parseFloat parses a string into a float if fails returns the default value 0.
func parseFloat(fl string) float64 {
	fll, _ := strconv.ParseFloat(digitsOnly(fl), 64)
	return fll
}

// parseInt parses a string into a int if fails returns the default value 0.
func parseInt(fl string) int {
	fll, _ := strconv.Atoi(digitsOnly(fl))
	return fll
}

// parseIntBase16 parses a string into a int using base16.
func parseIntBase16(fl string) int {
	fll, _ := strconv.ParseInt(fl, 16, 64)
	return int(fll)
}

// digitsOnly removes all non-digits characters in a string.
func digitsOnly(fl string) string {
	return nodigits.ReplaceAllString(fl, "")
}

//==============================================================================
