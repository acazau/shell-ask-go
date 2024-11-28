// pkg/fetch/fetch.go
package fetch

import (
	"fmt"
	"io"
	"net/http"
	"strings"
)

type FetchResult struct {
	URL     string
	Content string
	Error   error
}

func FetchURLs(urls []string) []FetchResult {
	results := make([]FetchResult, len(urls))

	for i, url := range urls {
		result := FetchResult{URL: url}

		// If URL is a markdown converter URL, handle differently
		if strings.HasPrefix(url, "https://r.jina.ai/") {
			actualURL := strings.TrimPrefix(url, "https://r.jina.ai/")
			result.Content, result.Error = fetchMarkdown(actualURL)
		} else {
			result.Content, result.Error = fetchURL(url)
		}

		results[i] = result
	}

	return results
}

func fetchURL(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("markdown conversion error: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

func fetchMarkdown(url string) (string, error) {
	resp, err := http.Get("https://r.jina.ai/" + url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("markdown conversion error: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}
