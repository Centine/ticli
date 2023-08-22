// check_darwin.go
//go:build darwin
// +build darwin

package checks

import (
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

var checks = []Check{}

func (c *DarwinChecker) DoSetup() error {
	// Deliberate no-op for now.
	// TODO: Download darwin setup script
	return nil
}

func (c *DarwinChecker) DoCleanup() error {
	// TODO: implement clean-up
	return nil
}

func (c *DarwinChecker) DoCheck() ([]CheckResult, error) {
	results := make([]CheckResult, 0, len(checks))

	for _, check := range checks {
		pass, detail := check.Fn()

		results = append(results, CheckResult{
			CheckName: check.Name,
			Status:    pass,
			Notes:     detail,
		})
	}

	return results, nil
}
