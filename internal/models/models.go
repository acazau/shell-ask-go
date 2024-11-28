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
}

// GetAllModels returns all available models, optionally including Ollama models
func GetAllModels(includeOllama bool) []ModelInfo {
	var allModels []ModelInfo

	// Add all models from ModelMap
	for _, models := range ModelMap {
		allModels = append(allModels, models...)
	}

	// Add Ollama models if requested
	if includeOllama {
		ollamaModels, err := GetOllamaModels()
		if err == nil {
			allModels = append(allModels, ollamaModels...)
		}
	}

	// Sort models by ID for consistent ordering
	sort.Slice(allModels, func(i, j int) bool {
		return allModels[i].ID < allModels[j].ID
	})

	return allModels
}

// GetCheapModel returns a cheaper alternative for expensive models
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
	default:
		return modelID
	}
}

func GetModels() ([]ModelInfo, error) {
	models := []ModelInfo{
		{
			Name:        "gpt-4",
			Description: "Most capable GPT-4 model",
			Family:      "GPT-4",
		},
		{
			Name:        "gpt-3.5-turbo",
			Description: "Most capable GPT-3.5 model",
			Family:      "GPT-3.5",
		},
	}
	return models, nil
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

	var models Models
	if err := json.Unmarshal(body, &models); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return models.Models, nil
}
