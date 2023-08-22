package checks

import (
	"fmt"
	"os/exec"
	"regexp"

	"github.com/Masterminds/semver"
	"github.com/centine/ticli/internal/config"
	"golang.org/x/exp/slices"
)

// Rule-based checks go in here.

type RulebasedChecker struct {
	ctx config.TicliContext
}

func newRuleBasedChecker(ctx config.TicliContext) Checker {
	return &RulebasedChecker{
		ctx: ctx,
	}
}

func (c *RulebasedChecker) DoSetup() error {
	// Deliberate no-op for now.
	return nil
}

func (c *RulebasedChecker) DoCleanup() error {
	return nil
}

func (gc *RulebasedChecker) DoCheck() ([]CheckResult, error) {
	platformOverrideFlag := "Linux" // FIXME
	checkResults := make([]CheckResult, 0)

	for _, dynCheck := range gc.ctx.Config.RuleBasedChecks {
		var result *CheckResult
		var err error

		switch dynCheck.Type {
		case "command_compare_versions":
			fmt.Printf("Running command_compare_versions: %s\n", dynCheck.Command)
			result, err = gc.performCommandCompareVersions(dynCheck)
		case "command_exists_in_path":
			fmt.Printf("Running command_compare_versions: %s\n", dynCheck.Command)
			result, err = gc.performCommandExistsInPath(dynCheck)

		default:
			fmt.Printf("Unknown dynamic check type: %s\n", dynCheck.Type)
		}

		if err != nil {
			fmt.Printf("Error running: %+v, error: %v\n", dynCheck, err)
		}

		fmt.Printf("Platform %s\nResult %+v\n", platformOverrideFlag, result)

	}
	return checkResults, nil // TODO: return errors
}

func (gc *RulebasedChecker) performCommandExistsInPath(check config.RuleBasedCheck) (*CheckResult, error) {
	path, err := exec.LookPath(check.Command)
	if err != nil {
		fmt.Printf("%s not found in PATH\n", check.Command)
	} else {
		fmt.Printf("%s found at path %s\n", check.Command, path)
	}

	_, err = exec.Command("sh", "-c", check.Command).Output()
	if err != nil {
		fmt.Println("Error executing command:", err)
		return nil, err
	}
	return &CheckResult{
		CheckName: "performCommandExistsInPath: " + check.Command,
		Status:    true,
		Notes:     "found in path " + path,
	}, nil
}

func (gc *RulebasedChecker) performCommandCompareVersions(check config.RuleBasedCheck) (*CheckResult, error) {
	cmdOutput, err := exec.Command("sh", "-c", check.Command).Output()
	if err != nil {
		exitError, ok := err.(*exec.ExitError)
		if ok {
			// The command failed to run; we can get the exit code
			exitCode := exitError.ExitCode()
			if !slices.Contains(check.AllowableExitCodes, exitCode) {
				fmt.Printf("Error executing command, output '%s', error: %v\n", cmdOutput, err)
				return nil, err

			}
		} else {
			fmt.Printf("Error executing command, output '%s', error: %v\n", cmdOutput, err)
			return nil, err
		}
	}

	re := regexp.MustCompile(check.OutputRegex)
	matches := re.FindAllStringSubmatch(string(cmdOutput), -1)

	for _, match := range matches {
		if len(match) < 3 {
			return nil, fmt.Errorf("failed to match regex: %s", re.String())
		}

		tool := match[1]
		version := match[2]
		fmt.Printf("!!!! Match %s %s\n", tool, version)
		for _, v := range gc.ctx.Config.ToolVersions {
			if v.Tool == tool || v.Platform == "*" {
				minVersion, err := semver.NewVersion(v.MinVersion)
				if err != nil {
					return nil, fmt.Errorf("error parsing semantic version: %v", err)
				}

				actualVersion, err := semver.NewVersion(version)
				if err != nil {
					return nil, fmt.Errorf("error parsing semantic version: %v", err)
				}

				if actualVersion.LessThan(minVersion) {
					fmt.Printf("%s version is less than the minimum required version (%s < %s)\n", tool, actualVersion, minVersion)
				} else {
					fmt.Printf("%s version is okay (%s >= %s)\n", tool, actualVersion, minVersion)
				}
			}
		}
	}
	return nil, nil
}
