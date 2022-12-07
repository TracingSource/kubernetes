/* Copyright 2016 The Bazel Authors. All rights reserved.


*/

package main

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/bazelbuild/bazel-gazelle/config"
	"github.com/bazelbuild/bazel-gazelle/rule"
)

func fixFile(c *config.Config, f *rule.File) error {
	outPath := findOutputPath(c, f)
	if err := os.MkdirAll(filepath.Dir(outPath), 0777); err != nil {
		return err
	}
	return ioutil.WriteFile(outPath, f.Format(), 0666)
}
