package util

import (
	gobuild "go/build"
	"path"
	"path/filepath"
	"reflect"
	"strings"
)

type empty struct{}

// CurrentPackage returns the go package of the current directory, or "" if it cannot
// be derived from the GOPATH.
func CurrentPackage() string {
	for _, root := range gobuild.Default.SrcDirs() {
		if pkg, ok := hasSubdir(root, "."); ok {
			return pkg
		}
	}
	return ""
}

func hasSubdir(root, dir string) (rel string, ok bool) {
	// ensure a tailing separator to properly compare on word-boundaries
	const sep = string(filepath.Separator)
	root = filepath.Clean(root)
	if !strings.HasSuffix(root, sep) {
		root += sep
	}

	// check whether root dir starts with root
	dir = filepath.Clean(dir)
	if !strings.HasPrefix(dir, root) {
		return "", false
	}

	// cut off root
	return filepath.ToSlash(dir[len(root):]), true
}

// BoilerplatePath uses the boilerplate in code-generator by calculating the relative path to it.
func BoilerplatePath() string {
	return path.Join(reflect.TypeOf(empty{}).PkgPath(), "/../../hack/boilerplate.go.txt")
}
