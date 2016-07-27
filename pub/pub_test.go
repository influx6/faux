package pub_test

import (
	"errors"
	"fmt"
	"sync/atomic"
	"testing"

	"github.com/influx6/faux/pub"
)

// succeedMark is the Unicode codepoint for a check mark.
const succeedMark = "\u2713"

// failedMark is the Unicode codepoint for an X mark.
const failedMark = "\u2717"

// TestAutoFn  validates the use of reflection with giving types to test use of
// the form in Pub.
func TestAutoFn(t *testing.T) {
	var count int64

	ctx := pub.NewCtx()
	pos := pub.Lift(func(r pub.Ctx, number int) {
		atomic.AddInt64(&count, 1)
	})(nil)

	pos(ctx, nil, errors.New("Ful"))
	pos(ctx, nil, 30)
	pos(ctx, nil, "Word") // -> This would not be seen. Has it does not match int type.
	pos(ctx, nil, 20)

	if atomic.LoadInt64(&count) != 2 {
		fatalFailed(t, "Total processed values is not equal, expected %d but got %d", 3000, count)
	}

	logPassed(t, "Total mouse data was processed with count %d", count)
}

// BenchmarkNodes benches the performance of using the Node api.
func BenchmarkNodes(b *testing.B) {
	b.ResetTimer()
	b.ReportAllocs()

	read := pub.Lift(func(r pub.Ctx, number int) {})(nil)

	ctx := pub.NewCtx()
	for i := 0; i < b.N; i++ {
		read(ctx, nil, i)
	}
}

func logPassed(t *testing.T, msg string, data ...interface{}) {
	t.Logf("%s %s", fmt.Sprintf(msg, data...), succeedMark)
}

func fatalFailed(t *testing.T, msg string, data ...interface{}) {
	t.Errorf("%s %s", fmt.Sprintf(msg, data...), failedMark)
}
