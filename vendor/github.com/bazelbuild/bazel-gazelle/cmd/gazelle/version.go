package main

import (
	"io/ioutil"
	"log"
	"path/filepath"
	"regexp"

	"github.com/bazelbuild/bazel-gazelle/config"
	"github.com/bazelbuild/bazel-gazelle/internal/version"
	"github.com/bazelbuild/bazel-gazelle/repo"
)

var minimumRulesGoVersion = version.Version{0, 19, 0}

// checkRulesGoVersion checks whether a compatible version of rules_go is
// being used in the workspace. A message will be logged if an incompatible
// version is found.
//
// Note that we can't always determine the version of rules_go in use. Also,
// if we find an incompatible version, we shouldn't bail out since the
// incompatibility may not matter in the current workspace.
func checkRulesGoVersion(repoRoot string) {
	const message = `Gazelle may not be compatible with this version of rules_go.
Update io_bazel_rules_go to a newer version in your WORKSPACE file.`

	rulesGoPath, err := repo.FindExternalRepo(repoRoot, config.RulesGoRepoName)
	if err != nil {
		return
	}
	defBzlPath := filepath.Join(rulesGoPath, "go", "def.bzl")
	defBzlContent, err := ioutil.ReadFile(defBzlPath)
	if err != nil {
		return
	}
	versionRe := regexp.MustCompile(`(?m)^RULES_GO_VERSION = ['"]([0-9.]*)['"]`)
	match := versionRe.FindSubmatch(defBzlContent)
	if match == nil {
		log.Printf("RULES_GO_VERSION not found in @%s//go:def.bzl.\n%s", config.RulesGoRepoName, message)
		return
	}
	vstr := string(match[1])
	v, err := version.ParseVersion(vstr)
	if err != nil {
		log.Printf("RULES_GO_VERSION %q could not be parsed in @%s//go:def.bzl.\n%s", vstr, config.RulesGoRepoName, message)
	}
	if v.Compare(minimumRulesGoVersion) < 0 {
		log.Printf("Found RULES_GO_VERSION %s. Minimum compatible version is %s.\n%s", v, minimumRulesGoVersion, message)
	}
}
