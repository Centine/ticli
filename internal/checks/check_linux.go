// check_linux.go
//go:build linux
// +build linux

package checks

import (
	"os/exec"

	"github.com/centine/ticli/internal/config"
)

type LinuxChecker struct {
	ctx config.TicliContext
}

func NewPlatformChecker(ctx config.TicliContext) Checker {
	return &LinuxChecker{
		ctx: ctx,
	}
}

var checks_platform = []Check{
	{CheckName: "Bash installed", Fn: checkBash},
}

func (c *LinuxChecker) DoSetup(cfg config.ConfigType) error {
	// Deliberate no-op for now.
	// TODO: Download linux setup script
	return nil
}

func (c *GenericChecker) DoCheck() ([]CheckResult, error) {
	return generic_check_iterator(checks_platform, c, c.ctx)
}

func (c *LinuxChecker) DoCleanup(cfg config.ConfigType) error {
	// TODO: implement clean-up
	return nil
}

func (c *LinuxChecker) CheckerName() string {
	return "PlatformChecker"
}

func checkBash(ctx config.TicliContext) (CheckStatus, string, error) {
	path, err := exec.LookPath("bash")
	if err != nil {
		return StatusFail, "No bash shell located", nil
	}
	return StatusSuccess, "Bash shell located at " + path, nil
}
