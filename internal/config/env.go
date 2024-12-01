// internal/config/env.go
package config

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
)

// EnvConfig stores environment configuration
type EnvConfig struct {
	OpenAIKey    string
	AnthropicKey string
	GoogleAPIKey string
	GroqAPIKey   string
	DefaultModel string
}

// LoadEnv loads environment variables from a .env file if it exists
// and falls back to system environment variables if not found in .env
func LoadEnv() (*EnvConfig, error) {
	config := &EnvConfig{}

	// Try to load from .env in current directory first
	if err := loadEnvFile(".env"); err == nil {
		// Successfully loaded .env from current directory
	} else {
		// If not found in current directory, try home directory
		if homeDir, err := os.UserHomeDir(); err == nil {
			envPath := filepath.Join(homeDir, ".env")
			if err := loadEnvFile(envPath); err == nil {
				// Successfully loaded .env from home directory
			}
		}
	}

	// Load configs with .env values taking precedence over system env vars
	config.OpenAIKey = getEnvWithFallback("OPENAI_API_KEY", "")
	config.AnthropicKey = getEnvWithFallback("ANTHROPIC_API_KEY", "")
	config.GoogleAPIKey = getEnvWithFallback("GOOGLE_API_KEY", "")
	config.GroqAPIKey = getEnvWithFallback("GROQ_API_KEY", "")
	config.DefaultModel = getEnvWithFallback("DEFAULT_MODEL", "gpt-4")

	return config, nil
}

// loadEnvFile reads a .env file and sets environment variables
func loadEnvFile(filepath string) error {
	file, err := os.Open(filepath)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Split on first = sign
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		// Remove quotes if present
		value = strings.Trim(value, `"'`)

		os.Setenv(key, value)
	}

	return scanner.Err()
}

// getEnvWithFallback gets an environment variable with a fallback value
func getEnvWithFallback(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
