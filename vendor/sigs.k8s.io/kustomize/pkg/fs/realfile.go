package fs

import (
	"os"
)

var _ File = &realFile{}

// realFile implements File using the local filesystem.
type realFile struct {
	file *os.File
}

// Close closes a file.
func (f *realFile) Close() error { return f.file.Close() }

// Read reads a file's content.
func (f *realFile) Read(p []byte) (n int, err error) { return f.file.Read(p) }

// Write writes bytes to a file
func (f *realFile) Write(p []byte) (n int, err error) { return f.file.Write(p) }

// Stat returns an interface which has all the information regarding the file.
func (f *realFile) Stat() (os.FileInfo, error) { return f.file.Stat() }
