package providers

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

type OpenAIProvider struct {
	client *openai.Client
	model  string
}

func NewOpenAIProvider(apiKey string, model string) (*OpenAIProvider, error) {
	client := openai.NewClient(option.WithAPIKey(apiKey))
	return &OpenAIProvider{
		client: client,
		model:  model,
	}, nil
}

func (p *OpenAIProvider) Complete(ctx context.Context, prompt string, stream bool) (io.ReadCloser, error) {
	if stream {
		return p.streamCompletion(ctx, prompt)
	}
	return p.completion(ctx, prompt)
}

func (p *OpenAIProvider) Name() string {
	return "openai"
}

func (p *OpenAIProvider) streamCompletion(ctx context.Context, prompt string) (io.ReadCloser, error) {
	var output strings.Builder

	stream := p.client.Chat.Completions.NewStreaming(ctx, openai.ChatCompletionNewParams{
		Messages: openai.F([]openai.ChatCompletionMessageParamUnion{
			openai.UserMessage(prompt),
		}),
		Model: openai.F(p.model),
	})

	for stream.Next() {
		evt := stream.Current()
		if len(evt.Choices) > 0 {
			output.WriteString(evt.Choices[0].Delta.Content)
		}
	}

	if err := stream.Err(); err != nil {
		return nil, fmt.Errorf("stream completion error: %w", err)
	}

	return io.NopCloser(strings.NewReader(output.String())), nil
}

func (p *OpenAIProvider) completion(ctx context.Context, prompt string) (io.ReadCloser, error) {
	completion, err := p.client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Messages: openai.F([]openai.ChatCompletionMessageParamUnion{
			openai.UserMessage(prompt),
		}),
		Model: openai.F(p.model),
	})
	if err != nil {
		return nil, fmt.Errorf("completion error: %w", err)
	}

	return io.NopCloser(strings.NewReader(completion.Choices[0].Message.Content)), nil
}
