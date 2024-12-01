package providers

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

type GeminiProvider struct {
	client *genai.Client
	model  string
}

func NewGeminiProvider(apiKey, model string) (*GeminiProvider, error) {
	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return nil, err
	}

	return &GeminiProvider{
		client: client,
		model:  model,
	}, nil
}

func (p *GeminiProvider) Complete(ctx context.Context, prompt string, stream bool) (io.ReadCloser, error) {
	model := p.client.GenerativeModel(p.model)

	resp, err := model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		return nil, err
	}

	if len(resp.Candidates) == 0 || resp.Candidates[0].Content == nil {
		return nil, fmt.Errorf("no content returned from Gemini API")
	}

	text := fmt.Sprintf("%v", resp.Candidates[0].Content.Parts[0])
	return io.NopCloser(strings.NewReader(text)), nil
}

func (p *GeminiProvider) Name() string {
	return "gemini"
}
