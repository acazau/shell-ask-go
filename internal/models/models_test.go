package models

import (
	"testing"
)

func TestModelInfo(t *testing.T) {
	expected := ModelInfo{
		ID:          "1",
		RealID:      "real-1",
		Name:        "Test Model",
		Description: "This is a test model.",
		Family:      "test-family",
	}

	actual := ModelInfo{
		ID:          "1",
		RealID:      "real-1",
		Name:        "Test Model",
		Description: "This is a test model.",
		Family:      "test-family",
	}

	if actual != expected {
		t.Errorf("ModelInfo mismatch: got %v, want %v", actual, expected)
	}
}

