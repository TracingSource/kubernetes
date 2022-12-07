// Package loader has a data loading interface and various implementations.
package loader

import (
	"sigs.k8s.io/kustomize/pkg/fs"
	"sigs.k8s.io/kustomize/pkg/git"
	"sigs.k8s.io/kustomize/pkg/ifc"
)

// NewLoader returns a Loader.
func NewLoader(path string, fSys fs.FileSystem) (ifc.Loader, error) {
	repoSpec, err := git.NewRepoSpecFromUrl(path)
	if err == nil {
		return newLoaderAtGitClone(
			repoSpec, fSys, nil, git.ClonerUsingGitExec)
	}
	root, err := demandDirectoryRoot(fSys, path)
	if err != nil {
		return nil, err
	}
	return newLoaderAtConfirmedDir(
		root, fSys, nil, git.ClonerUsingGitExec), nil
}
