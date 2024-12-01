package cli

import (
	"testing"

	"github.com/acazau/shell-ask-go/internal/config"
)

func TestNewCLI(t *testing.T) {
	cfg := &config.Config{
		DefaultModel: "gpt-4",
		OpenAIKey:    "test-key",
		// Add other necessary fields if needed
	}
	cli := NewCLI(cfg)
	if cli == nil {
		t.Errorf("Expected CLI instance, got nil")
	}
}

func TestRootCmd(t *testing.T) {
	cfg := &config.Config{
		DefaultModel: "gpt-4",
		OpenAIKey:    "test-key",
		// Add other necessary fields if needed
	}
	cli := NewCLI(cfg)
	rootCmd := cli.RootCmd()
	if rootCmd == nil {
		t.Errorf("Expected root command, got nil")
	}
}
