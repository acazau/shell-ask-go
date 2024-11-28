// internal/providers/groq.go
package providers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type GroqProvider struct {
	apiKey string
	model  string
}

func NewGroqProvider(apiKey, model string) *GroqProvider {
	return &GroqProvider{
		apiKey: apiKey,
		model:  model,
	}
}

type groqRequest struct {
	Model    string        `json:"model"`
	Messages []groqMessage `json:"messages"`
	Stream   bool          `json:"stream"`
}

type groqMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

func (p *GroqProvider) Complete(ctx context.Context, prompt string, stream bool) (io.ReadCloser, error) {
	reqBody := groqRequest{
		Model: strings.TrimPrefix(p.model, "groq-"),
		Messages: []groqMessage{
			{Role: "user", Content: prompt},
		},
		Stream: stream,
	}

	body, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", "https://api.groq.com/openai/v1/chat/completions", strings.NewReader(string(body)))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", p.apiKey))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("groq API error: %s", resp.Status)
	}

	return resp.Body, nil
}

func (p *GroqProvider) Name() string {
	return "groq"
}
