// pkg/version/version.go
package version

import (
	"fmt"
	"runtime"
)

var (
	// Version is the current version of shell-ask
	Version = "0.1.0"

	// GitCommit is the git commit hash of the current build
	GitCommit = "unknown"

	// BuildDate is the date when this binary was built
	BuildDate = "unknown"
)

// GetVersionInfo returns formatted version information
func GetVersionInfo() string {
	return fmt.Sprintf("shell-ask version %s\n"+
		"git commit: %s\n"+
		"build date: %s\n"+
		"go version: %s\n"+
		"platform: %s/%s\n",
		Version,
		GitCommit,
		BuildDate,
		runtime.Version(),
		runtime.GOOS,
		runtime.GOARCH,
	)
}
