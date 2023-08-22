// check_linux.go
//go:build linux
// +build linux

package checks

import "github.com/centine/ticli/internal/config"

type LinuxChecker struct {
	ctx config.TicliContext
}

func NewPlatformChecker(ctx config.TicliContext) Checker {
	return &LinuxChecker{
		ctx: ctx,
	}
}

func (c *LinuxChecker) DoSetup(cfg config.ConfigType) error {
	// Deliberate no-op for now.
	// TODO: Download linux setup script
	return nil
}

func (c *LinuxChecker) DoCheck(cfg config.ConfigType) ([]CheckResult, error) {
	// Linux-specific implementation here
	return results, nil
}

func (c *LinuxChecker) DoCleanup(cfg config.ConfigType) error {
	// TODO: implement clean-up
	return nil
}
