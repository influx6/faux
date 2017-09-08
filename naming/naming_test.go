package naming_test

import (
	"testing"

	"github.com/influx6/faux/naming"
)

// TestSimpleNamer validates the use of the provided namer to match the
// giving template rule set provided.
func TestSimpleNamer(t *testing.T) {
	namer := naming.GenerateNamer(naming.SimpleNamer{}, "API-%s-%s")

	firstName := namer("Trappa")
	secondName := namer("Honey")

	if firstName == secondName {
		t.Fatalf("Should have successfully generated new unique generation names")
	}
	t.Logf("Should have successfully generated new unique generation names")
}
