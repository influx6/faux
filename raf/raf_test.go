package raf_test

import (
	"testing"
	"time"

	"github.com/influx6/faux/raf"
)

func TestRAF(t *testing.T) {
	var count int

	id := raf.RequestAnimationFrame(func(df float64) {
		count++
	})

	<-time.After(1 * time.Second)
	raf.CancelAnimationFrame(id)

	if count > 62 {
		t.Fatalf("Expected 61 or 62 count for animation calls: %d", count)
	}
	t.Logf("Expected 61 or 62 count for animation calls: %d", count)
}
