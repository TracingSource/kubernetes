package git

import (
	"bytes"
	"os/exec"

	"github.com/pkg/errors"
	"sigs.k8s.io/kustomize/pkg/fs"
)

// Cloner is a function that can clone a git repo.
type Cloner func(repoSpec *RepoSpec) error

// ClonerUsingGitExec uses a local git install, as opposed
// to say, some remote API, to obtain a local clone of
// a remote repo.
func ClonerUsingGitExec(repoSpec *RepoSpec) error {
	gitProgram, err := exec.LookPath("git")
	if err != nil {
		return errors.Wrap(err, "no 'git' program on path")
	}
	repoSpec.cloneDir, err = fs.NewTmpConfirmedDir()
	if err != nil {
		return err
	}
	cmd := exec.Command(
		gitProgram,
		"clone",
		repoSpec.CloneSpec(),
		repoSpec.cloneDir.String())
	var out bytes.Buffer
	cmd.Stdout = &out
	err = cmd.Run()
	if err != nil {
		return errors.Wrapf(err, "trouble cloning %s", repoSpec.raw)
	}
	if repoSpec.ref == "" {
		return nil
	}
	cmd = exec.Command(gitProgram, "checkout", repoSpec.ref)
	cmd.Dir = repoSpec.cloneDir.String()
	err = cmd.Run()
	if err != nil {
		return errors.Wrapf(
			err, "trouble checking out href %s", repoSpec.ref)
	}
	return nil
}

// DoNothingCloner returns a cloner that only sets
// cloneDir field in the repoSpec.  It's assumed that
// the cloneDir is associated with some fake filesystem
// used in a test.
func DoNothingCloner(dir fs.ConfirmedDir) Cloner {
	return func(rs *RepoSpec) error {
		rs.cloneDir = dir
		return nil
	}
}
