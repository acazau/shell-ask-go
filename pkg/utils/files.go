// pkg/utils/files.go
package utils

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// ReadFiles reads the content of multiple files and returns their combined content
func ReadFiles(files []string) (string, error) {
	var contents []string
	for _, file := range files {
		content, err := os.ReadFile(filepath.Clean(strings.TrimSpace(file)))
		if err != nil {
			return "", fmt.Errorf("failed to read file %s: %w", file, err)
		}
		contents = append(contents, fmt.Sprintf("=== %s ===\n%s", file, string(content)))
	}
	return strings.Join(contents, "\n\n"), nil
}

// FetchURLs fetches content from multiple URLs and returns their combined content
func FetchURLs(urls []string) (string, error) {
	var contents []string
	client := &http.Client{}

	for _, url := range urls {
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return "", fmt.Errorf("failed to create request for %s: %w", url, err)
		}

		// Add a user agent to avoid being blocked by some websites
		req.Header.Set("User-Agent", "Shell-Ask-Go/1.0")

		resp, err := client.Do(req)
		if err != nil {
			return "", fmt.Errorf("failed to fetch URL %s: %w", url, err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return "", fmt.Errorf("failed to fetch URL %s: status code %d", url, resp.StatusCode)
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return "", fmt.Errorf("failed to read response body from %s: %w", url, err)
		}

		contents = append(contents, fmt.Sprintf("=== %s ===\n%s", url, string(body)))
	}

	return strings.Join(contents, "\n\n"), nil
}
