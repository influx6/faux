package linux

import (
	"fmt"

	"github.com/influx6/faux/context"
	"github.com/influx6/faux/exec"
	"github.com/influx6/faux/metrics"
)

// pkg constant types
const (
	InstallAction PackageAction = iota
	RemoveAction
	PurgeAction
)

// PackageAction defines a int type to represent a package action for a package installer.
type PackageAction int

// String returns the name of the action.
func (ap PackageAction) String() string {
	switch ap {
	case InstallAction:
		return "install"
	case RemoveAction:
		return "remove"
	case PurgeAction:
		return "purge"
	}

	return "unknown"
}

// UpdateApt runs necessary commands to install `sudo` package on ubuntu/devian systems
func UpdateApt(ctx context.Context, m metrics.Metrics) (int, error) {
	if err := exec.New(exec.Command("if ! type sudo; then exit 1; fi")).Exec(ctx, m); err != nil {
		return exec.New(exec.Async(), exec.Command("sudo apt-get -y update")).ExecWithExitCode(ctx, m)
	}
	return exec.New(exec.Async(), exec.Command("apt-get -y update")).ExecWithExitCode(ctx, m)
}

// InstallSudo runs necessary commands to install `sudo` package on ubuntu/devian systems
func InstallSudo(ctx context.Context, m metrics.Metrics, upstart bool) (int, error) {
	if err := exec.New(exec.Command("if ! type sudo; then exit 1; fi")).Exec(ctx, m); err != nil {
		return DebianPackageInstall("sudo", InstallAction, upstart).ExecWithExitCode(ctx, m)
	}
	return 0, nil
}

// DebianPackageInstall returns a exec.Command that is executed to install/remove a giving ubuntu package.
func DebianPackageInstall(pkgName string, action PackageAction, upstartbased bool, cmds ...exec.CommanderOption) *exec.Commander {
	var command string

	if action == PurgeAction {
		action = RemoveAction
	}

	switch upstartbased {
	case true:
		command = fmt.Sprintf("DEBIAN_FRONTEND=noninteractive sudo -E apt-get %+s -y -o Dpkg::Options::=\"--force-confnew\" %s", action, pkgName)
	case false:
		command = fmt.Sprintf("DEBIAN_FRONTEND=noninteractive sudo -E apt-get %+s -y %s", action, pkgName)
	}

	return exec.ApplyImmediate(exec.New(exec.Command(command), exec.Async()), cmds...)
}
