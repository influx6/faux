package maker_test

import (
	"testing"

	"github.com/ardanlabs/kit/tests"
	"github.com/influx6/faux/maker"
)

func init() {
	tests.Init("")
}

func TestMaker(t *testing.T) {
	m := maker.New(nil)

	if err := m.Register("sub", func(m int) int {
		return m * 2
	}); err != nil {
		t.Fatalf("\t%s\tShould have registered new make command[%s]", tests.Failed, "sub")
	}
	t.Logf("\t%s\tShould have registered new make command[%s]", tests.Success, "sub")

	res, err := m.Create("sub", 3)
	if err != nil {
		t.Fatalf("\t%s\tShould have successfully create with make command[%s]", tests.Failed, "sub")
	}
	t.Logf("\t%s\tShould have successfully create with make command[%s]", tests.Success, "sub")

	val, ok := res.(int)
	if !ok {
		t.Logf("\t\tResult: %+v", val)
		t.Fatalf("\t%s\tShould have successfully received a int as result for command[%s]", tests.Failed, "sub")
	}
	t.Logf("\t%s\tShould have successfully received a int as result for command[%s]", tests.Success, "sub")

	if val != 6 {
		t.Fatalf("\t%s\tShould have successfully received 6 for command[%s] with param 3", tests.Failed, "sub")
	}
	t.Logf("\t%s\tShould have successfully received 6 for command[%s] with param 3", tests.Success, "sub")
}
