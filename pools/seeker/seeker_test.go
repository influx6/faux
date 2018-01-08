package seeker_test

import (
	"bytes"
	"testing"

	"github.com/influx6/faux/pools/seeker"
	"github.com/influx6/faux/tests"
)

func TestBufferedPeeker(t *testing.T) {
	content := []byte("Thunder world, Reckage before the dawn")
	buff := seeker.NewBufferedPeeker(content)

	if buff.Length() != len(content) {
		tests.Failed("Should have same length has content")
	}
	tests.Passed("Should have same length has content")

	buff.Peek(2)
	if buff.Area() != len(content) {
		tests.Failed("Should have same length has content")
	}
	tests.Passed("Should have same length has content")

	next := buff.Next(2)
	if !bytes.Equal(next, content[:2]) {
		tests.Failed("Should match sub elements of same area")
	}
	tests.Passed("Should match sub elements of same area")

	next = buff.Next(100)
	if !bytes.Equal(next, content[2:]) {
		tests.Failed("Should match sub elements of rest of slice")
	}
	tests.Passed("Should match sub elements of rest of slice")

	if len(buff.Peek(2)) != 0 {
		tests.Failed("Should have index way past length of slice")
	}
	tests.Passed("Should have index way past length of slice")

	buff.Reverse(5)
	if len(buff.Peek(2)) == 0 {
		tests.Failed("Should have index way back within slice")
	}
	tests.Passed("Should have index way back within slice")

	buff.Reverse(100)
	next = buff.Next(2)
	if !bytes.Equal(next, content[:2]) {
		tests.Failed("Should match sub elements of same area")
	}
	tests.Passed("Should match sub elements of same area")
}
