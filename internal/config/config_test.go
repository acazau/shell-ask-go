// internal/config/config_test.go
package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	// Create temporary config directory
	tempDir := t.TempDir()
	configDir := filepath.Join(tempDir, ".config", "shell-ask")
	err := os.MkdirAll(configDir, 0755)
	if err != nil {
		t.Fatal(err)
	}

	// Create test config file
	testConfig := Config{
		DefaultModel: "gpt-4",
		OpenAIKey:    "test-key",
		Commands: []CustomCommand{
			{
				Name:        "test",
				Description: "Test command",
				Prompt:      "Test prompt",
			},
		},
	}
	configData, err := json.Marshal(testConfig)
	if err != nil {
		t.Fatal(err)
	}

	err = os.WriteFile(filepath.Join(configDir, "config.json"), configData, 0644)
	if err != nil {
		t.Fatal(err)
	}

	// Set environment variable for test
	os.Setenv("HOME", tempDir)

	// Load config
	config, err := Load()
	if err != nil {
		t.Fatal(err)
	}

	// Debug print loaded config
	t.Logf("Loaded config: %+v", config)

	// Verify config
	if config.DefaultModel != testConfig.DefaultModel {
		t.Errorf("expected default model %s, got %s", testConfig.DefaultModel, config.DefaultModel)
	}
}
