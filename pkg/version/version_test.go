package version

import (
	"fmt"
	"runtime"
	"testing"
)

func TestGetVersionInfo(t *testing.T) {
	expected := "shell-ask version 0.1.0\ngit commit: unknown\nbuild date: unknown\ngo version: %s\nplatform: %s/%s\n"
	actual := GetVersionInfo()

	if actual != fmt.Sprintf(expected, runtime.Version(), runtime.GOOS, runtime.GOARCH) {
		t.Errorf("GetVersionInfo() = %v; want %v", actual, expected)
	}
}
