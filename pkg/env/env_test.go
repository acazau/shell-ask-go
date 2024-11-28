package env

import (
	"os"
	"testing"
)

func TestGetConfigDir(t *testing.T) {
	expected := os.Getenv("HOME") + "/.config/shell-ask"
	actual, err := GetConfigDir()
	if err != nil {
		t.Errorf("GetConfigDir failed: %v", err)
	}
	if actual != expected {
		t.Errorf("GetConfigDir mismatch: got %s, want %s", actual, expected)
	}
}

func TestGetCacheDir(t *testing.T) {
	expected := os.Getenv("HOME") + "/.cache/shell-ask"
	actual, err := GetCacheDir()
	if err != nil {
		t.Errorf("GetCacheDir failed: %v", err)
	}
	if actual != expected {
		t.Errorf("GetCacheDir mismatch: got %s, want %s", actual, expected)
	}
}

