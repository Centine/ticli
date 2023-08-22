package checks

import (
	"errors"

	"github.com/centine/ticli/internal/config"
)

// Scriptbased checks go in here.
type ScriptBasedChecker struct {
	ctx config.TicliContext
}

func newScriptBasedChecker(ctx config.TicliContext) Checker {
	return &ScriptBasedChecker{
		ctx: ctx,
	}
}

func (c *ScriptBasedChecker) DoSetup() error {
	// Deliberate no-op for now.
	return nil
}

func (c *ScriptBasedChecker) DoCleanup() error {
	return nil
}

func (gc *ScriptBasedChecker) DoCheck() ([]CheckResult, error) {
	return nil, errors.New("not implemented")
}
