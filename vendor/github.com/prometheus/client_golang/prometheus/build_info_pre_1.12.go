// +build !go1.12

package prometheus

// readBuildInfo is a wrapper around debug.ReadBuildInfo for Go versions before
// 1.12. Remove this whole file once the minimum supported Go version is 1.12.
func readBuildInfo() (path, version, sum string) {
	return "unknown", "unknown", "unknown"
}
