package hostpath

import (
	"os"

	v1 "k8s.io/api/core/v1"
	"k8s.io/kubernetes/pkg/volume/util/hostutil"
)

type hostPathTypeChecker interface {
	Exists() bool
	IsFile() bool
	MakeFile() error
	IsDir() bool
	MakeDir() error
	IsBlock() bool
	IsChar() bool
	IsSocket() bool
	GetPath() string
}

type fileTypeChecker struct {
	path   string
	exists bool
	hu     hostutil.HostUtils
}

func (ftc *fileTypeChecker) Exists() bool {
	exists, err := ftc.hu.PathExists(ftc.path)
	return exists && err == nil
}

func (ftc *fileTypeChecker) IsFile() bool {
	if !ftc.Exists() {
		return false
	}
	return !ftc.IsDir()
}

func (ftc *fileTypeChecker) MakeFile() error {
	return makeFile(ftc.path)
}

func (ftc *fileTypeChecker) IsDir() bool {
	if !ftc.Exists() {
		return false
	}
	pathType, err := ftc.hu.GetFileType(ftc.path)
	if err != nil {
		return false
	}
	return string(pathType) == string(v1.HostPathDirectory)
}

func (ftc *fileTypeChecker) MakeDir() error {
	return makeDir(ftc.path)
}

func (ftc *fileTypeChecker) IsBlock() bool {
	blkDevType, err := ftc.hu.GetFileType(ftc.path)
	if err != nil {
		return false
	}
	return string(blkDevType) == string(v1.HostPathBlockDev)
}

func (ftc *fileTypeChecker) IsChar() bool {
	charDevType, err := ftc.hu.GetFileType(ftc.path)
	if err != nil {
		return false
	}
	return string(charDevType) == string(v1.HostPathCharDev)
}

func (ftc *fileTypeChecker) IsSocket() bool {
	socketType, err := ftc.hu.GetFileType(ftc.path)
	if err != nil {
		return false
	}
	return string(socketType) == string(v1.HostPathSocket)
}

func (ftc *fileTypeChecker) GetPath() string {
	return ftc.path
}

func newFileTypeChecker(path string, hu hostutil.HostUtils) hostPathTypeChecker {
	return &fileTypeChecker{path: path, hu: hu}
}

// makeDir creates a new directory.
// If pathname already exists as a directory, no error is returned.
// If pathname already exists as a file, an error is returned.
func makeDir(pathname string) error {
	err := os.MkdirAll(pathname, os.FileMode(0755))
	if err != nil {
		if !os.IsExist(err) {
			return err
		}
	}
	return nil
}

// makeFile creates an empty file.
// If pathname already exists, whether a file or directory, no error is returned.
func makeFile(pathname string) error {
	f, err := os.OpenFile(pathname, os.O_CREATE, os.FileMode(0644))
	if f != nil {
		f.Close()
	}
	if err != nil {
		if !os.IsExist(err) {
			return err
		}
	}
	return nil
}
