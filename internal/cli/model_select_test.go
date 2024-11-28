package cli

import (
	"testing"
)

func TestSelectModel(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"openai", "gpt-4"},
		{"other", "default-model"},
	}

	for _, test := range tests {
		result := SelectModel(test.input)
		if result != test.expected {
			t.Errorf("SelectModel(%q) = %q, want %q", test.input, result, test.expected)
		}
	}
}
