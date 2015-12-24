package tasks_test

import (
	"encoding/json"
	"testing"

	"github.com/ardanlabs/kit/tests"
	"github.com/influx6/faux/pubro"
	"github.com/influx6/faux/tasks"
)

// TestTaskNames validates that the giving tasks are properly registered with
// the giving names.
func TestTaskNames(t *testing.T) {
	hasTask("tasks/watchCommand", t)
	hasTask("tasks/static", t)
	hasTask("tasks/markdown2Templates", t)
	hasTask("tasks/goBinary", t)
	hasTask("tasks/jsClient", t)
}

// TestTaskBuild validates we can build all provided tasks using the provided
// arguments.
func TestTaskBuild(t *testing.T) {
	tests.ResetLog()
	defer tests.DisplayLog()

	makeTask("tasks/watchCommand", tasks.WatchCommandDirective{
		Dir:      "./",
		Commands: []string{"echo 'word'"},
	}, t)

	makeTask("tasks/markdown2Templates", tasks.MarkdownDirective{
		InDir:    "./",
		OutDir:   "./mk",
		Ext:      ".tmpl",
		Sanitize: true,
	}, t)

	makeTask("tasks/goBinary", tasks.GoBinaryDirective{
		Package: "github.com/influx6/flux",
		Name:    "relay",
		OutDir:  "/tmp",
	}, t)

	makeTask("tasks/jsClient", tasks.JSClientDirective{
		Package: "github.com/influx6/haiku",
	}, t)

	makeTask("tasks/static", tasks.StaticDirective{
		InDir:       "./",
		OutDir:      "./mk",
		PackageName: "mk",
		FileName:    "moku.go",
	}, t)
}

// TestBadArgumentTaskBuild validates the failure of building a publisher with
// inaccurate arguments.
func TestBadArgumentTaskBuild(t *testing.T) {
	tests.ResetLog()
	defer tests.DisplayLog()

	makeFailedTask("tasks/watchCommand", tasks.WatchCommandDirective{}, t)

	makeFailedTask("tasks/static", tasks.StaticDirective{}, t)

	makeFailedTask("tasks/markdown2Templates", tasks.MarkdownDirective{}, t)

	makeFailedTask("tasks/tasks/goBinary", tasks.GoBinaryDirective{}, t)

	makeFailedTask("tasks/jsClient", tasks.JSClientDirective{}, t)
}

// hasTask validates the giving task exists and prints the appropriate Pass/Fail
// message.
func hasTask(pubName string, t *testing.T) {
	if !pubro.Has(pubName) {
		t.Fatalf("\t%s\tShould have publisher[%s] registered.", tests.Failed, pubName)
	} else {
		t.Logf("\t%s\tShould have publisher[%s] registered.", tests.Success, pubName)
	}
}

// makeTask validates the giving argument can satisify the building of a publisher,
// printing the appropriate Pass/Fail message.
func makeTask(pubName string, args interface{}, t *testing.T) {
	defer func() {
		if err := recover(); err != nil {
			t.Fatalf("\t%s\tShould have successfully built publisher[%s]: %s.\nArgument{%q}", tests.Failed, pubName, err, query(args))
			return
		}
		t.Logf("\t%s\tShould have successfully built publisher[%s].", tests.Success, pubName)
	}()

	tm := pubro.New(pubName, args)
	tm.Close()
}

// makeFailedTask validates the giving argument can not satisify the building
// of a publisher, printing the appropriate Pass/Fail message.
func makeFailedTask(pubName string, args interface{}, t *testing.T) {
	defer func() {
		if err := recover(); err == nil {
			t.Fatalf("\t%s\tShould have failed to build publisher[%s].\nArgument{%q}", tests.Failed, pubName, query(args))
			return
		}
		t.Logf("\t%s\tShould have failed to build publisher[%s].", tests.Success, pubName)
	}()

	tm := pubro.New(pubName, args)
	tm.Close()
}

// query provides a string version of the value.
func query(value interface{}) string {
	json, err := json.Marshal(value)
	if err != nil {
		return ""
	}

	return string(json)
}
