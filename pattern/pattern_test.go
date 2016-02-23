package pattern_test

import (
	"testing"

	"github.com/influx6/faux/pattern"
)

func TestPriority(t *testing.T) {
	if pattern.CheckPriority(`/admin/id`) != 0 {
		t.Fatal(`/admin/id priority is not 0`)
	}
	if pattern.CheckPriority(`/admin/:id`) != 2 {
		t.Fatal(`/admin/:id priority is not 2`)
	}

	if pattern.CheckPriority(`/admin/{id:[\d+]}/:name`) != 1 {
		t.Fatal(`/admin/:id priority is not 1`)
	}
}

func TestEndless(t *testing.T) {
	if !pattern.IsEndless(`/admin/id/*`) {
		t.Fatal(`/admin/id/* is not endless`)
	}
	if pattern.IsEndless(`/admin/id*`) {
		t.Fatal(`/admin/id* is not endless`)
	}
}

func TestPicker(t *testing.T) {
	if pattern.HasPick(`id`) {
		t.Fatal(`/admin/id has no picker`)
	}
	if !pattern.HasPick(`:id`) {
		t.Fatal(`/admin/:id has picker`)
	}
}

func TestSpecialChecker(t *testing.T) {
	if !pattern.HasParam(`/admin/{id:[\d+]}`) {
		t.Fatal(`/admin/{id:[\d+]} is special`)
	}
	if pattern.HasParam(`/admin/id`) {
		t.Fatal(`/admin/id is not special`)
	}
	if !pattern.HasParam(`/admin/:id`) {
		t.Fatal(`/admin/:id is special`)
	}
	if pattern.HasKeyParam(`/admin/:id`) {
		t.Fatal(`/admin/:id is special`)
	}
}

func TestClassicMuxPicker(t *testing.T) {
	cpattern := `/name/:id`
	r := pattern.New(cpattern)

	if r == nil {
		t.Fatalf("invalid array: %+s", r)
	}

	param, state := r.Validate(`/name/12`)
	if !state {
		t.Fatalf("incorrect pattern: %+s %t", param, state)
	}

}

func TestClassicMux(t *testing.T) {
	cpattern := `/name/{id:[\d+]}/`

	r := pattern.New(cpattern)

	if r == nil {
		t.Fatalf("invalid array: %+s", r)
	}

	param, state := r.Validate(`/name/12/d`)

	if state {
		t.Fatalf("incorrect pattern: %+s %t", param, state)
	}

}
