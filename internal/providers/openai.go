// internal/providers/openai.go
package providers

import (
	"context"
	"fmt"
	"io"
	"strings"

	openai "github.com/sashabaranov/go-openai"
)

type OpenAIProvider struct {
	client *openai.Client
	model  string
}

func NewOpenAIProvider(apiKey string, model string) *OpenAIProvider {
	return &OpenAIProvider{
		client: openai.NewClient(apiKey),
		model:  model,
	}
}

func (p *OpenAIProvider) Complete(ctx context.Context, prompt string, stream bool) (io.ReadCloser, error) {
	if stream {
		reader, err := p.streamCompletion(ctx, prompt)
		if err != nil {
			return nil, err
		}
		return reader, nil
	}
	reader, err := p.completion(ctx, prompt)
	if err != nil {
		return nil, err
	}
	return reader, nil
}

func (p *OpenAIProvider) Name() string {
	return "openai"
}

func (p *OpenAIProvider) streamCompletion(ctx context.Context, prompt string) (io.ReadCloser, error) {
	stream, err := p.client.CreateChatCompletionStream(
		ctx,
		openai.ChatCompletionRequest{
			Model: p.model,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: prompt,
				},
			},
		},
	)
	if err != nil {
		return nil, fmt.Errorf("stream completion error: %w", err)
	}
	defer stream.Close()

	pr, pw := io.Pipe()
	go func() {
		defer pw.Close()
		var result strings.Builder
		for {
			response, err := stream.Recv()
			if err != nil {
				return
			}
			result.WriteString(response.Choices[0].Delta.Content)
			_, _ = pw.Write([]byte(response.Choices[0].Delta.Content))
		}
	}()
	return pr, nil
}

func (p *OpenAIProvider) completion(ctx context.Context, prompt string) (io.ReadCloser, error) {
	resp, err := p.client.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model: p.model,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: prompt,
				},
			},
		},
	)
	if err != nil {
		return nil, fmt.Errorf("completion error: %w", err)
	}

	content := resp.Choices[0].Message.Content
	pr, pw := io.Pipe()
	go func() {
		defer pw.Close()
		_, _ = pw.Write([]byte(content))
	}()
	return pr, nil
}
