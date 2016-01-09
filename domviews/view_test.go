package domviews

import (
	"fmt"
	"strings"
	"testing"

	"github.com/influx6/faux/domtrees"
	"github.com/influx6/faux/domtrees/attrs"
	"github.com/influx6/faux/domtrees/elems"
)

var success = "\u2713"
var failed = "\u2717"

var treeRenderlen = 272

type videoList []map[string]string

func (v videoList) Render(m ...string) domtrees.Markup {
	dom := elems.Div()
	for _, data := range v {
		dom.Augment(elems.Video(
			attrs.Src(data["src"]),
			elems.Text(data["name"]),
		))
	}
	return dom
}

func TestView(t *testing.T) {
	videos := NewView(videoList([]map[string]string{
		map[string]string{
			"src":  "https://youtube.com/xF5R32YF4",
			"name": "Joyride Lewis!",
		},
		map[string]string{
			"src":  "https://youtube.com/dox32YF4",
			"name": "Wonderlust Bombs!",
		},
	}))

	bo := videos.RenderHTML()

	if len(bo) != treeRenderlen {
		fatalFailed(t, "Rendered result with invalid length, expected %d but got %d -> \n %s", treeRenderlen, len(bo), bo)
	}

	logPassed(t, "Rendered result accurated with length %d", treeRenderlen)
}

type item string

func (i item) Render(m ...string) domtrees.Markup {
	return elems.Span(elems.Text(fmt.Sprintf("+ %s", i)))
}

func TestSequenceView(t *testing.T) {
	items := SequenceView(SequenceMeta{Tag: "div"}, item("Book"), item("Funch"), item("Fudder"))

	out := string(items.RenderHTML())

	if !strings.Contains(out, "+ Book") {
		t.Errorf("\t%s\tShould contain %q inside rendered output", failed, "+ Book")
	}

	if !strings.Contains(out, "+ Book") {
		t.Errorf("\t%s\tShould contain %q inside rendered output", failed, "+ Funch")
	}

	if !strings.Contains(out, "+ Book") {
		t.Errorf("\t%s\tShould contain %q inside rendered output", failed, "+ Fudder")
	}

	t.Logf("\t%s\tShould contain %q inside rendered output", success, []string{"+ Book", "+ Funch", "+ Fudder"})
}
