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

			ws.Add(1)

			cache := vfx.NewDeferWriterCache()
			stat := vfx.TimeStat(1*time.Second, "ease-in", false, false)
			stat2 := vfx.TimeStat(2*time.Second, "ease-in", false, false)

			defer cache.Clear(stat)
			defer cache.Clear(stat2)

			go func() {
				defer ws.Done()
				for i := 0; i < stat.TotalIterations(); i++ {
					j := stat.CurrentIteration()
					cache.Store(stat, j, buildNDeferWriters(i+1)...)
					stat.NextIteration(0)
				}
			}()

			ws.Wait()
			ws.Add(2)

			go func() {
				defer ws.Done()
				for i := 0; i < 10; i++ {
					wd := cache.Writers(stat, i)
					if len(wd) != (i + 1) {
						t.Fatalf("\t%s\tShould have writer lists for step: %d of len: %d", tests.Failed, i, i*4)
					}
					t.Logf("\t%s\tShould have writer lists for step: %d of len: %d", tests.Success, i, i*4)
				}
			}()

			go func() {
				defer ws.Done()
				for i := 0; i < stat2.TotalIterations(); i++ {
					j := stat2.CurrentIteration()
					cache.Store(stat2, j, buildNDeferWriters(i+1)...)
					stat2.NextIteration(0)
				}
			}()

			ws.Wait()

			w4 := cache.Writers(stat, 4)
			cache.ClearIteration(stat, 4)
			w42 := cache.Writers(stat, 4)

			if len(w42) >= len(w4) {
				t.Fatalf("\t%s\tShould have cleared writers of 4 iteration", tests.Failed)
			}
			t.Logf("\t%s\tShould have cleared writers of 4 iteration", tests.Success)

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
