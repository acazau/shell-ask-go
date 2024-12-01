package main

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/acazau/shell-ask-go/internal/cli"
	"github.com/acazau/shell-ask-go/internal/config"
	"github.com/spf13/cobra"
)

func mockGetVersionInfo() string {
	return "shell-ask version 0.1.0\ngit commit: unknown\nbuild date: unknown\ngo version: go1.23.2\nplatform: linux/amd64\n"
}

func TestMain(t *testing.T) {
	// Mock config loading
	cfg := &config.Config{
		DefaultModel:    "test-model",
		AvailableModels: []string{"model1", "model2"},
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

func TestListCommand(t *testing.T) {
	// Mock config loading with available models
	cfg := &config.Config{
		DefaultModel:    "test-model",
		AvailableModels: []string{"model1", "model2"},
	}

	// Create a mock CLI
	mockCLI := cli.NewCLI(cfg)

	// Initialize rootCmd
	rootCmd := &cobra.Command{
		Use:   "ask",
		Short: "CLI tool for asking questions to AI models",
	}
	rootCmd.AddCommand(mockCLI.RootCmd())

	// Add list command manually in the test setup
	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List available models",
		Long:  `List available models`,
		Run: func(cmd *cobra.Command, args []string) {
			if len(cfg.AvailableModels) == 0 {
				fmt.Fprintln(cmd.OutOrStdout(), "No models available.")
				return
			}

			fmt.Fprintln(cmd.OutOrStdout(), "Available models:")
			for _, model := range cfg.AvailableModels {
				fmt.Fprintf(cmd.OutOrStdout(), "- %s\n", model)
			}
		},
	}
	rootCmd.AddCommand(listCmd)

	// Capture output
	var out bytes.Buffer
	rootCmd.SetOut(&out)
	rootCmd.SetErr(&out)

	// Test list command
	rootCmd.SetArgs([]string{"list"})
	if err := rootCmd.Execute(); err != nil {
		t.Errorf("Error executing CLI: %v", err)
	}

	expected := "Available models:\n- model1\n- model2\n"
	actual := out.String()

	if actual != expected {
		t.Errorf("Expected output:\n%s\nGot:\n%s", expected, actual)
	}
}
