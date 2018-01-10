package done_test

import (
	"testing"

	"bytes"
	"io"

	"github.com/influx6/faux/pools/done"
	"github.com/influx6/faux/tests"
)

// BenchmarkDonePool benchmarks speed and memory allocation using NewBytesPool.
func BenchmarkDonePool(b *testing.B) {
	b.StopTimer()
	b.ReportAllocs()

	by := done.NewDonePool(1024, 10)
	item := []byte("c")

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		writer := by.Get(i*100, nil)
		writer.Write(item)
		writer.Close()
	}

	b.StopTimer()
}

func TestExceedLimitError(b *testing.T) {

	by := done.NewDonePool(20, 10)
	writer := by.Get(10, nil)
	defer writer.Close()

	_, err := writer.Write([]byte("34332234fsdj"))
	if err == nil {
		tests.Failed("Should have received error for write exceed bountry")
	}
	if err != done.ErrLimitExceeded {
		tests.FailedWithError(err, "Expected ErrLimitExceeded")
	}
	tests.Passed("Should have received error for write exceed bountry")
}

func TestLimitWrite(b *testing.T) {

	by := done.NewDonePool(512, 10)
	writer := by.Get(10, func(w int, src io.WriterTo) error {
		var bu bytes.Buffer
		if _, err := src.WriteTo(&bu); err != nil {
			tests.FailedWithError(err, "Should have received src data without issue")
		}
		tests.Passed("Should have received src data without issue")

		if bu.Len() != w {
			tests.Failed("Received data must match written count")
		}
		tests.Passed("Received data must match written count")
		return nil
	})

	if _, err := writer.Write([]byte("34332234tv")); err != nil {
		tests.FailedWithError(err, "Should have written data without error")
	}
	tests.Passed("Should have written data without error")

	if err := writer.Close(); err != nil {
		tests.FailedWithError(err, "Should have closed without error")
	}
	tests.Passed("Should have closed without error")

}
