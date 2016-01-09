package domviews

import (
	"fmt"
	"testing"
)

// succeedMark is the Unicode codepoint for a check mark.
const succeedMark = "\u2713"

// failedMark is the Unicode codepoint for an X mark.
const failedMark = "\u2717"

func TestStateEngineAll(t *testing.T) {
	var engine = NewStateEngine()

	// home :=
	home := engine.AddState("home")

	home.UseActivator(func() {
		logPassed(t, "Sucessfully activated Home")
	})

	home.Engine().AddState(".").UseActivator(func() {
		logPassed(t, "Sucessfully activated border")
	})

	home.Engine().AddState("swatch").UseActivator(func() {
		logPassed(t, "Sucessfully activated swatch")
	})

	err := engine.All(".home.swatch")

	if err != nil {
		fatalFailed(t, "Unable to run full state: %s", err)
	}

}

func TestStateEnginePartial(t *testing.T) {
	var engine = NewStateEngine()

	home := engine.AddState("home")

	home.UseActivator(func() {
		fatalFailed(t, "Should not have activated home")
	})

	home.Engine().AddState(".").UseActivator(func() {
		fatalFailed(t, "Should not have activated border")
	})

	home.Engine().AddState("swatch").UseActivator(func() {
		logPassed(t, "Sucessfully activated swatch")
	})

	err := engine.Partial(".home.swatch")

	if err != nil {
		fatalFailed(t, "Unable to run partial state: %s", err)
	}

}

func TestStateEngineDeactivate(t *testing.T) {
	var engine = NewStateEngine()

	home := engine.AddState("home")

	home.UseActivator(func() {
		logPassed(t, "Sucessfully activated home")
	})

	home.Engine().AddState("swatch").UseActivator(func() {
		logPassed(t, "Sucessfully activated swatch")
	}).UseDeactivator(func() {
		logPassed(t, "Sucessfully deactivated swatch")
	})

	err := engine.All(".home.swatch")

	if err != nil {
		fatalFailed(t, "Unable to run full state: %s", err)
	}

	err = engine.All(".home")

	if err != nil {
		fatalFailed(t, "Unable to run deactivate state: %s", err)
	}

}

func TestStateEngineRoot(t *testing.T) {
	var engine = NewStateEngine()

	home := engine.AddState(".")

	home.UseActivator(func() {
		logPassed(t, "Sucessfully activated home")
	})

	home.Engine().AddState(".").UseActivator(func() {
		logPassed(t, "Sucessfully activated swatch")
	}).UseDeactivator(func() {
		logPassed(t, "Sucessfully deactivated swatch")
	})

	err := engine.All(".")

	if err != nil {
		fatalFailed(t, "Unable to run full state: %s", err)
	}
}

func logPassed(t *testing.T, msg string, data ...interface{}) {
	t.Logf("%s %s", fmt.Sprintf(msg, data...), succeedMark)
}

func fatalFailed(t *testing.T, msg string, data ...interface{}) {
	t.Errorf("%s %s", fmt.Sprintf(msg, data...), failedMark)
}
