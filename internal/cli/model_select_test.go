package cli

import (
	"testing"
)

func TestSelectModel(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"openai", "gpt-4"},
		{"other", "gpt-3.5-turbo"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := SelectModel(tt.input)
			if got != tt.want {
				t.Errorf("SelectModel(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}
