package buffer_test

import (
	"testing"

	"github.com/influx6/faux/pools/buffer"
)

// BenchmarkBytesPool benchmarks speed and memory allocation using NewBytesPool.
func BenchmarkBytesPool(b *testing.B) {
	b.StopTimer()
	b.ReportAllocs()

	by := buffer.NewBytesPool(1024, 10)
	item := []byte("c")

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		writer := by.Get(i * 100)
		writer.Write(item)
		by.Put(writer)
	}

	b.StopTimer()
}
