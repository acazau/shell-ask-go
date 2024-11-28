// internal/providers/provider.go
package providers

import (
	"context"
	"io"
)

type Provider interface {
	Complete(ctx context.Context, prompt string, stream bool) (io.ReadCloser, error)
	Name() string
}
