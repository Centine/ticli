// check_darwin.go
//go:build darwin
// +build darwin

package checks

import (
	"os/exec"

	"github.com/centine/ticli/internal/config"
)

type DarwinChecker struct {
	ctx config.TicliContext
}

func NewPlatformChecker(ctx config.TicliContext) Checker {
	return &DarwinChecker{
		ctx: ctx,
	}
}

var checks_platform = []Check{
	{CheckName: "Bash installed", Fn: checkBash},
}

func (c *DarwinChecker) DoSetup() error {
	// Deliberate no-op for now.
	return nil
}

func (c *DarwinChecker) DoCleanup() error {
	// Deliberate no-op for now.
	return nil
}

func (c *DarwinChecker) DoCheck() ([]CheckResult, error) {
	return generic_check_iterator(checks_platform, c, c.ctx)
}

func (c *DarwinChecker) CheckerName() string {
	return "PlatformChecker"
}

func checkBash(ctx config.TicliContext) (CheckStatus, string, error) {
	path, err := exec.LookPath("bash")
	if err != nil {
		return StatusFail, "No bash shell located", nil
	}
	return StatusSuccess, "Bash shell located at " + path, nil
}
