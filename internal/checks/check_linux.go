// check_linux.go
//go:build linux
// +build linux

package checks

type LinuxChecker struct{}

func (c LinuxChecker) DoCheck() []CheckResult {
	// Linux-specific implementation here
}

func NewChecker() Checker {
	return &LinuxChecker{}
}
