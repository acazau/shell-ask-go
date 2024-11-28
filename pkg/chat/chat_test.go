// pkg/chat/chat_test.go
package chat

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestSaveAndLoadChat(t *testing.T) {
	// Create test chat
	testChat := &Chat{
		Messages: []Message{
			{Role: "user", Content: "test message"},
			{Role: "assistant", Content: "test response"},
		},
		Model: "gpt-4",
	}

	// Save chat
	err := SaveChat(testChat)
	if err != nil {
		t.Fatal(err)
	}

	// Get cache dir
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		t.Fatal(err)
	}

	// Read saved chat file
	chatPath := filepath.Join(cacheDir, "shell-ask", "chat.json")
	data, err := os.ReadFile(chatPath)
	if err != nil {
		t.Fatal(err)
	}

	// Parse and verify
	var savedChat Chat
	err = json.Unmarshal(data, &savedChat)
	if err != nil {
		t.Fatal(err)
	}

	if len(savedChat.Messages) != len(testChat.Messages) {
		t.Errorf("expected %d messages, got %d", len(testChat.Messages), len(savedChat.Messages))
	}
}
