package falto_test

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/ardanlabs/kit/tests"
	"github.com/influx6/faux/falto"
)

var context = "tests"

type appd struct{ wg sync.WaitGroup }

// Write simulates a file write operation with possibly completion overhead of a
// few seconds to allow proper simulation of writing capability of ringbuffers.
func (a *appd) Write(data []byte) (int, error) {
	defer a.wg.Done()
	fmt.Printf("Data: %s", data)
	time.Sleep(1 * time.Second)
	return len(data), nil
}

// TestRingBuffer validates the behaviour and stability of the ring buffer as a
// log-only low-latency write mechanism.
func TestRingBuffer(t *testing.T) {
	tests.ResetLog()
	defer tests.DisplayLog()

	var wg sync.WaitGroup

	wg.Add(4)
	ring := falto.NewRingBuffer(falto.RingDirective{
		Capacity:    100,
		SleepResort: 2 * time.Second,
	}, &appd{wg})

	ring.Write(context, []byte("Alex"))
	ring.Write(context, []byte("back"))
	ring.Write(context, []byte("jug"))
	ring.Write(context, []byte("slug"))

	wg.Wait()
}

// BenchmarkRingBuffer validates the speed capability of the ringbuffer in handling
// massive writes with low-latency capability.
func BenchmarkRingBuffer(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()

	ring := falto.NewRingBuffer(falto.RingDirective{
		Capacity:    100,
		SleepResort: 2 * time.Second,
	}, &appd{})

	for i := 0; i < b.N; i++ {
		ring.Write(context, []byte(fmt.Sprintf("%d", i)))
	}
}

func init() {
	tests.Init("")
}
