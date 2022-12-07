/* Copyright 2016 The Bazel Authors. All rights reserved.


*/

package main

import (
	"os"

	"github.com/bazelbuild/bazel-gazelle/config"
	"github.com/bazelbuild/bazel-gazelle/rule"
)

func printFile(c *config.Config, f *rule.File) error {
	content := f.Format()
	_, err := os.Stdout.Write(content)
	return err
}
