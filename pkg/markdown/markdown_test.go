package markdown

import (
	"testing"
)

func TestRenderMarkdown(t *testing.T) {
	input := "# Hello World\nThis is a test."
	expectedOutput := "# Hello World\nThis is a test."

	output, err := RenderMarkdown(input)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if output != expectedOutput {
		t.Errorf("Expected output to be '%s', but got '%s'", expectedOutput, output)
	}
}
