package stacks_test

import (
	"testing"

	"github.com/ardanlabs/kit/tests"
	"github.com/influx6/faux/stacks"
)

func init() {
	tests.Init("")
}

var context = "tests"

// TestStackExtract tests the extraction API for stack traces.
func TestStackExtract(t *testing.T) {
	tests.ResetLog()
	defer tests.DisplayLog()

	var pkg = "github.com/influx6/faux/stacks"
	var pkgRoot = "github.com/influx6/faux"

	trace := make(stacks.Trace, 0)
	stacks.Run(context, false, 0, &trace)

	stack, index := trace.FindMethod("PullTrace")
	if index == -1 {
		t.Fatalf("\t%s\tShould have found trace for Add method", tests.Failed)
	} else {
		t.Logf("\t%s\tShould have found trace for Add method", tests.Success)
	}

	if stack.Package != pkg {
		t.Errorf("\t%s\tShould have package matching %s", tests.Failed, pkg)
	} else {
		t.Logf("\t%s\tShould have package matching %s", tests.Success, pkg)
	}

	if stack.Root != pkgRoot {
		t.Errorf("\t%s\tShould have package root matching %s", tests.Failed, pkgRoot)
	} else {
		t.Logf("\t%s\tShould have package root matching %s", tests.Success, pkgRoot)
	}
	
	// if stack.LineNumber != 106 {
	// 	t.Errorf("\t%s\tShould have method %s matching line number 105: %d", tests.Failed, stack.MethodName, stack.LineNumber)
	// } else {
	// 	t.Logf("\t%s\tShould have method %s matching line number 106: %d", tests.Success, stack.MethodName, stack.LineNumber)
	// }
}
