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

}

func hasTask(pubName string, t *testing.T) {
	if !pubro.Has(pubName) {
		t.Fatalf("\t%s\tShould have publisher[%s] registered", tests.Failed, pubName)
	} else {
		t.Logf("\t%s\tShould have publisher[%s] registered", tests.Success, pubName)
	}
}

func makeTask(pubName string, args interface{}, t *testing.T) {
	defer func() {
		if err := recover(); err != nil {
			t.Fatalf("\t%s\tShould have built publisher[%s] using Argument{%q}", tests.Failed, pubName, query(args))
			return
		}
		t.Logf("\t%s\tShould have built publisher[%s] using Argument{%q}", tests.Success, pubName, query(args))
	}()

	tm := pubro.New("tasks/watchCommand", args)
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
