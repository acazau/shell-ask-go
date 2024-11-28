package fetch

import (
	"errors"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func loadTestData(t *testing.T, filename string) string {
	content, err := os.ReadFile(filepath.Join("testdata", filename))
	if err != nil {
		t.Fatalf("Failed to load test data %s: %v", filename, err)
	}
	return string(content)
}

func TestFetchURLs(t *testing.T) {
	// Create a custom HTTP client with mocked responses
	originalClient := http.DefaultClient
	defer func() { http.DefaultClient = originalClient }()

	mockClient := &http.Client{
		Transport: &roundTripperFunc{
			RoundTripFunc: func(req *http.Request) (*http.Response, error) {
				switch req.URL.String() {
				case "http://example.com":
					return &http.Response{
						StatusCode: 200,
						Body:       io.NopCloser(strings.NewReader(loadTestData(t, "example.com.html"))),
					}, nil
				case "http://example.com/markdown":
					return &http.Response{
						StatusCode: 451,
						Status:     "451 Unavailable For Legal Reasons",
					}, nil
				case "https://example.com":
					return &http.Response{
						StatusCode: 200,
						Body:       io.NopCloser(strings.NewReader(loadTestData(t, "example.com.simple.txt"))),
					}, nil
				case "https://r.jina.ai/markdown":
					return &http.Response{
						StatusCode: 200,
						Body:       io.NopCloser(strings.NewReader(loadTestData(t, "markdown.txt"))),
					}, nil
				default:
					return &http.Response{
						StatusCode: 404,
						Status:     "404 Not Found",
					}, nil
				}
			},
		},
	}
	http.DefaultClient = mockClient

	tests := []struct {
		name     string
		urls     []string
		expected []FetchResult
	}{
		{
			name: "Single_URL",
			urls: []string{"http://example.com"},
			expected: []FetchResult{
				{URL: "http://example.com", Content: loadTestData(t, "example.com.html"), Error: nil},
			},
		},
		{
			name: "Markdown_URL",
			urls: []string{"http://example.com/markdown"},
			expected: []FetchResult{
				{URL: "http://example.com/markdown", Content: "", Error: errors.New("markdown conversion error: 451 Unavailable For Legal Reasons")},
			},
		},
		{
			name: "Multiple URLs",
			urls: []string{"https://example.com", "https://r.jina.ai/markdown"},
			expected: []FetchResult{
				{URL: "https://example.com", Content: loadTestData(t, "example.com.simple.txt"), Error: nil},
				{URL: "https://r.jina.ai/markdown", Content: loadTestData(t, "markdown.txt"), Error: nil},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results := FetchURLs(tt.urls)
			for i, result := range results {
				assert.Equal(t, tt.expected[i].URL, result.URL)
				assert.Equal(t, tt.expected[i].Content, result.Content)
				assert.Equal(t, tt.expected[i].Error, result.Error)
			}
		})
	}
}

type roundTripperFunc struct {
	RoundTripFunc func(*http.Request) (*http.Response, error)
}

func (f *roundTripperFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f.RoundTripFunc(req)
}
