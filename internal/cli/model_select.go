// internal/cli/model_select.go
package cli

import (
	"context"
	"fmt"
)

type ModelInfo struct {
	ID     string `json:"id"`
	RealID *string `json:"real_id,omitempty"`
}

var MODEL_MAP = map[string][]ModelInfo{
	"gpt": {
		{ID: "gpt-3.5-turbo"},
		{ID: "gpt-4-turbo"},
		{ID: "gpt-4o"},
		{ID: "gpt-4o-mini"},
	},
	"openai": {
		{ID: "gpt-4", RealID: stringPtr("gpt-4")},
		{ID: "gpt-3.5-turbo", RealID: stringPtr("gpt-3.5-turbo")},
	},
	"other": {
		{ID: "gpt-3.5-turbo", RealID: stringPtr("gpt-3.5-turbo")},
	},
	"claude": {
		{ID: "claude-3-haiku", RealID: stringPtr("claude-3-haiku-20240307")},
		{ID: "claude-3-sonnet", RealID: stringPtr("claude-3-sonnet-20240229")},
		{ID: "claude-3-opus", RealID: stringPtr("claude-3-opus-20240229")},
		{ID: "claude-3.5-haiku", RealID: stringPtr("claude-3-5-haiku-latest")},
		{ID: "claude-3.5-sonnet", RealID: stringPtr("claude-3-5-sonnet-latest")},
	},
	"gemini": {
		{ID: "gemini-1.5-pro", RealID: stringPtr("gemini-1.5-pro-latest")},
		{ID: "gemini-1.5-flash", RealID: stringPtr("gemini-1.5-flash-latest")},
		{ID: "gemini-pro"},
	},
	"groq": {
		{ID: "groq-llama3", RealID: stringPtr("groq-llama3-70b-8192")},
		{ID: "groq-llama3-8b", RealID: stringPtr("groq-llama3-8b-8192")},
		{ID: "groq-llama3-4b", RealID: stringPtr("groq-llama3-4b-8192")},
	},
}

func stringPtr(s string) *string {
	return &s
}

// SelectModel selects a model based on the provided input or interactive selection.
func SelectModel(input string) string {
	// Check if input matches any model ID
	for _, ms := range MODEL_MAP {
		for _, model := range ms {
			if model.ID == input {
				if model.RealID != nil {
					return *model.RealID
				}
				return model.ID
			}
		}
	}

	// Check if input matches any provider
	models, exists := MODEL_MAP[input]
	if exists && len(models) > 0 {
		return models[0].ID // Return the first model of the provider
	}

	// If no match, prompt for interactive selection
	return interactiveModelSelect()
}

func interactiveModelSelect() string {
    providers := getProviderList()
    
    selectedProvider := selectOption("Choose a provider:", providers)

    models := getModelsForProvider(selectedProvider)
    
    selectedModel := selectOption(fmt.Sprintf("Choose a %s model:", selectedProvider), models)

    // Return the selected model ID directly
    return selectedModel
}

func getProviderList() []string {
    var providers []string
    for provider := range MODEL_MAP {
        providers = append(providers, provider)
    }
    return providers
}

func getModelsForProvider(provider string) []string {
    models, exists := MODEL_MAP[provider]
    if !exists {
        return []string{}
    }
    modelIDs := make([]string, 0, len(models))
    for _, model := range models {
        modelIDs = append(modelIDs, model.ID)
    }
    return modelIDs
}

func selectOption(prompt string, options []string) string {
    for {
        fmt.Println(prompt)
        for i, option := range options {
            fmt.Printf("%d: %s\n", i+1, option)
        }
        var choice int
        fmt.Print("Enter the number of your choice: ")
        _, err := fmt.Scanln(&choice)
        if err != nil || choice < 1 || choice > len(options) {
            fmt.Println("Invalid choice, please try again.")
            continue
        }
        return options[choice-1]
    }
}

// GetAllModels returns all available models.
func GetAllModels(ctx context.Context, includeOllama bool) ([]ModelInfo, error) {
	var models []ModelInfo

	for _, ms := range MODEL_MAP {
		models = append(models, ms...)
	}

	if includeOllama {
		ollamaModels, err := getOllamaModels(ctx)
		if err != nil {
			return nil, err
		}
		models = append(models, ollamaModels...)
	}

	return models, nil
}

// getOllamaModels fetches models from Ollama.
func getOllamaModels(ctx context.Context) ([]ModelInfo, error) {
	// Placeholder for fetching models from Ollama
	return []ModelInfo{}, nil
}
