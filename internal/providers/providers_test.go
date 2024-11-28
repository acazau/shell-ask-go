// internal/providers/providers_test.go
package providers

import (
	"testing"
)

func TestOpenAIProvider(t *testing.T) {
	provider := NewOpenAIProvider("test-key", "gpt-4")
	if provider.Name() != "openai" {
		t.Errorf("expected provider name openai, got %s", provider.Name())
	}
}

func TestAnthropicProvider(t *testing.T) {
	provider := NewAnthropicProvider("test-key", "claude-3")
	if provider.Name() != "anthropic" {
		t.Errorf("expected provider name anthropic, got %s", provider.Name())
	}
}
