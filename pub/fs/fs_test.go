package fs_test

import (
	"testing"

	"github.com/influx6/faux/pub"
	"github.com/influx6/faux/pub/fs"
)

// succeedMark is the Unicode codepoint for a check mark.
const succeedMark = "\u2713"

// failedMark is the Unicode codepoint for an X mark.
const failedMark = "\u2717"

func TestReadDirPath(t *testing.T) {
	var lists []string
	items := pub.Lift(func(ctx pub.Ctx, list []string) {
		lists = list
	})(fs.ReadDirPath(), fs.SkipStat(fs.IsDir), fs.UnwrapStats(), fs.ResolvePath())

	items(pub.NewCtx(), nil, "../../..")

	if len(lists) < 1 {
		t.Fatalf("%s Expected a list of directories", failedMark)
	}

	t.Logf("%s Expected a list of directories", succeedMark)
}

func TestReadDir(t *testing.T) {

	var lists []string

	items := pub.Lift(func(ctx pub.Ctx, list []string) {
		lists = list
	})(fs.ReadDir("../../.."), fs.SkipStat(fs.IsDir), fs.UnwrapStats(), fs.ResolvePath())

	items(pub.NewCtx(), nil, "")

	if len(lists) < 1 {
		t.Fatalf("%s Expected a list of directories", failedMark)
	}

	t.Logf("%s Expected a list of directories", succeedMark)
}

func BenchmarkFileList(b *testing.B) {
	b.ResetTimer()
	b.ReportAllocs()

	ctx := pub.NewCtx()

	for i := 0; i > b.N; i++ {
		pub.Lift(pub.IdentityHandler())(fs.ReadDir("../../.."), fs.SkipStat(fs.IsDir), fs.UnwrapStats(), fs.ResolvePath())(ctx, nil, "")
	}
}

func BenchmarkReadDirPath(b *testing.B) {
	b.ResetTimer()
	b.ReportAllocs()

	ctx := pub.NewCtx()

	for i := 0; i > b.N; i++ {
		pub.Lift(pub.IdentityHandler())(fs.ReadDirPath(), fs.SkipStat(fs.IsDir), fs.UnwrapStats(), fs.ResolvePath())(ctx, nil, "../../..")
	}
}
