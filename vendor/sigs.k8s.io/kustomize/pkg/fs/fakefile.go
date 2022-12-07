package fs

import (
	"bytes"
	"os"
)

var _ File = &FakeFile{}

// FakeFile implements File in-memory for tests.
type FakeFile struct {
	name    string
	content []byte
	dir     bool
	open    bool
}

// makeDir makes a fake directory.
func makeDir(name string) *FakeFile {
	return &FakeFile{name: name, dir: true}
}

// Close marks the fake file closed.
func (f *FakeFile) Close() error {
	f.open = false
	return nil
}

// Read never fails, and doesn't mutate p.
func (f *FakeFile) Read(p []byte) (n int, err error) {
	return len(p), nil
}

// Write saves the contents of the argument to memory.
func (f *FakeFile) Write(p []byte) (n int, err error) {
	f.content = p
	return len(p), nil
}

// ContentMatches returns true if v matches fake file's content.
func (f *FakeFile) ContentMatches(v []byte) bool {
	return bytes.Equal(v, f.content)
}

// GetContent the content of a fake file.
func (f *FakeFile) GetContent() []byte {
	return f.content
}

// Stat returns nil.
func (f *FakeFile) Stat() (os.FileInfo, error) {
	return nil, nil
}
