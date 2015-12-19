package tasks_test

import (
	"testing"

	"github.com/ardanlabs/kit/tests"
	"github.com/influx6/faux/pubro"
	"github.com/influx6/faux/tasks"
)

func TestWatchCommand(t *testing.T) {
	tests.ResetLog()
	defer tests.DisplayLog()

	pubName := "tasks/watchCommand"

	if !pubro.Has(pubName) {
		t.Fatalf("\t%s\tShould have publisher[%s] registered", tests.Failed, pubName)
	} else {
		t.Logf("\t%s\tShould have publisher[%s] registered", tests.Success, pubName)
	}

	m, _ := pubro.Get(pubName)
	t.Logf("Meta: %s", m)

	wc := pubro.New(pubName, tasks.WatchCommandConfig{
		Dir:      "./",
		Commands: []string{"echo 'word'"},
	})

	wc.Close()
}
