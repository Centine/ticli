package utility

import (
	"testing"
)

func TestIsPlatformCompatible(t *testing.T) {
	tests := []struct {
		wanted   string
		actual   string
		expected bool
		testCase string
	}{
		{"*", "linux", true, "Wildcard wanted"},
		{"linux", "*", true, "Wildcard actual"},
		{"linux", "linux", true, "Same platform (lowercase)"},
		{"Linux", "linux", true, "Same platform (mixed case)"},
		{"linux", "windows", false, "Different platforms"},
		{"darwin", "Darwin", true, "Same platform (case insensitive)"},
		{"wsl", "WSL", true, "Same platform (uppercase)"},
		{"windows", "linux", false, "Incompatible platforms"},
	}

	for _, test := range tests {
		actualResult := is_platform_compatible(test.wanted, test.actual)
		if actualResult != test.expected {
			t.Errorf("Test %s failed: got %v, want %v", test.testCase, actualResult, test.expected)
		}
	}
}
