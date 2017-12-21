package linux

import (
	"bytes"

	"context"

	"github.com/influx6/faux/exec"
	"github.com/influx6/faux/metrics"
	"github.com/influx6/faux/os/osinfo"
)

// Info retrieves the OSRelease details related to the operating system.
// Specifically useful for debian/linux systems.
func Info(ctx context.Context, m metrics.Metrics) (*osinfo.Info, error) {
	if data, err := useETC(ctx, m); err == nil {
		return osinfo.NewInfo(data)
	}

	data, err := useUsrLib(ctx, m)
	if err != nil {
		return nil, err
	}

	return osinfo.NewInfo(data)
}

func useETC(ctx context.Context, m metrics.Metrics) ([]byte, error) {
	var outs bytes.Buffer

	lsCmd := exec.New(exec.Command("cat /etc/os-release"), exec.Sync(), exec.Output(&outs))
	if err := lsCmd.Exec(ctx, m); err != nil {
		return nil, err
	}

	return outs.Bytes(), nil
}

func useUsrLib(ctx context.Context, m metrics.Metrics) ([]byte, error) {
	var outs bytes.Buffer

	lsCmd := exec.New(exec.Command("cat /usr/lib/os-release"), exec.Sync(), exec.Output(&outs))
	if err := lsCmd.Exec(ctx, m); err != nil {
		return nil, err
	}

	return outs.Bytes(), nil
}
