package kustomize

import (
	"io"

	"k8s.io/cli-runtime/pkg/kustomize/k8sdeps"
	"sigs.k8s.io/kustomize/pkg/commands/build"
	"sigs.k8s.io/kustomize/pkg/fs"
)

// RunKustomizeBuild runs kustomize build given a filesystem and a path
func RunKustomizeBuild(out io.Writer, fSys fs.FileSystem, path string) error {
	f := k8sdeps.NewFactory()
	o := build.NewOptions(path, "")
	return o.RunBuild(out, fSys, f.ResmapF, f.TransformerF)
}
