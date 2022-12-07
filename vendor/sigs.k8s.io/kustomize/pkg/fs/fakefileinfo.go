package fs

import (
	"os"
	"time"
)

var _ os.FileInfo = &Fakefileinfo{}

// Fakefileinfo implements Fakefileinfo using a fake in-memory filesystem.
type Fakefileinfo struct {
	*FakeFile
}

// Name returns the name of the file
func (fi *Fakefileinfo) Name() string { return fi.name }

// Size returns the size of the file
func (fi *Fakefileinfo) Size() int64 { return int64(len(fi.content)) }

// Mode returns the file mode
func (fi *Fakefileinfo) Mode() os.FileMode { return 0777 }

// ModTime returns the modification time
func (fi *Fakefileinfo) ModTime() time.Time { return time.Time{} }

// IsDir returns if it is a directory
func (fi *Fakefileinfo) IsDir() bool { return fi.dir }

// Sys should return underlying data source, but it now returns nil
func (fi *Fakefileinfo) Sys() interface{} { return nil }
