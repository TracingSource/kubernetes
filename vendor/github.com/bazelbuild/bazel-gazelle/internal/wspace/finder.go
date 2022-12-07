/* Copyright 2016 The Bazel Authors. All rights reserved.

*/

// Package wspace provides functions to locate and modify a bazel WORKSPACE file.
package wspace

import (
	"os"
	"path/filepath"
	"strings"
)

const workspaceFile = "WORKSPACE"

// Find searches from the given dir and up for the WORKSPACE file
// returning the directory containing it, or an error if none found in the tree.
func Find(dir string) (string, error) {
	dir, err := filepath.Abs(dir)
	if err != nil {
		return "", err
	}

	for {
		_, err = os.Stat(filepath.Join(dir, workspaceFile))
		if err == nil {
			return dir, nil
		}
		if !os.IsNotExist(err) {
			return "", err
		}
		if strings.HasSuffix(dir, string(os.PathSeparator)) { // stop at root dir
			return "", os.ErrNotExist
		}
		dir = filepath.Dir(dir)
	}
}
