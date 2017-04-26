package utils_test

import (
	"testing"
	"time"

	"github.com/influx6/faux/utils"
)

func TestGetDuration(t *testing.T) {
	tm, err := utils.GetDuration("20ns")
	if err != nil {
		t.Fatalf("\t\u2717\t Should have sucessfully generated time.Duration measure: %q", err.Error())
	}
	t.Logf("\t\u2713\t Should have sucessfully generated time.Duration measure")

	if tm != (20 * time.Nanosecond) {
		t.Fatalf("\t\u2717\t Should have generated 20ns, but got: %q", tm)
	}
	t.Logf("\t\u2713\t Should have generated 20ns, but got: %q", tm)

	tm, err = utils.GetDuration("20s")
	if err != nil {
		t.Fatalf("\t\u2717\t Should have sucessfully generated time.Duration measure: %q", err.Error())
	}
	t.Logf("\t\u2713\t Should have sucessfully generated time.Duration measure")

	if tm != (20 * time.Second) {
		t.Fatalf("\t\u2717\t Should have generated 20s, but got: %q", tm)
	}
	t.Logf("\t\u2713\t Should have generated 20s, but got: %q", tm)

	tm, err = utils.GetDuration("20ms")
	if err != nil {
		t.Fatalf("\t\u2717\t Should have sucessfully generated time.Duration measure: %q", err.Error())
	}
	t.Logf("\t\u2713\t Should have sucessfully generated time.Duration measure")

	if tm != (20 * time.Millisecond) {
		t.Fatalf("\t\u2717\t Should have generated 20ms, but got: %q", tm)
	}
	t.Logf("\t\u2713\t Should have generated 20ms, but got: %q", tm)

	tm, err = utils.GetDuration("20m")
	if err != nil {
		t.Fatalf("\t\u2717\t Should have sucessfully generated time.Duration measure: %q", err.Error())
	}
	t.Logf("\t\u2713\t Should have sucessfully generated time.Duration measure")

	if tm != (20 * time.Minute) {
		t.Fatalf("\t\u2717\t Should have generated 20m, but got: %q", tm)
	}
	t.Logf("\t\u2713\t Should have generated 20m, but got: %q", tm)

	tm, err = utils.GetDuration("20h")
	if err != nil {
		t.Fatalf("\t\u2717\t Should have sucessfully generated time.Duration measure: %q", err.Error())
	}
	t.Logf("\t\u2713\t Should have sucessfully generated time.Duration measure")

	if tm != (20 * time.Hour) {
		t.Fatalf("\t\u2717\t Should have generated 20h, but got: %q", tm)
	}
	t.Logf("\t\u2713\t Should have generated 20h, but got: %q", tm)
}
