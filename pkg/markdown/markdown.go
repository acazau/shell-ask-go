// pkg/markdown/markdown.go
package markdown

import (
	"fmt"

	"github.com/charmbracelet/glamour"
)

type CLI struct {}

// Renderer handles markdown rendering with consistent styling
type Renderer struct {
	renderer *glamour.TermRenderer
}

func NewRenderer() (*Renderer, error) {
	r, err := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(100),
	)
	if err != nil {
		return nil, err
	}
	return &Renderer{renderer: r}, nil
}

func (r *Renderer) Render(markdown string) (string, error) {
	return r.renderer.Render(markdown)
}

func RenderMarkdown(input string) (string, error) {
	return input, nil
}

func RunWithMarkdown(cli *CLI, provider string, prompt string, noStream bool, commandOnly bool, fn func() error) error {
	err := fn()
	if err != nil {
		return fmt.Errorf("error during execution: %w", err)
	}

	return nil
}
