// internal/providers/ollama.go
package providers

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/ollama/ollama/api"
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


func (p *OllamaProvider) Complete(ctx context.Context, prompt string, stream bool) (io.ReadCloser, error) {
	client, err := api.ClientFromEnvironment()
	if err != nil {
		return nil, err
	}

	streamPtr := &stream
	req := &api.GenerateRequest{
		Model:  strings.TrimPrefix(p.model, "ollama-"),
		Prompt: prompt,
		Stream: streamPtr,
	}

	respFunc := func(resp api.GenerateResponse) error {
		// Only print the response here; GenerateResponse has a number of other
		// interesting fields you want to examine.
		fmt.Println(resp.Response)
		return nil
	}

	err = client.Generate(ctx, req, respFunc)
	if err != nil {
		return nil, err
	}

	// Since we are not returning the response body directly, we need to handle it differently.
	// For now, we will return an empty io.ReadCloser and handle the response in the respFunc.
	return io.NopCloser(strings.NewReader("")), nil
}

func (p *OllamaProvider) Name() string {
	return "ollama"
}
