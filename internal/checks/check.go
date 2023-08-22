// check.go
package checks

// Entry point to checks, defers to platform, rulebased, script and other from here

import (
	"fmt"
	"log"

	"github.com/centine/ticli/internal/config"
)

// Define your struct
type CheckResult struct {
	CheckName string
	Status    bool
	Notes     string
}

type CheckFunc func() (bool, string)

type Check struct {
	Name     string
	Platform string
	Fn       CheckFunc
}

type Checker interface {
	DoSetup() error
	DoCleanup() error
	DoCheck() ([]CheckResult, error)
}

// Entry point to performing checks
func PerformChecks(ctx config.TicliContext) {
	checkers := []Checker{
		NewPlatformChecker(ctx),
		newRuleBasedChecker(ctx),
		newScriptBasedChecker(ctx),
		// You can add more checkers here
	}

	results, err := runChecker(ctx, checkers)
	if err != nil {
		log.Fatalf("Error performing checks: %s\n", err)
	}
	fmt.Printf("Platform %s\nResult %+v\n", ctx.SelectedPlatform, results)

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
