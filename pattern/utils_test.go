package pattern_test

import (
	"strings"
	"testing"

	"github.com/influx6/faux/pattern"
)

func TestTrimEndSlashe(t *testing.T) {

	cases := map[string]struct {
		original, expected string
	}{
		"Single slash at the end of the string": struct {
			original, expected string
		}{
			"/users/6/influx6/", "/users/6/influx6",
		},
		"Two slashes at the end of the string": struct {
			original, expected string
		}{
			"/users/6/influx6//", "/users/6/influx6",
		},
	}

	for k, v := range cases {
		if got := pattern.TrimEndSlashe(v.original); !strings.EqualFold(got, v.expected) {
			t.Fatalf(`Test for %s failed .\n. Expected %s.. Got %s`,
				k, v.expected, got)
		}
	}
}

func TestTrimSlashes(t *testing.T) {

	cases := map[string]struct {
		original, expected string
	}{
		"Slash present only at the end of the string": struct {
			original, expected string
		}{
			"users/6/influx6/", "users/6/influx6",
		},
		"Slashes present as a prefix and suffix": struct {
			original, expected string
		}{
			"/users/6/influx6/", "users/6/influx6",
		},
	}

	for k, v := range cases {

		if got := pattern.TrimSlashes(v.original); !strings.EqualFold(got, v.expected) {
			t.Fatalf(`
				Test for %s failed \n..Expected %s. Got %s`,
				k, v.expected, got)
		}
	}
}
