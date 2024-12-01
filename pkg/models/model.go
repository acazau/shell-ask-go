package models

import (
	"context"
	"fmt"
)

// ModelInfo represents information about an AI model
type ModelInfo struct {
	ID     string  `json:"id"`
	RealID *string `json:"real_id,omitempty"`
	Description string `json:"description,omitempty"`
}

// ModelRegistry manages available AI models
type ModelRegistry struct {
	models map[string][]ModelInfo
}

// NewModelRegistry creates a new ModelRegistry
func NewModelRegistry() *ModelRegistry {
	return &ModelRegistry{
		models: map[string][]ModelInfo{
			"gpt": {
				{ID: "gpt-3.5-turbo", Description: "OpenAI's GPT-3.5 Turbo model"},
				{ID: "gpt-4-turbo", Description: "OpenAI's GPT-4 Turbo model"},
				{ID: "gpt-4o", Description: "OpenAI's GPT-4o model"},
			},
			"openai": {
				{ID: "gpt-4", RealID: stringPtr("gpt-4"), Description: "OpenAI's GPT-4 model"},
				{ID: "gpt-3.5-turbo", RealID: stringPtr("gpt-3.5-turbo"), Description: "OpenAI's GPT-3.5 Turbo model"},
			},
			"claude": {
				{ID: "claude-3-haiku", RealID: stringPtr("claude-3-haiku-20240307"), Description: "Anthropic's Claude 3 Haiku model"},
				{ID: "claude-3-sonnet", RealID: stringPtr("claude-3-sonnet-20240229"), Description: "Anthropic's Claude 3 Sonnet model"},
			},
		},
	}
}

// GetModelsForProvider returns models for a specific provider
func (mr *ModelRegistry) GetModelsForProvider(provider string) []ModelInfo {
	return mr.models[provider]
}

// GetAllModels returns all available models
func (mr *ModelRegistry) GetAllModels() []ModelInfo {
	var allModels []ModelInfo
	for _, providerModels := range mr.models {
		allModels = append(allModels, providerModels...)
	}
	return allModels
}

// SelectModel selects a model based on input
func (mr *ModelRegistry) SelectModel(input string) (string, error) {
	for _, models := range mr.models {
		for _, model := range models {
			if model.ID == input {
				if model.RealID != nil {
					return *model.RealID, nil
				}
				return model.ID, nil
			}
		}
	}

	// Check if input matches a provider
	if models, exists := mr.models[input]; exists && len(models) > 0 {
		return models[0].ID, nil
	}

	return "", fmt.Errorf("no model found for input: %s", input)
}

func stringPtr(s string) *string {
	return &s
}

// GetOllamaModels fetches models from Ollama (placeholder)
func GetOllamaModels(ctx context.Context) ([]ModelInfo, error) {
	return []ModelInfo{}, nil
}
