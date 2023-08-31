// check.go
package checks

// Entry point to checks, defers to platform, rulebased, script and other from here

import (
	"log"

	"github.com/centine/ticli/internal/config"
	"github.com/pterm/pterm"
)

// CheckStatus represents the status of a check.
type CheckStatus int

const (
	// StatusSuccess indicates that the check passed.
	StatusSuccess CheckStatus = iota
	// StatusFail indicates that the check failed.
	StatusFail
	// StatusWarning indicates that the check passed, but with warnings.
	StatusWarning
	// StatusInconclusive indicates that the check was inconclusive.
	StatusInconclusive
	// StatusSkipped indicates that the check was skipped deliberately, either through filtering or configuration.
	StatusSkipped
	// StatusNotApplicable indicates that the check was not applicable to the current platform.
	StatusNotApplicable
	// StatusUnknown indicates that the check status is unknown, hopefully the notes field has more information.
	StatusUnknown
)

// CheckResult represents the result of a check.
type CheckResult struct {
	// Name of check, for identification purposes
	CheckName string
	// Result of check
	Status CheckStatus
	// Optional information to relay to user
	Notes string
	// Source of check, e.g. "script", "platform", "rule"
	Source string
}

type CheckFunc func(ctx config.TicliContext) (CheckStatus, string, error)

type Check struct {
	// Name of check, for identification purposes
	CheckName string
	// Platform compatibility, e.g. "darwin"
	Platform string
	// Function to execute
	Fn CheckFunc
}

type Checker interface {
	DoSetup() error
	DoCleanup() error
	DoCheck() ([]CheckResult, error)
	CheckerName() string
}

// Entry point to performing checks
func PerformChecks(ctx config.TicliContext) {
	checkers := []Checker{
		NewGenericChecker(ctx),
		NewPlatformChecker(ctx),
		newRuleBasedChecker(ctx),
		newScriptBasedChecker(ctx),
		// You can add more checkers here
	}

	pterm.Info.Printf("Performing checks for %s\n", ctx.SelectedPlatform)
	results, err := runChecker(ctx, checkers)
	if err != nil {
		log.Fatalf("Error performing checks: %s\n", err)
	}

	// Convert your data to a 2D slice
	data := [][]string{{"Check Name", "Status", "Source", "Notes"}}
	for _, r := range results {
		data = append(data, []string{r.CheckName, r.Status.String(), r.Source, trimString(r.Notes, 130)})
	}

	// Print data as a table
	pterm.DefaultTable.WithHasHeader().WithData(data).Render()

}

func runChecker(ctx config.TicliContext, checkers []Checker) ([]CheckResult, error) {
	var results []CheckResult

	for _, checker := range checkers {
		defer checker.DoCleanup()

		err := checker.DoSetup()
		if err != nil {
			return nil, err
		}

		result, err := checker.DoCheck()
		if err != nil {
			return nil, err
		}

		results = append(results, result...)
	}

	return results, nil
}
