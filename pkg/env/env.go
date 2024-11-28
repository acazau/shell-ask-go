// pkg/env/env.go
package env

import (
	"os"
	"runtime"
)

// IsInteractive checks if the program is running in an interactive terminal
func IsInteractive() bool {
	if runtime.GOOS == "windows" {
		return false
	}
	fileInfo, _ := os.Stdout.Stat()
	return (fileInfo.Mode() & os.ModeCharDevice) != 0
}

// IsPiped checks if the program is receiving piped input
func IsPiped() bool {
	fileInfo, _ := os.Stdin.Stat()
	return (fileInfo.Mode() & os.ModeCharDevice) == 0
}

// GetConfigDir returns the configuration directory path
func GetConfigDir() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return configDir + "/shell-ask", nil
}

// GetCacheDir returns the cache directory path
func GetCacheDir() (string, error) {
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		return "", err
	}
	return cacheDir + "/shell-ask", nil
}

// EnsureConfigDir ensures the configuration directory exists
func EnsureConfigDir() error {
	dir, err := GetConfigDir()
	if err != nil {
		return err
	}
	return os.MkdirAll(dir, 0755)
}

// EnsureCacheDir ensures the cache directory exists
func EnsureCacheDir() error {
	dir, err := GetCacheDir()
	if err != nil {
		return err
	}
	return os.MkdirAll(dir, 0755)
}
