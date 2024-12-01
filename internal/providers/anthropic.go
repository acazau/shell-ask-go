// internal/providers/anthropic.go
package providers

import (
	"context"
	"io"
	"strings"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
)

type AnthropicProvider struct {
	client *anthropic.Client
	model  string
}

// NewAnthropicProvider creates a new Anthropic provider
// model should be one of: claude-3-opus-20240229, claude-3-sonnet-20240229, claude-3-haiku-20240307
func NewAnthropicProvider(apiKey, model string) *AnthropicProvider {
	client := anthropic.NewClient(
		option.WithAPIKey(apiKey),
	)
	return &AnthropicProvider{
		client: client,
		model:  model,
	}
}

func (p *AnthropicProvider) Complete(ctx context.Context, prompt string, stream bool) (io.ReadCloser, error) {
	req := anthropic.MessageNewParams{
		MaxTokens: anthropic.Int(1024),
		Messages: anthropic.F([]anthropic.MessageParam{
			anthropic.NewUserMessage(anthropic.NewTextBlock("You are a helpful AI assistant.\n\n" + prompt)),
		}),
		Model:         anthropic.F(p.model),
		StopSequences: anthropic.F([]string{"```\n"}),
	}

	if stream {
		stream := p.client.Messages.NewStreaming(ctx, req)

		reader, writer := io.Pipe()

		go func() {
			defer writer.Close()
			for stream.Next() {
				event := stream.Current()

				switch delta := event.Delta.(type) {
				case anthropic.ContentBlockDeltaEventDelta:
					if delta.Text != "" {
						writer.Write([]byte(delta.Text))
					}
				case anthropic.MessageDeltaEventDelta:
					if delta.StopSequence != "" {
						writer.Write([]byte(delta.StopSequence))
					}
				}
			}

			if stream.Err() != nil {
				writer.CloseWithError(stream.Err())
			}
		}()

		return reader, nil
	}

	resp, err := p.client.Messages.New(ctx, req)
	if err != nil {
		return nil, err
	}

	if len(resp.Content) > 0 {
		return io.NopCloser(strings.NewReader(resp.Content[0].Text)), nil
	}
	return io.NopCloser(strings.NewReader("")), nil
}

func (p *AnthropicProvider) Name() string {
	return "anthropic"
}
