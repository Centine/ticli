// check.go
package checks

// Define your struct
type CheckResult struct {
	CheckName string
	Status    string
	Notes     string
}

type CheckFunc func() (bool, string)

type Check struct {
	Name     string
	Platform string
	Fn       CheckFunc
}

type Checker interface {
	DoCheck() []CheckResult
}

// Runs platform-specific checks
func PerformChecks() {
	PerformPlatformSpecificChecks()
	// putils.TableFromStructSlice(*pterm.DefaultTable.WithHasHeader(), results).Render()

}

func PerformPlatformSpecificChecks() []CheckResult {
	checker := NewChecker()
	result := checker.DoCheck()
	return result

}
