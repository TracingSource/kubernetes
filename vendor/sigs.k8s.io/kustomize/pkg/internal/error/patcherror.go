package error

import (
	"fmt"
)

// PatchError represents error during Patch.
type PatchError struct {
	KustomizationPath string
	PatchFilepath     string
	ErrorMsg          string
}

func (e PatchError) Error() string {
	return fmt.Sprintf("Kustomization file [%s] encounters a patch error for [%s]: %s\n", e.KustomizationPath, e.PatchFilepath, e.ErrorMsg)
}
