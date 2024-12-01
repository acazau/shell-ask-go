// internal/models/models.go
package models

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strings"
)

type ModelInfo struct {
	ID          string
	RealID      string
	Name        string
	Description string
	Family      string
}

type Models struct {
	Models []ModelInfo
}

var ModelMap = map[string][]ModelInfo{
	"gpt": {
		{ID: "gpt-3.5-turbo"},
		{ID: "gpt-4-turbo"},
		{ID: "gpt-4"},
		{ID: "gpt-4-32k"},
	},
	"claude": {
		{ID: "claude-3-haiku", RealID: "claude-3-haiku-20240307"},
		{ID: "claude-3-sonnet", RealID: "claude-3-sonnet-20240229"},
		{ID: "claude-3-opus", RealID: "claude-3-opus-20240229"},
	},
	"gemini": {
		{ID: "gemini-pro"},
		{ID: "gemini-1.5-pro", RealID: "gemini-1.5-pro-latest"},
		{ID: "gemini-1.5-flash", RealID: "gemini-1.5-flash-latest"},
	},
	"groq": {
		{ID: "groq-llama3", RealID: "llama3-70b-8192"},
		{ID: "groq-mixtral", RealID: "mixtral-8x7b-32768"},
		{ID: "groq-gemma", RealID: "gemma-7b-it"},
	},
	"copilot": {
		{ID: "copilot-chat", Description: "GitHub Copilot Chat"},
	},
	"ollama": {}, // Will be populated dynamically
}

// SelectModel selects a model based on input or returns default
func SelectModel(input string) string {
	// Check if it's an Ollama model format (contains ":")
	if strings.Contains(input, ":") && ValidateOllamaModel(input) {
		return input // Return as-is for valid Ollama models
	}

	// Check direct model ID match
	for _, models := range ModelMap {
		for _, model := range models {
			if model.ID == input {
				if model.RealID != "" {
					return model.RealID
				}
				return model.ID
			}
		}
	}

	// Check provider match
	if models, exists := ModelMap[input]; exists && len(models) > 0 {
		if models[0].RealID != "" {
			return models[0].RealID
		}
		return models[0].ID
	}

	return "gpt-4o-mini" // default model
}

func ValidateOllamaModel(name string) bool {
	// Basic validation for Ollama model format (name:tag)
	parts := strings.Split(name, ":")
	return len(parts) == 2 && len(parts[0]) > 0 && len(parts[1]) > 0
}

func GetAllModels(includeOllama bool) []ModelInfo {
	var allModels []ModelInfo
	for _, models := range ModelMap {
		allModels = append(allModels, models...)
	}

	if includeOllama {
		ollamaModels, err := GetOllamaModels()
		if err == nil {
			allModels = append(allModels, ollamaModels...)
			// Update the ollama section of ModelMap
			ModelMap["ollama"] = ollamaModels
		}
	}

	sort.Slice(allModels, func(i, j int) bool {
		return allModels[i].ID < allModels[j].ID
	})
	return allModels
}

func GetCheapModel(modelID string) string {
	switch {
	case strings.HasPrefix(modelID, "gpt-4"):
		return "gpt-3.5-turbo"
	case strings.HasPrefix(modelID, "claude-3-opus"):
		return "claude-3-haiku"
	case strings.HasPrefix(modelID, "gemini-1.5"):
		return "gemini-pro"
	case strings.HasPrefix(modelID, "groq-llama3"):
		return "groq-gemma"
	case strings.Contains(modelID, ":"):
		// For Ollama models, return the smallest available model
		models, err := GetOllamaModels()
		if err == nil && len(models) > 0 {
			return models[0].ID // Return the first available model
		}
		return modelID
	default:
		return modelID
	}
}

func GetOllamaModels() ([]ModelInfo, error) {
	resp, err := http.Get("http://localhost:11434/api/tags")
	if err != nil {
		return nil, fmt.Errorf("failed to get Ollama models: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var result struct {
		Models []struct {
			Name    string `json:"name"`
			Details string `json:"details"`
		} `json:"models"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	var models []ModelInfo
	for _, m := range result.Models {
		models = append(models, ModelInfo{
			ID:          m.Name,
			Name:        m.Name,
			Description: m.Details,
			Family:      "ollama",
		})
	}

	return models, nil
}

func GetModelInfo(name string) *ModelInfo {
	// First check predefined models
	for _, mi := range ModelMap {
		for _, model := range mi {
			if model.Name == name {
				return &model
			}
		}
	}

	// If it looks like an Ollama model, try to validate it
	if strings.Contains(name, ":") && ValidateOllamaModel(name) {
		return &ModelInfo{
			ID:     name,
			Name:   name,
			Family: "ollama",
		}
	}

	return nil
}
