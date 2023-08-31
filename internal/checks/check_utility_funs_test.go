package checks

import "testing"

func TestCheckStatusString(t *testing.T) {
	tests := []struct {
		input    CheckStatus
		expected string
	}{
		{StatusSuccess, "SUCCESS"},
		{StatusFail, "FAIL"},
		{StatusWarning, "WARNING"},
		{StatusInconclusive, "INCONCLUSIVE"},
		{StatusSkipped, "SKIPPED"},
		{StatusNotApplicable, "NOT APPLICABLE"},
		{StatusUnknown, "UNKNOWN"},
		{CheckStatus(1000), "UNKNOWN"}, // testing for unknown values
	}

	for _, test := range tests {
		output := test.input.String()
		if output != test.expected {
			t.Errorf("Expected %s but got %s for input %v", test.expected, output, test.input)
		}
	}
}

func TestParseCheckStatus(t *testing.T) {
	tests := []struct {
		input      string
		expected   CheckStatus
		shouldFail bool
	}{
		{"SUCCESS", StatusSuccess, false},
		{"FAIL", StatusFail, false},
		{"WARNING", StatusWarning, false},
		{"INCONCLUSIVE", StatusInconclusive, false},
		{"SKIPPED", StatusSkipped, false},
		{"NOT APPLICABLE", StatusNotApplicable, false},
		{"UNKNOWN", StatusUnknown, false},
		{"INVALID", -1, true}, // testing for an invalid value
	}

	for _, test := range tests {
		output, err := ParseCheckStatus(test.input)
		if test.shouldFail {
			if err == nil {
				t.Errorf("Expected an error for input %s", test.input)
			}
		} else if output != test.expected {
			t.Errorf("Expected %v but got %v for input %s", test.expected, output, test.input)
		}
	}
}
