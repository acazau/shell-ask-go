// internal/providers/stream.go
package providers

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"strings"
)

func ProcessRequest(ctx context.Context, provider Provider, prompt string, stream bool) error {
	reader, err := provider.Complete(ctx, prompt, stream)
	if err != nil {
		return fmt.Errorf("failed to complete request: %w", err)
	}
	defer func() {
		if closer, ok := reader.(io.Closer); ok {
			closer.Close()
		}
	}()

	if !stream {
		_, err = io.Copy(os.Stdout, reader)
		return err
	}

	scanner := bufio.NewScanner(reader)
	scanner.Split(bufio.ScanWords)

	var buffer bytes.Buffer
	wordCount := 0

	for scanner.Scan() {
		word := scanner.Text()
		buffer.WriteString(word)
		wordCount++

		if wordCount > 0 && !strings.ContainsAny(word, ".,!?:;") {
			buffer.WriteString(" ")
		}

		if wordCount >= 10 || strings.ContainsAny(word, ".!?") {
			fmt.Print(buffer.String())
			buffer.Reset()
			wordCount = 0
		}
	}

	if buffer.Len() > 0 {
		fmt.Print(buffer.String())
	}
	fmt.Println()

	return scanner.Err()
}
