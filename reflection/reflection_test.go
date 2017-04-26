package reflection_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/influx6/faux/reflection"
)

// succeedMark is the Unicode codepoint for a check mark.
const succeedMark = "\u2713"

// failedMark is the Unicode codepoint for an X mark.
const failedMark = "\u2717"

type bull string

type speaker interface {
	Speak() string
}

// mosnter provides a basic struct test case type.
type monster struct {
	Name  string
	Items []bull
}

// Speak returns the sound the monster makes.
func (m *monster) Speak() string {
	return "Raaaaaaarggg!"
}

func get(t *testing.T, sm speaker) {
	name, embedded, err := reflection.StructAndEmbeddedTypeNames(sm)
	if err != nil {
		t.Fatalf("\t%s\tShould be able to retrieve field names arguments lists: %s", failedMark, err)
	} else {
		t.Logf("\tName: %s\n", name)
		t.Logf("\tFields: %+q\n", embedded)
		t.Logf("\t%s\tShould be able to retrieve function arguments lists", succeedMark)
	}
}

// TestGetArgumentsType validates reflection API GetArgumentsType functions
// results.
func TestGetArgumentsType(t *testing.T) {
	f := func(m monster) string {
		return fmt.Sprintf("Monster[%s] is ready!", m.Name)
	}

	args, err := reflection.GetFuncArgumentsType(f)
	if err != nil {
		t.Fatalf("\t%s\tShould be able to retrieve function arguments lists: %s", failedMark, err)
	} else {
		t.Logf("\t%s\tShould be able to retrieve function arguments lists", succeedMark)
	}

	name, embedded, err := reflection.StructAndEmbeddedTypeNames(monster{Name: "Bob"})
	if err != nil {
		t.Fatalf("\t%s\tShould be able to retrieve field names arguments lists: %s", failedMark, err)
	} else {
		t.Logf("\tName: %s\n", name)
		t.Logf("\tFields: %+q\n", embedded)
		t.Logf("\t%s\tShould be able to retrieve function arguments lists", succeedMark)
	}

	get(t, &monster{Name: "Bob"})

	newVals := reflection.MakeArgumentsValues(args)
	if nlen, alen := len(newVals), len(args); nlen != alen {
		t.Fatalf("\t%s\tShould have matching new values lists for arguments", failedMark)
	}
	t.Logf("\t%s\tShould have matching new values lists for arguments", succeedMark)

	mstring := reflect.TypeOf((*monster)(nil)).Elem()

	if mstring.Kind() != newVals[0].Kind() {
		t.Fatalf("\t%s\tShould be able to match argument kind: %s", failedMark)
	}
	t.Logf("\t%s\tShould be able to match argument kind", succeedMark)

}

func TestMatchFUncArgumentTypeWithValues(t *testing.T) {
	f := func(m monster) string {
		return fmt.Sprintf("Monster[%s] is ready!", m.Name)
	}

	var vals []reflect.Value
	vals = append(vals, reflect.ValueOf(monster{Name: "FireHouse"}))

	if index := reflection.MatchFuncArgumentTypeWithValues(f, vals); index != -1 {
		t.Fatalf("\t%s\tShould have matching new values lists for arguments: %d", failedMark, index)
	}
	t.Logf("\t%s\tShould have matching new values lists for arguments", succeedMark)

}
