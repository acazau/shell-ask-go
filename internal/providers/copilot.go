// internal/providers/copilot.go
package providers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

const (
	copilotCompletionAPI = "https://api.githubcopilot.com/chat/completions"
)

type CopilotProvider struct {
	client *http.Client
	token  string
	model  string
}

func NewCopilotProvider(token string, model string) (*CopilotProvider, error) {
	// Define valid models
	validModels := map[string]string{
		"gpt-4":             "gpt-4",
		"4":                 "gpt-4",
		"gpt-4o":            "gpt-4o",
		"4o":                "gpt-4o",
		"o1-mini":           "o1-mini",
		"o1-preview":        "o1-preview",
		"claude-3.5-sonnet": "claude-3.5-sonnet",
	}

	// Clean up model name
	model = strings.TrimPrefix(model, "copilot-")

	// Look up the canonical model name
	if canonicalModel, ok := validModels[model]; ok {
		model = canonicalModel
	} else if model == "" {
		model = "gpt-4" // default model
	} else {
		return nil, fmt.Errorf("unsupported Copilot model: %s. Valid models are: gpt-4, gpt-4o, o1-mini, o1-preview, claude-3.5-sonnet", model)
	}

	return &CopilotProvider{
		client: &http.Client{},
		token:  token,
		model:  model,
	}, nil
}

type copilotRequest struct {
	Intent      bool             `json:"intent"`
	Model       string           `json:"model"`
	N           int              `json:"n"`
	Stream      bool             `json:"stream"`
	Temperature float32          `json:"temperature"`
	TopP        int              `json:"top_p"`
	Messages    []copilotMessage `json:"messages"`
	MaxTokens   int              `json:"max_tokens"`
}

type copilotMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type copilotResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

func (p *CopilotProvider) Complete(ctx context.Context, prompt string, stream bool) (io.ReadCloser, error) {
	reqBody := copilotRequest{
		Intent:      true,
		Model:       p.model,
		N:           1,
		Stream:      stream,
		Temperature: 0.1,
		TopP:        1,
		MaxTokens:   8192,
		Messages: []copilotMessage{
			{Role: "user", Content: prompt},
		},
	}

	body, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	fmt.Printf("Debug - Request Body: %s\n", string(body))

	req, err := http.NewRequestWithContext(
		ctx,
		"POST",
		copilotCompletionAPI,
		bytes.NewReader(body),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Add required headers
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", p.token))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("OpenAI-Intent", "conversation-panel")
	req.Header.Set("OpenAI-Organization", "github-copilot")
	req.Header.Set("Editor-Version", "vscode/1.88.0")
	req.Header.Set("Editor-Plugin-Version", "copilot-chat/0.14.2024032901")
	req.Header.Set("User-Agent", "GitHubCopilotChat/0.14.2024032901")
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Accept-Encoding", "gzip,deflate,br")
	req.Header.Set("X-GitHub-Api-Version", "2023-07-07")
	req.Header.Set("Copilot-Integration-Id", "vscode-chat")

	resp, err := p.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		defer resp.Body.Close()
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("copilot API error: %s - %s", resp.Status, string(body))
	}

	if !stream {
		var response copilotResponse
		if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
			resp.Body.Close()
			return nil, fmt.Errorf("failed to decode response: %w", err)
		}

		if len(response.Choices) == 0 {
			resp.Body.Close()
			return nil, fmt.Errorf("no completion choices returned")
		}

		return io.NopCloser(strings.NewReader(response.Choices[0].Message.Content)), nil
	}

	return resp.Body, nil
}

func (p *CopilotProvider) Name() string {
	return "copilot"
}
