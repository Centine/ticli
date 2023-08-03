// check_windows.go
//go:build windows
// +build windows

// internal/windows/windows.go
package checks

import (
	"fmt"
	"os/exec"

	"golang.org/x/sys/windows"
)

type WindowsChecker struct{}

func NewChecker() Checker {
	return &WindowsChecker{}
}

var checks = []Check{
	{
		Name: "Check PowerShell installed",
		Fn:   CheckPSInstalled,
	},
	{
		Name: "Check Windows version",
		Fn:   CheckWindowsVersion,
	},
}

func CheckPSInstalled() (bool, string) {
	fmt.Println("Executing Windows-native specific checks...")

	cmd := exec.Command("powershell", "-Command", "echo 'Testing powershell'")
	if err := cmd.Run(); err != nil {
		return false, "Powershell is not installed or not in the PATH"
	} else {
		return true, "Powershell is installed"
	}
}

func CheckWindowsVersion() (bool, string) {
	maj, min, patch := windows.RtlGetNtVersionNumbers()
	return true, fmt.Sprintf("Windows version is %s.%s.%s", maj, min, patch)
}

func (c WindowsChecker) DoCheck() []CheckResult {
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
