package stream

import (
	"bytes"
	"testing"
)

func TestOpenAIStreamProcessor_ProcessChunk(t *testing.T) {
	tests := []struct {
		name     string
		input    []byte
		expected string
	}{
		{
			name:     "empty input",
			input:    []byte{},
			expected: "",
		},
		{
			name:     "data prefix only",
			input:    []byte("data: "),
			expected: "",
		},
		{
			name:     "valid JSON with content",
			input:    []byte(`data: {"choices":[{"delta":{"content":"Hello"}}]}`),
			expected: "Hello",
		},
		{
			name:     "invalid JSON",
			input:    []byte(`data: {"choices":[{"delta":{}}]}`),
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &OpenAIStreamProcessor{}
			result, err := p.ProcessChunk(tt.input)
			if err != nil && tt.expected != "" {
				t.Errorf("ProcessChunk() error = %v, wantErr false", err)
			}
			if result != tt.expected {
				t.Errorf("ProcessChunk() got = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestProcessStream(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name: "single chunk",
			input: `data: {"choices":[{"delta":{"content":"Hello"}}]}
`,
			expected: "Hello",
		},
		{
			name: "multiple chunks",
			input: `data: {"choices":[{"delta":{"content":"Hello"}}]}
data: {"choices":[{"delta":{"content":" World!"}}]}
`,
			expected: "Hello World!",
		},
		{
			name:     "empty input",
			input:    "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := bytes.NewBufferString(tt.input)
			var writer bytes.Buffer
			p := &OpenAIStreamProcessor{}

			err := ProcessStream(reader, &writer, p)
			if err != nil {
				t.Errorf("ProcessStream() error = %v", err)
			}

			result := writer.String()
			if result != tt.expected {
				t.Errorf("ProcessStream() got = %v, want %v", result, tt.expected)
			}
		})
	}
}
