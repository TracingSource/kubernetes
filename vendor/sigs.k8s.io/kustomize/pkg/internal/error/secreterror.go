package error

import "fmt"

// SecretError represents error with a secret.
type SecretError struct {
	KustomizationPath string
	// ErrorMsg is an error message
	ErrorMsg string
}

func (e SecretError) Error() string {
	return fmt.Sprintf("Kustomization file [%s] encounters a secret error: %s\n", e.KustomizationPath, e.ErrorMsg)
}
