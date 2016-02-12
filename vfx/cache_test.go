package vfx_test

import (
	"sync"
	"testing"
	"time"

	"github.com/ardanlabs/kit/tests"
	"github.com/influx6/faux/vfx"
)

// TestDeferWriterCache validates the operation and write safetiness of DeferWriterCache.
func TestDeferWriterCache(t *testing.T) {
	tests.ResetLog()
	defer tests.DisplayLog()

	t.Logf("Given the need to use a DeferWriterCache")
	{

		t.Logf("\tWhen giving two stat keys")
		{

			var ws sync.WaitGroup

			ws.Add(3)

			cache := vfx.NewDeferWriterCache()
			stat := vfx.TimeStat(1*time.Minute, "ease-in", false, false)
			stat2 := vfx.TimeStat(2*time.Minute, "ease-in", false, false)

			go func() {
				defer ws.Done()
				for i := 0; i < 10; i++ {
					cache.Store(stat, i, buildNDeferWriters(i*4)...)
				}
			}()

			go func() {
				defer ws.Done()
				for i := 0; i < 10; i++ {
					cache.Store(stat2, i, buildNDeferWriters(i*4)...)
				}
			}()

			ws.Wait()

			for i := 0; i < 10; i++ {
				wd := cache.Writers(stat2, i)
				if len(wd) == i*4 {
					t.Fatalf("\t%s\tShould have writer lists of len: %d", tests.Failed, i*4)
				}
				t.Logf("\t%s\tShould have writer lists of len: %d", tests.Success, i*4)
			}
		}
	}

}

//==============================================================================

// wr implements vfx.DeferWriter interface.
type wr struct{}

// Write writes out the writers details.
func (w *wr) Write() {}

func buildNDeferWriters(size int) vfx.DeferWriters {
	var ws vfx.DeferWriters

	for i := 0; i < size; i++ {
		ws = append(ws, &wr{})
	}

	return ws
}

//==============================================================================
