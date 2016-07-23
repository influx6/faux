package reflection_test

import (
	"fmt"
	"testing"

	"github.com/influx6/faux/reflection"
)

// succeedMark is the Unicode codepoint for a check mark.
const succeedMark = "\u2713"

// failedMark is the Unicode codepoint for an X mark.
const failedMark = "\u2717"


// mosnter provides a basic struct test case type.
type monster struct {
	Name string
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

	newVals := reflection.MakeArgumentsValues(args)
	t.Logf("%+s", newVals)

	if nlen, alen := len(newVals), len(args); nlen != alen {
		t.Fatalf("\t%s\tShould have matching new values lists for arguments: %d %d", failedMark, nlen, alen)
	} else {
		t.Logf("\t%s\tShould have matching new values lists for arguments: %d %d", succeedMark, nlen, alen)
	}

}

func TestMatchFUncArgumentTypeWithValues(t *testing.T){
	f := func(m monster) string {
		return fmt.Sprintf("Monster[%s] is ready!", m.Name)
	}

	var vals []reflect.Value
	vals = append(vals, reflect.ValueOf("strongHold"))

	if err := reflection.MatchFuncArgumentTypeWithValues(f, ); err != nil {
		t.Fatalf("\t%s\tShould have matching new values lists for arguments: %d %d", failedMark, nlen, alen)
	} else {
		t.Logf("\t%s\tShould have matching new values lists for arguments: %d %d", succeedMark, nlen, alen)
	}

}
