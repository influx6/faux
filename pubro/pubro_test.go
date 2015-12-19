package pubro_test

import (
	"fmt"
	"strings"
	"sync"
	"testing"

	"github.com/ardanlabs/kit/tests"
	"github.com/influx6/faux/pub"
	"github.com/influx6/faux/pubro"
)

func init() {
	tests.Init("")
}

type monster struct {
	Name string
}

func mon(m monster) pub.Publisher {
	return pub.Simple(func(p pub.Publisher, data interface{}) {
		p.Reply(fmt.Sprintf("Monster[%s] just ate %+s", m.Name, data))
	})
}

func init() {
}

// pubName is name of the register publisher.
var pubName = "pubro/mon"

// TestRegister tests the registering API function for pubro.
func TestRegister(t *testing.T) {
	tests.ResetLog()
	defer tests.DisplayLog()

	pubro.Register(pubro.Meta{
		Name:   pubName,
		Desc:   "mon builds a publisher that tells you what a monster eats",
		Inject: mon,
	})

	if !pubro.Has(pubName) {
		t.Fatalf("\t%s\tShould have a publisher inject with name[%s]", tests.Failed, pubName)
	} else {
		t.Logf("\t%s\tShould have a publisher inject with name[%s]", tests.Success, pubName)
	}

}

// TestGet tests the retrieving API function Get for pubro.
func TestGet(t *testing.T) {
	tests.ResetLog()
	defer tests.DisplayLog()

	meta, err := pubro.Get(pubName)
	if err != nil {
		t.Fatalf("\t%s\tShould have a publisher inject with name[%s]", tests.Failed, pubName)
	} else {
		t.Logf("\t%s\tShould have a publisher inject with name[%s]", tests.Success, pubName)
	}

	pkg := "github.com/influx6/faux/pubro"
	if meta.Package != pkg {
		t.Errorf("\t%s\tShould have a package path[%s]", tests.Failed, pkg)
	} else {
		t.Logf("\t%s\tShould have a package path[%s]", tests.Success, pkg)
	}

}

// TestNew tests the registering API function for pubro.
func TestNew(t *testing.T) {
	tests.ResetLog()
	defer tests.DisplayLog()

	var wg sync.WaitGroup

	pubmon := pubro.New(pubName, monster{Name: "Bob"})

	wg.Add(1)
	pubmon.React(pub.IdiomaticMuxer(func(err error, data interface{}) (interface{}, error) {
		wg.Done()
		if eats, ok := data.(string); ok && !strings.Contains(eats, "Mushroom") {
			t.Fatalf("\t%s\tShould have a 'Mushrooms' in reply: %s", tests.Failed, eats)
		} else {
			t.Logf("\t%s\tShould have a 'Mushrooms' in reply: %s", tests.Success, eats)
		}

		return data, err
	}), true)

	pubmon.Send("Mushrooms")
	wg.Wait()
}
