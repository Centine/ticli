package checks

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"log"
	"os/exec"

	"github.com/centine/ticli/internal/config"
	"github.com/centine/ticli/internal/utility"
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

func (c *ScriptBasedChecker) CheckerName() string {
	return "Script"
}

func (c *ScriptBasedChecker) DoSetup() error {
	utility.DownloadAndVerify(&c.ctx)

	return nil
}

func (c *ScriptBasedChecker) DoCleanup() error {
	return nil
}

func (gc *ScriptBasedChecker) DoCheck() ([]CheckResult, error) {
	// FIXME: Hardcoded path for darwin
	scriptPath := gc.ctx.TicliDir + "/scriptbundles/darwin_checks.sh"
	records, err := executeScript(scriptPath)
	if err != nil {
		panic(err)
	}

	var results []CheckResult
	for _, record := range records {
		if len(record) != 3 {
			log.Printf("Expected 3 columns in CSV output, got %d. Record %v\n", len(record), record) // todo: better error handling
		}

		checkName := record[0]
		notes := record[2]
		status, err := ParseCheckStatus(record[1])
		if err != nil {
			status = StatusUnknown
			notes = fmt.Sprintf("Error parsing status: %v. Original message: %v", err, record[2])
		}
		cr := newCheckResult(checkName, status, notes, withSource(gc.CheckerName()))
		results = append(results, cr)

	}
	return results, nil
}

func executeScript(scriptPath string) ([][]string, error) {
	cmd := exec.Command("/bin/bash", scriptPath)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, err
	}

	reader := csv.NewReader(bytes.NewReader(output))
	reader.Comma = ';'
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	return records, nil
}
