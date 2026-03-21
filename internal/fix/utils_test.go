package fix

import "testing"

func TestTranslateExecType(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{ExecTypeNew, "New (Accepted)"},
		{ExecTypePartialFill, "Partially Filled"},
		{ExecTypeFill, "Filled (Complete)"},
		{ExecTypeCanceled, "Canceled"},
		{ExecTypeRejected, "Rejected"},
		{"Unknown", "Unknown"},
	}

	for _, test := range tests {
		result := translateExecType(test.input)
		if result != test.expected {
			t.Errorf("translateExecType(%s) = %s; want %s", test.input, result, test.expected)
		}
	}
}

func TestTranslateOrdStatus(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"0", "New"},
		{"1", "Partially Filled"},
		{"2", "Filled"},
		{"4", "Canceled"},
		{"8", "Rejected"},
		{"99", "99"},
	}

	for _, test := range tests {
		result := translateOrdStatus(test.input)
		if result != test.expected {
			t.Errorf("translateOrdStatus(%s) = %s; want %s", test.input, result, test.expected)
		}
	}
}
