// internal/providers/factory.go
package providers

import (
	"fmt"
	"os"
	"strings"

	"github.com/acazau/shell-ask-go/internal/config"
	"github.com/acazau/shell-ask-go/internal/copilot"
)

// parseModelID splits a model ID into provider and model parts
// Format: provider/model or just model
func parseModelID(modelID string) (provider, model string) {
	parts := strings.SplitN(modelID, "/", 2)
	if len(parts) == 2 {
		return parts[0], parts[1]
	}
	return "", modelID
}

func InitializeProvider(modelID string) (Provider, error) {
	providerPrefix, modelName := parseModelID(modelID)
	model := modelName // For Copilot, we want to keep the original model name

	// If provider is explicitly specified, use it
	if providerPrefix != "" {
		switch providerPrefix {
		case "openai":
			return NewOpenAIProvider(os.Getenv("OPENAI_API_KEY"), model)
		case "anthropic":
			return NewAnthropicProvider(os.Getenv("ANTHROPIC_API_KEY"), model), nil
		case "gemini":
			return NewGeminiProvider(os.Getenv("GOOGLE_API_KEY"), model)
		case "groq":
			return NewGroqProvider(os.Getenv("GROQ_API_KEY"), model), nil
		case "ollama":
			return NewOllamaProvider("http://localhost:11434", model), nil
		case "copilot":
			copilotClient := copilot.New(config.GetConfigDir())
			token, err := copilotClient.GetAPIToken()
			if err != nil {
				return nil, fmt.Errorf("failed to get Copilot token: %w", err)
			}
			provider, err := NewCopilotProvider(token, model)
			if err != nil {
				return nil, fmt.Errorf("failed to initialize Copilot provider: %w", err)
			}
			return provider, nil
		default:
			return nil, fmt.Errorf("unsupported provider: %s", providerPrefix)
		}
	}

	// If no provider specified, infer from model name
	switch {
	case strings.HasPrefix(model, "openai"):
		return NewOpenAIProvider(os.Getenv("OPENAI_API_KEY"), model)
	case strings.HasPrefix(model, "anthropic"):
		return NewAnthropicProvider(os.Getenv("ANTHROPIC_API_KEY"), model), nil
	case strings.HasPrefix(model, "gemini"):
		return NewGeminiProvider(os.Getenv("GOOGLE_API_KEY"), model)
	case strings.HasPrefix(model, "groq"):
		return NewGroqProvider(os.Getenv("GROQ_API_KEY"), model), nil
	case strings.HasPrefix(model, "llama"):
		return NewOllamaProvider("http://localhost:11434", model), nil
	case strings.HasPrefix(model, "copilot-"):
		copilotClient := copilot.New(config.GetConfigDir())
		token, err := copilotClient.GetAPIToken()
		if err != nil {
			return nil, fmt.Errorf("failed to get Copilot token: %w", err)
		}
		provider, err := NewCopilotProvider(token, model)
		if err != nil {
			return nil, fmt.Errorf("failed to initialize Copilot provider: %w", err)
		}
		return provider, nil
	default:
		return nil, fmt.Errorf("unsupported model: %s", model)
	}
}
