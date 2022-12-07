// Package error has contextual error types.
package error

import "fmt"

// ConfigmapError represents error with a configmap.
type ConfigmapError struct {
	Path     string
	ErrorMsg string
}

func (e ConfigmapError) Error() string {
	return fmt.Sprintf("Kustomization file [%s] encounters a configmap error: %s\n", e.Path, e.ErrorMsg)
}
