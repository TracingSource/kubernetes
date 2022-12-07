// +build go1.12

package prometheus

import "runtime/debug"

// readBuildInfo is a wrapper around debug.ReadBuildInfo for Go 1.12+.
func readBuildInfo() (path, version, sum string) {
	path, version, sum = "unknown", "unknown", "unknown"
	if bi, ok := debug.ReadBuildInfo(); ok {
		path = bi.Main.Path
		version = bi.Main.Version
		sum = bi.Main.Sum
	}
	return
}
