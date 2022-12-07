package util

import (
	"io/ioutil"
	"os"
	"path/filepath"
)

// IoUtil is a mockable util for common IO operations
type IoUtil interface {
	ReadFile(filename string) ([]byte, error)
	ReadDir(dirname string) ([]os.FileInfo, error)
	Lstat(name string) (os.FileInfo, error)
	EvalSymlinks(path string) (string, error)
}

type osIOHandler struct{}

//NewIOHandler Create a new IoHandler implementation
func NewIOHandler() IoUtil {
	return &osIOHandler{}
}

func (handler *osIOHandler) ReadFile(filename string) ([]byte, error) {
	return ioutil.ReadFile(filename)
}
func (handler *osIOHandler) ReadDir(dirname string) ([]os.FileInfo, error) {
	return ioutil.ReadDir(dirname)
}
func (handler *osIOHandler) Lstat(name string) (os.FileInfo, error) {
	return os.Lstat(name)
}
func (handler *osIOHandler) EvalSymlinks(path string) (string, error) {
	return filepath.EvalSymlinks(path)
}
