// Package fs provides a file system abstraction layer.
package fs

import (
	"io"
	"os"
)

// FileSystem groups basic os filesystem methods.
type FileSystem interface {
	Create(name string) (File, error)
	Mkdir(name string) error
	MkdirAll(name string) error
	RemoveAll(name string) error
	Open(name string) (File, error)
	IsDir(name string) bool
	CleanedAbs(path string) (ConfirmedDir, string, error)
	Exists(name string) bool
	Glob(pattern string) ([]string, error)
	ReadFile(name string) ([]byte, error)
	WriteFile(name string, data []byte) error
}

// File groups the basic os.File methods.
type File interface {
	io.ReadWriteCloser
	Stat() (os.FileInfo, error)
}
