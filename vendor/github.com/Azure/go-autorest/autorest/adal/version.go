package adal

import (
	"fmt"
	"runtime"
)

const number = "v1.0.0"

var (
	ua = fmt.Sprintf("Go/%s (%s-%s) go-autorest/adal/%s",
		runtime.Version(),
		runtime.GOARCH,
		runtime.GOOS,
		number,
	)
)

// UserAgent returns a string containing the Go version, system architecture and OS, and the adal version.
func UserAgent() string {
	return ua
}

// AddToUserAgent adds an extension to the current user agent
func AddToUserAgent(extension string) error {
	if extension != "" {
		ua = fmt.Sprintf("%s %s", ua, extension)
		return nil
	}
	return fmt.Errorf("Extension was empty, User Agent remained as '%s'", ua)
}
