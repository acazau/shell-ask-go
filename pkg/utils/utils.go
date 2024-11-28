// pkg/utils/utils.go
package utils

import (
	"io"
	"os"
)

func IsPiped() bool {
	stat, _ := os.Stdin.Stat()
	return (stat.Mode() & os.ModeCharDevice) == 0
}

func ReadPipe() (string, error) {
	if !IsPiped() {
		return "", nil
	}

	bytes, err := io.ReadAll(os.Stdin)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}
