// pkg/chat/chat.go
package chat

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type Chat struct {
	Messages []Message `json:"messages"`
	Model    string    `json:"model"`
}

func SaveChat(chat *Chat) error {
	dir, err := os.UserCacheDir()
	if err != nil {
		return err
	}

	path := filepath.Join(dir, "shell-ask", "chat.json")
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}

	data, err := json.Marshal(chat)
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}
