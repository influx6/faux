package tmplutil_test

import (
	"bytes"
	"testing"

	"github.com/influx6/faux/tests"
	"github.com/influx6/faux/tmplutil"
)

func TestGroup(t *testing.T) {
	tg := tmplutil.New()

	tg.Add("box.tml", []byte(`Alex went to {{.Place}},{{template "area" .}}`))
	tg.Add("area.tml", []byte(`{{define "area"}} then went to {{.Name}}{{end}}`))

	tml, err := tg.New("box.tml", "area.tml")
	if err != nil {
		tests.Failed("Should have succesffully created template: %+q", err)
	}
	tests.Passed("Should have succesffully created template.")

	tml = tml.Lookup("box.tml")
	if tml == nil {
		tests.Failed("Should have successfully retrieved template 'box.tml'")
	}
	tests.Passed("Should have successfully retrieved template 'box.tml'")

	var buf bytes.Buffer
	if err := tml.Execute(&buf, struct {
		Name  string
		Place string
	}{
		Name:  "Rico",
		Place: "London",
	}); err != nil {
		tests.Failed("Should have successfully executed template: %+q", err)
	}
	tests.Passed("Should have successfully executed template")

	if buf.String() != "Alex went to London, then went to Rico" {
		tests.Failed("Should have succesffuly matched expected response but got %+q", buf.String())
	}
	tests.Passed("Should have succesffuly matched expected response.")
}
