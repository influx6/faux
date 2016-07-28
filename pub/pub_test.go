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
	pos := pub.Lift(func(r pub.Ctx, number int) int {
		atomic.AddInt64(&count, 1)
		return number * 2
	})()

	err := errors.New("Ful")
	_, e := pos(ctx, err, nil)
	if e != err {
		fatalFailed(t, "Should have recieved err %s but got %s", err, e)
	}
	logPassed(t, "Should have recieved err %s", err)

	res, _ := pos(ctx, nil, 30)
	if res != 60 {
		fatalFailed(t, "Should have returned %d given %d", 60, 30)
	}
	logPassed(t, "Should have returned %d given %d", 60, 30)

	pos(ctx, nil, "Word") // -> This would not be seen. Has it does not match int type.

	res, _ = pos(ctx, nil, 20)
	if res != 40 {
		fatalFailed(t, "Should have returned %d given %d", 40, 20)
	}
	logPassed(t, "Should have returned %d given %d", 40, 20)

	if atomic.LoadInt64(&count) != 2 {
		fatalFailed(t, "Total processed values is not equal, expected %d but got %d", 2, count)
	}
	logPassed(t, "Total processed values was with count %d", count)
}

// BenchmarkNodes benches the performance of using the Node api.
func BenchmarkNodes(b *testing.B) {
	b.ResetTimer()
	b.ReportAllocs()

	ctx := pub.NewCtx()
	read := pub.Lift(func(r pub.Ctx, number int) int {
		return number * 2
	})()

	for i := 0; i < b.N; i++ {
		read(ctx, nil, i)
	}
}

// BenchmarkNoReflect benches the performance of using the Node api.
func BenchmarkNoReflect(b *testing.B) {
	b.ResetTimer()
	b.ReportAllocs()

	ctx := pub.NewCtx()
	read := pub.Lift(func(r pub.Ctx, number interface{}) interface{} {
		return number.(int) * 2
	})()

	for i := 0; i < b.N; i++ {
		read(ctx, nil, i)
	}
}

func logPassed(t *testing.T, msg string, data ...interface{}) {
	t.Logf("%s %s", fmt.Sprintf(msg, data...), succeedMark)
}

func fatalFailed(t *testing.T, msg string, data ...interface{}) {
	t.Fatalf("%s %s", fmt.Sprintf(msg, data...), failedMark)
}
