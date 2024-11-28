// pkg/stream/stream.go
package stream

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"strings"
)


type StreamProcessor interface {
	ProcessChunk([]byte) (string, error)
}

type OpenAIStreamProcessor struct{}

type openAIResponse struct {
	Choices []struct {
		Delta struct {
			Content string `json:"content"`
		} `json:"delta"`
	} `json:"choices"`
}

func (p *OpenAIStreamProcessor) ProcessChunk(chunk []byte) (string, error) {
	// Skip empty lines
	if len(chunk) == 0 {
		return "", nil
	}

	// Remove "data: " prefix if present
	data := strings.TrimPrefix(string(chunk), "data: ")

	// Parse JSON
	var resp openAIResponse
	if err := json.Unmarshal([]byte(data), &resp); err != nil {
		return "", err
	}

	if len(resp.Choices) == 0 {
		return "", nil
	}

	return resp.Choices[0].Delta.Content, nil
}

// Process a stream and write it to the output writer
func ProcessStream(reader io.Reader, writer io.Writer, processor StreamProcessor) error {
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		text, err := processor.ProcessChunk(scanner.Bytes())
		if err != nil {
			return err
		}
		if text != "" {
			if _, err := fmt.Fprint(writer, text); err != nil {
				return err
			}
		}
	}
	return scanner.Err()
}

