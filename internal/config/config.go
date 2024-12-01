// internal/config/config.go
package config

import (
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

type Config struct {
	DefaultModel    string          `json:"default_model" mapstructure:"default_model"`
	AvailableModels []string        `json:"available_models" mapstructure:"available_models"` // New field
	OpenAIKey       string          `json:"openai_api_key" mapstructure:"openai_api_key"`
	OpenAIURL       string          `json:"openai_api_url" mapstructure:"openai_api_url"`
	AnthropicKey    string          `json:"anthropic_api_key" mapstructure:"anthropic_api_key"`
	GeminiKey       string          `json:"gemini_api_key" mapstructure:"gemini_api_key"`
	GroqKey         string          `json:"groq_api_key" mapstructure:"groq_api_key"`
	OllamaHost      string          `json:"ollama_host" mapstructure:"ollama_host"`
	Commands        []CustomCommand `json:"commands" mapstructure:"commands"`
}

type CustomCommand struct {
	Name         string                 `mapstructure:"command"`
	Description  string                 `mapstructure:"description"`
	Example      string                 `mapstructure:"example"`
	Prompt       string                 `mapstructure:"prompt"`
	Variables    map[string]interface{} `mapstructure:"variables"`
	RequireStdin bool                   `mapstructure:"require_stdin"`
}

func Load() (*Config, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	// Set up viper
	v := viper.New()

	// Global config
	globalConfigPath := filepath.Join(home, ".config", "shell-ask")
	v.AddConfigPath(globalConfigPath)

	// Local config
	v.AddConfigPath(".")

	v.SetConfigName("config")
	v.SetConfigType("json")

	// Environment variables
	v.AutomaticEnv()
	v.SetEnvPrefix("SHELL_ASK")

	// Read config
	var config Config
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, err
		}
	}

	if err := v.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}

func GetConfigDir() string {
	configDir, err := os.UserConfigDir()
	if err != nil {
		// Fallback to home directory if user config dir is not available
		homeDir, err := os.UserHomeDir()
		if err != nil {
			// Last resort: use current directory
			return ".shell-ask"
		}
		return filepath.Join(homeDir, ".shell-ask")
	}
	return filepath.Join(configDir, "shell-ask")
}
