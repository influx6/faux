package js_test

import (
	"testing"

	"github.com/influx6/faux/js"
	"github.com/influx6/flux"
)

var session = js.New(nil, true, false)

func TestJSPkgBundler(t *testing.T) {

	js, jsmap, err := session.BuildPkg("github.com/influx6/reactors/builders/base", "base")

	if err != nil {
		flux.FatalFailed(t, "Error build gopherjs dir: %s", err)
	}

	if jsmap.Len() > 50 {
		flux.LogPassed(t, "Successfully built js.map package: %d", jsmap.Len())
	}

	if js.Len() > 50 {
		flux.LogPassed(t, "Successfully built js package: %d", js.Len())
	}

}

func TestJSDirBundler(t *testing.T) {
	js, jsmap, err := session.BuildDir("./js", "github.com/influx6/reactors/builders/base", "base")

	if err != nil {
		flux.FatalFailed(t, "Error build gopherjs dir: %s", err)
	}

	if jsmap.Len() > 50 {
		flux.LogPassed(t, "Successfully built js.map package: %d", jsmap.Len())
	}

	if js.Len() > 50 {
		flux.LogPassed(t, "Successfully built js package: %d", js.Len())
	}
}
