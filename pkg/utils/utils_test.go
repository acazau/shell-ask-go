// pkg/utils/utils_test.go
package utils

import (
	"os"
	"testing"
)

func TestIsPiped(t *testing.T) {
	// Save original stdin
	originalStdin := os.Stdin
	defer func() { os.Stdin = originalStdin }()

	// Create a pipe
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}

	// Write some data to the pipe
	testData := []byte("test input")
	go func() {
		defer w.Close()
		w.Write(testData)
	}()

	// Set stdin to read end of pipe
	os.Stdin = r

	// Test ReadPipe
	input, err := ReadPipe()
	if err != nil {
		t.Fatal(err)
	}

	if input != string(testData) {
		t.Errorf("expected input %s, got %s", string(testData), input)
	}
}
