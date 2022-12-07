package coverage

import (
	"io"
)

// This is an implementation of testing.testDeps. It doesn't need to do anything, because
// no tests are actually run. It does need a concrete implementation of at least ImportPath,
// which is called unconditionally when running tests.
type fakeTestDeps struct{}

func (fakeTestDeps) ImportPath() string {
	return ""
}

func (fakeTestDeps) MatchString(pat, str string) (bool, error) {
	return false, nil
}

func (fakeTestDeps) StartCPUProfile(io.Writer) error {
	return nil
}

func (fakeTestDeps) StopCPUProfile() {}

func (fakeTestDeps) StartTestLog(io.Writer) {}

func (fakeTestDeps) StopTestLog() error {
	return nil
}

func (fakeTestDeps) WriteHeapProfile(io.Writer) error {
	return nil
}

func (fakeTestDeps) WriteProfileTo(string, io.Writer, int) error {
	return nil
}
