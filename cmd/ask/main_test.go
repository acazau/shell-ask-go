package main

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/spf13/cobra"
	"github.com/acazau/shell-ask-go/internal/cli"
	"github.com/acazau/shell-ask-go/internal/config"
)

func mockGetVersionInfo() string {
	return "shell-ask version 0.1.0\ngit commit: unknown\nbuild date: unknown\ngo version: go1.23.2\nplatform: linux/amd64\n"
}

func TestMain(t *testing.T) {
	// Mock config loading
	cfg := &config.Config{
		DefaultModel: "test-model",
	}

	// Create a mock CLI
	mockCLI := cli.NewCLI(cfg)

	// Initialize rootCmd
	rootCmd := &cobra.Command{
		Use:   "ask",
		Short: "CLI tool for asking questions to AI models",
	}
	rootCmd.AddCommand(mockCLI.RootCmd())

	// Add version command
	rootCmd.AddCommand(&cobra.Command{
		Use:   "version",
		Short: "Print version information",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Fprint(cmd.OutOrStdout(), mockGetVersionInfo())
		},
	})

	// Capture output
	var out bytes.Buffer
	rootCmd.SetOut(&out)
	rootCmd.SetErr(&out)

	// Test version command
	rootCmd.SetArgs([]string{"version"})
	if err := rootCmd.Execute(); err != nil {
		t.Errorf("Error executing CLI: %v", err)
	}

	expected := mockGetVersionInfo()
	actual := out.String()

	if actual != expected {
		t.Errorf("Expected output:\n%s\nGot:\n%s", expected, actual)
	}
}
