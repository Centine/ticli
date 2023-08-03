// check_darwin.go
//go:build darwin
// +build darwin

package checks

type DarwinChecker struct{}

var checks = []Check{}

func (c DarwinChecker) DoCheck() []CheckResult {
	results := make([]CheckResult, 0, len(checks))

	for _, check := range checks {
		pass, detail := check.Fn()
		status := "Fail"
		if pass {
			status = "Pass"
		}
		results = append(results, CheckResult{
			CheckName: check.Name,
			Status:    status,
			Notes:     detail,
		})
	}

	return results
}

func NewChecker() Checker {
	return &DarwinChecker{}
}
