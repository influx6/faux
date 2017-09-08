package context_test

import (
	"testing"

	"github.com/influx6/faux/context"
	"github.com/influx6/faux/tests"
)

func TestValueBag(t *testing.T) {
	bag := context.ValueBag()
	bag.Set("k", "z")
	bag.Set("m", "z")
	bag.Set("k", "a")
	bag.Set("m", "v")

	if val, _ := bag.GetString("m"); val != "v" {
		tests.Failed("Should match expected value")
	}
	tests.Passed("Should match expected value")

	if val, _ := bag.GetString("k"); val != "a" {
		tests.Failed("Should match expected value")
	}
	tests.Passed("Should match expected value")
}
