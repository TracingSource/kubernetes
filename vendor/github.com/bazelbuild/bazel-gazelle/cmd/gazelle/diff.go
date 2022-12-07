/* Copyright 2016 The Bazel Authors. All rights reserved.


*/

package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/bazelbuild/bazel-gazelle/config"
	"github.com/bazelbuild/bazel-gazelle/rule"
	"github.com/pmezard/go-difflib/difflib"
)

var exitError = fmt.Errorf("encountered changes while running diff")

func diffFile(c *config.Config, f *rule.File) error {
	rel, err := filepath.Rel(c.RepoRoot, f.Path)
	if err != nil {
		return fmt.Errorf("error getting old path for file %q: %v", f.Path, err)
	}
	rel = filepath.ToSlash(rel)

	date := "1970-01-01 00:00:00.000000000 +0000"
	diff := difflib.UnifiedDiff{
		Context:  3,
		FromDate: date,
		ToDate:   date,
	}

	if oldContent, err := ioutil.ReadFile(f.Path); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("error reading original file: %v", err)
	} else if err != nil {
		diff.FromFile = "/dev/null"
	} else if err == nil {
		diff.A = difflib.SplitLines(string(oldContent))
		if c.ReadBuildFilesDir == "" {
			path, err := filepath.Rel(c.RepoRoot, f.Path)
			if err != nil {
				return fmt.Errorf("error getting old path for file %q: %v", f.Path, err)
			}
			diff.FromFile = filepath.ToSlash(path)
		} else {
			diff.FromFile = f.Path
		}
	}

	newContent := f.Format()
	diff.B = difflib.SplitLines(string(newContent))
	outPath := findOutputPath(c, f)
	if c.WriteBuildFilesDir == "" {
		path, err := filepath.Rel(c.RepoRoot, f.Path)
		if err != nil {
			return fmt.Errorf("error getting new path for file %q: %v", f.Path, err)
		}
		diff.ToFile = filepath.ToSlash(path)
	} else {
		diff.ToFile = outPath
	}

	uc := getUpdateConfig(c)
	var out io.Writer = os.Stdout
	if uc.patchPath != "" {
		out = &uc.patchBuffer
	}
	if err := difflib.WriteUnifiedDiff(out, diff); err != nil {
		return fmt.Errorf("error diffing %s: %v", f.Path, err)
	}
	if ds, _ := difflib.GetUnifiedDiffString(diff); ds != "" {
		return exitError
	}

	return nil
}
