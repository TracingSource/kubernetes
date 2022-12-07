// +build !coverage

package coverage

// InitCoverage is illegal when not running with coverage.
func InitCoverage(name string) {
	panic("Called InitCoverage when not built with coverage instrumentation.")
}

// FlushCoverage is a no-op when not running with coverage.
func FlushCoverage() {

}
