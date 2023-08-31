package checks

import (
	"errors"
	"strings"

	"github.com/centine/ticli/internal/config"
)

// Contains utility functions for performing checks, and parsing in and output

// Returns the string representation of the given CheckStatus.
func (s CheckStatus) String() string {
	switch s {
	case StatusSuccess:
		return "SUCCESS"
	case StatusFail:
		return "FAIL"
	case StatusWarning:
		return "WARNING"
	case StatusInconclusive:
		return "INCONCLUSIVE"
	case StatusSkipped:
		return "SKIPPED"
	case StatusNotApplicable:
		return "NOT APPLICABLE"
	case StatusUnknown:
		return "UNKNOWN"
	default:
		return "UNKNOWN"
	}
}

// ParseCheckStatus parses the given string s and returns the appropriate CheckStatus.
func ParseCheckStatus(s string) (CheckStatus, error) {
	switch strings.ToUpper(s) {
	case "SUCCESS":
		return StatusSuccess, nil
	case "FAIL":
		return StatusFail, nil
	case "WARNING":
		return StatusWarning, nil
	case "INCONCLUSIVE":
		return StatusInconclusive, nil
	case "SKIPPED":
		return StatusSkipped, nil
	case "NOT APPLICABLE":
		return StatusNotApplicable, nil
	case "UNKNOWN":
		return StatusUnknown, nil
	default:
		return -1, errors.New("invalid status")
	}
}

func bool2CheckStatus(b bool) CheckStatus {
	if b {
		return StatusSuccess
	}
	return StatusFail
}

func generic_check_iterator(checks []Check, checker Checker, ctx config.TicliContext) ([]CheckResult, error) {
	var results []CheckResult

	for _, check := range checks {
		status, notes, err := check.Fn(ctx)
		if err != nil {
			return nil, err
		}
		result := CheckResult{
			CheckName: check.CheckName,
			Status:    status,
			Notes:     notes,
			Source:    checker.CheckerName(),
		}
		results = append(results, result)
	}
	return results, nil
}

// Trims the given string to the given max length, and replaces all newlines with spaces.
func trimString(s string, maxLength int) string {
	// Replace all newlines
	s = strings.ReplaceAll(s, "\n", "")

	// Trim to max length
	if len(s) > maxLength {
		s = s[:maxLength]
	}

	return s
}

// Experimental and over-engineered way of passing options to NewCheckResult ;)

type checkResultOption func(*CheckResult)

func withSource(source string) checkResultOption {
	return func(cr *CheckResult) {
		cr.Source = source
	}
}

func newCheckResult(name string, status CheckStatus, notes string, opts ...checkResultOption) CheckResult {
	cr := CheckResult{
		CheckName: name,
		Status:    status,
		Notes:     notes,
	}

	for _, opt := range opts {
		opt(&cr)
	}

	return cr
}
