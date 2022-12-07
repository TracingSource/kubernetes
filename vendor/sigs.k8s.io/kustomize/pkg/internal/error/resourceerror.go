package error

import "fmt"

// ResourceError represents error in a resource.
type ResourceError struct {
	KustomizationPath string
	ResourceFilepath  string
	ErrorMsg          string
}

func (e ResourceError) Error() string {
	return fmt.Sprintf("Kustomization file [%s] encounters a resource error for [%s]: %s\n", e.KustomizationPath, e.ResourceFilepath, e.ErrorMsg)
}
