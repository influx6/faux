package exec_test

import (
	"bytes"
	"context"
	"os"
	"testing"
	"time"

	"github.com/influx6/faux/exec"
	"github.com/influx6/faux/metrics"
	"github.com/influx6/faux/tests"
)

func TestLsCommand(t *testing.T) {
	var outs, errs bytes.Buffer
	lsCmd := exec.New(exec.Command("ls ./.."), exec.Sync(), exec.Output(&outs), exec.Err(&errs))
	ctx, cn := context.WithTimeout(context.Background(), 20*time.Second)
	defer cn()

	if err := lsCmd.Exec(ctx, metrics.New()); err != nil {
		tests.Info("Output: %+q", outs.Bytes())
		tests.Info("Errs: %+q", errs.Bytes())
		tests.Failed("Should have succcesfully executed command: %+q", err)
	}
	tests.Passed("Should have succcesfully executed command")
	tests.Info("Output: %+q", outs.Bytes())
}

func TestCommandBinaryWithWget(t *testing.T) {
	var outs, errs bytes.Buffer
	lsCmd := exec.New(exec.Command("wget -nv -O - https://get.docker.com/ | sh"), exec.Async(), exec.Output(&outs), exec.Err(&errs))

	if err := lsCmd.Exec(context.Background(), metrics.New()); err != nil {
		tests.Info("Errs: %+q", errs.Bytes())
		tests.Info("Output: %+q", outs.Bytes())
		tests.Failed("Should have succcesfully executed command: %+q", err)
	}
	tests.Passed("Should have succcesfully executed command")
	tests.Info("Output: %+q", outs.Bytes())
}

func TestCommandWithWget(t *testing.T) {
	var outs, errs bytes.Buffer
	lsCmd := exec.New(exec.Commands("wget", "-nv", "-O", "-", "https://get.docker.com/", "|", "sh"), exec.Async(), exec.Output(&outs), exec.Err(&errs))

	if err := lsCmd.Exec(context.Background(), metrics.New()); err != nil {
		tests.Info("Errs: %+q", errs.Bytes())
		tests.Info("Output: %+q", outs.Bytes())
		tests.Failed("Should have succcesfully executed command: %+q", err)
	}
	tests.Passed("Should have succcesfully executed command")
	tests.Info("Output: %+q", outs.Bytes())
}

func TestCommandWithIf(t *testing.T) {
	var outs, errs bytes.Buffer
	lsCmd := exec.New(exec.Command("if ! type budo; then DEBIAN_FRONTEND=noninteractive echo 'bob'; fi"), exec.Sync(), exec.Output(&outs), exec.Err(&errs))
	ctx, cn := context.WithTimeout(context.Background(), 20*time.Second)
	defer cn()

	if err := lsCmd.Exec(ctx, metrics.New()); err != nil {
		tests.Info("Output: %+q", outs.Bytes())
		tests.Info("Errs: %+q", errs.Bytes())
		tests.Failed("Should have succcesfully executed command: %+q", err)
	}
	tests.Passed("Should have succcesfully executed command")
	tests.Info("Output: %+q", outs.Bytes())
}

func TestCommandWithIfSudo(t *testing.T) {
	var outs, errs bytes.Buffer
	lsCmd := exec.New(exec.Command("if ! type sudo; then apt-get update && DEBIAN_FRONTEND=noninteractive eho 'bob'; fi"), exec.Sync(), exec.Output(&outs), exec.Err(&errs))
	ctx, cn := context.WithTimeout(context.Background(), 20*time.Second)
	defer cn()

	if err := lsCmd.Exec(ctx, metrics.New()); err != nil {
		tests.Info("Output: %+q", outs.Bytes())
		tests.Info("Errs: %+q", errs.Bytes())
		tests.Failed("Should have succcesfully executed command: %+q", err)
	}
	tests.Passed("Should have succcesfully executed command")
	tests.Info("Output: %+q", outs.Bytes())
}

func TestCatCommand(t *testing.T) {
	var outs, errs bytes.Buffer
	lsCmd := exec.New(exec.Command("cat /etc/hosts"), exec.Sync(), exec.Output(&outs), exec.Err(&errs), exec.Binary("/bin/bash", "-c"))
	ctx, cn := context.WithTimeout(context.Background(), 50*time.Second)
	defer cn()

	if err := lsCmd.Exec(ctx, metrics.New()); err != nil {
		tests.Info("Output: %+q", outs.Bytes())
		tests.Info("Errs: %+q", errs.Bytes())
		tests.Failed("Should have succcesfully executed command: %+q", err)
	}
	tests.Passed("Should have succcesfully executed command")

	if outs.Len() == 0 {
		tests.Failed("Should have reeived contents from command")
	}
	tests.Passed("Should have reeived contents from command")
}

func TestWgetCommand(t *testing.T) {
	defer os.Remove("index.html")

	var outs, errs bytes.Buffer
	lsCmd := exec.New(exec.Command("wget www.google.com"), exec.Sync(), exec.Output(&outs), exec.Err(&errs), exec.Binary("/bin/bash", "-c"))
	ctx, cn := context.WithTimeout(context.Background(), 50*time.Second)
	defer cn()

	if err := lsCmd.Exec(ctx, metrics.New()); err != nil {
		tests.Info("Output: %+q", outs.Bytes())
		tests.Info("Errs: %+q", errs.Bytes())
		tests.Failed("Should have succcesfully executed command: %+q", err)
	}
	tests.Passed("Should have succcesfully executed command")
	tests.Info("Output: %+q", outs.Bytes())
}
