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

	pkg := "testing"
	if meta.Package != pkg {
		t.Errorf("\t%s\tShould have a package path[%s]: %s", tests.Failed, pkg, meta.Package)
	} else {
		t.Logf("\t%s\tShould have a package path[%s]: %s", tests.Success, pkg, meta.Package)
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

// TestDosMake tests we can make a lists of publishers using the Do API for pubro.
func TestDosMake(t *testing.T) {
	tests.ResetLog()
	defer tests.DisplayLog()

	pubs := pubro.NewDos()

	// Add instruction to build a rat monster.
	pubs.MustAdd(pubro.Do{
		Tag:  "rat",
		Name: pubName,
		Use:  monster{Name: "beans"},
	})

	// Add instruction to build a fox monster.
	pubs.MustAdd(pubro.Do{
		Tag:  "fox",
		Name: pubName,
		Use:  monster{Name: "rat"},
	})

	res, err := pubs.Make()
	if err != nil {
		t.Fatalf("\t%s\tShould have receive a map of results from build: %s.", tests.Failed, err)
	} else {
		t.Logf("\t%s\tShould have receive a map of results from build.", tests.Success)
	}

	if !res.Has("fox") {
		t.Fatalf("\t%s\tShould have result with tag[%s].", tests.Failed, "fox")
	} else {
		t.Logf("\t%s\tShould have result with tag[%s].", tests.Success, "fox")
	}

	if !res.Has("rat") {
		t.Fatalf("\t%s\tShould have result with tag[%s].", tests.Failed, "rat")
	} else {
		t.Logf("\t%s\tShould have result with tag[%s].", tests.Success, "rat")
	}

	if _, ok := res.Pub("rat").(pub.Publisher); !ok {
		t.Fatalf("\t%s\tShould have received a pub.Publisher type.", tests.Failed)
	} else {
		t.Logf("\t%s\tShould have received a pub.Publisher type.", tests.Success)
	}

	if _, ok := res.Pub("fox").(pub.Publisher); !ok {
		t.Fatalf("\t%s\tShould have received a pub.Publisher type.", tests.Failed)
	} else {
		t.Logf("\t%s\tShould have received a pub.Publisher type.", tests.Success)
	}
}

// TestBadDosMake validates build failure using wrong argument type.
func TestBadDosMake(t *testing.T) {
	tests.ResetLog()
	defer tests.DisplayLog()

	pubs := pubro.NewDos()

	// Add instruction to build a rat monster.
	pubs.MustAdd(pubro.Do{
		Tag:  "rat",
		Name: pubName,
		Use:  "beans",
	})

	_, err := pubs.Make()
	if err == nil {
		t.Fatalf("\t%s\tShould fail to build pub result: \n%s.", tests.Failed, err)
	} else {
		t.Logf("\t%s\tShould fail to build pub result: \n%s.", tests.Success, err)
	}

}
