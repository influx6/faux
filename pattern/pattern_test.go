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

func TestNamePattern(t *testing.T) {
	r := pattern.New(`/name/:id`)
	path := `/name/12`

	param, _, state := r.Validate(path)
	if !state {
		t.Fatalf("Failed: Should have matched path: %s>%#v Params: %#v", r.Pattern(), path, param)
	}

	t.Logf("Passed: Should have matched path: %s->%#v Params: %#v", r.Pattern(), path, param)
}

func TestParam(t *testing.T) {
	r := pattern.New(`:id`)

	param, _, state := r.Validate(`12`)
	if !state {
		t.Fatalf("Failed: Should have matched path: %s->%#v", r.Pattern(), "/12")
	}

	val, ok := param["id"]
	if !ok {
		t.Fatalf("Failed: Should have matched with parameter : %s->%s but %#v", "id", "12", param)
	}

	t.Logf("Passed: Should have matched path: %s->%#v Params: %#v", r.Pattern(), val, param)
}

func TestEndlessPattern(t *testing.T) {
	r := pattern.New(`/github.com/influx6/*`)

	param, rem, state := r.Validate(`/github.com/influx6/examples/views#blob`)
	if !state {
		t.Fatalf("Incorrect pattern: %+s %t", param, state)
	}

	expected := "/examples/views/#blob"
	if rem != expected {
		t.Fatalf("Incorrect remaining path(Expected: %s Found: %s)", expected, rem)
	}
}

func TestAsterick(t *testing.T) {
	r := pattern.New(`*`)

	param, rem, state := r.Validate(`/github.com/influx6/examples/views#blob`)
	if !state {
		t.Fatalf("Incorrect pattern: %+s %t", param, state)
	}

	t.Logf("Correct pattern: %+s %+s", param, rem)
}

func TestHashedWithRemainder(t *testing.T) {
	r := pattern.New(`/colors/*`)

	param, rem, state := r.Validate(`#colors/red`)
	if !state {
		t.Fatalf("incorrect pattern: %+s %t", param, state)
	}

	if rem != "/red" {
		t.Fatalf("incorrect remainer: Expected[%s] Got[%s]", "/red", rem)
	}
}

func TestHashedPatternWithRemainder(t *testing.T) {
	r := pattern.New(`/github.com/influx6/examples/*`)

	param, rem, state := r.Validate(`/github.com/influx6/examples#views`)
	if !state {
		t.Fatalf("incorrect pattern: %+s %t", param, state)
	}

	if rem != "/#views" {
		t.Fatalf("incorrect remainer: Expected[%s] Got[%s]", "#views", rem)
	}
}

func TestHashedPattern(t *testing.T) {
	r := pattern.New(`/github.com/influx6/examples#views`)

	param, _, state := r.Validate(`/github.com/influx6/examples/views`)
	if !state {
		t.Fatalf("incorrect pattern: %+s %t", param, state)
	}
}

func TestRegExpPattern(t *testing.T) {
	r := pattern.New(`/name/{id:[\d+]}/`)

	param, _, state := r.Validate(`/name/12/d`)

	if state {
		t.Fatalf("incorrect pattern: %+s %t", param, state)
	}
}
