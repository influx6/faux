package bag_test

import (
	"testing"

	"github.com/influx6/faux/bag"
	"github.com/influx6/faux/tests"
)

func TestValueBag(t *testing.T) {
	bag := bag.NewValueBag()
	bag = bag.WithValue("k", "z")
	bag = bag.WithValue("m", "z")
	bag = bag.WithValue("k", "a")
	bag.WithValue("m", "v")

	if val, _ := bag.GetString("m"); val != "z" {
		tests.Failed("Should match expected value")
	}
	tests.Passed("Should match expected value")

	if val, _ := bag.GetString("k"); val != "a" {
		tests.Failed("Should match expected value")
	}
	tests.Passed("Should match expected value")
}
