package cli

import (
	"testing"
)

func TestSearch(t *testing.T) {
	tests := []struct {
		query    string
		expected string
	}{
		{"example", "example result"},
		{"other", "no results found"},
	}

	for _, test := range tests {
		result := Search(test.query)
		if result != test.expected {
			t.Errorf("Search(%q) = %q, want %q", test.query, result, test.expected)
		}
	}
}
