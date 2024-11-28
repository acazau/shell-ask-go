// internal/providers/ollama.go
package providers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type OllamaProvider struct {
	host  string
	model string
}

func NewOllamaProvider(host, model string) *OllamaProvider {
	if host == "" {
		host = "http://localhost:11434"
	}
	return &OllamaProvider{
		host:  host,
		model: model,
	}
}

type ollamaRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"`
}

func (p *OllamaProvider) Complete(ctx context.Context, prompt string, stream bool) (io.ReadCloser, error) {
	reqBody := ollamaRequest{
		Model:  strings.TrimPrefix(p.model, "ollama-"),
		Prompt: prompt,
		Stream: stream,
	}

	body, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", fmt.Sprintf("%s/api/generate", p.host), strings.NewReader(string(body)))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("ollama API error: %s", resp.Status)
	}

	return resp.Body, nil
}

func (p *OllamaProvider) Name() string {
	return "ollama"
}
