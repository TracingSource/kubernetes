package autorest

import (
	"fmt"
	"runtime"
)

const number = "v13.0.0"

var (
	userAgent = fmt.Sprintf("Go/%s (%s-%s) go-autorest/%s",
		runtime.Version(),
		runtime.GOARCH,
		runtime.GOOS,
		number,
	)
)

// UserAgent returns a string containing the Go version, system architecture and OS, and the go-autorest version.
func UserAgent() string {
	return userAgent
}

// Version returns the semantic version (see http://semver.org).
func Version() string {
	return number
}
