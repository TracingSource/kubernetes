package error

import (
	"fmt"
)

// KustomizationError represents an error with a kustomization.
type KustomizationError struct {
	KustomizationPath string
	ErrorMsg          string
}

func (ke KustomizationError) Error() string {
	return fmt.Sprintf("Kustomization File [%s]: %s\n", ke.KustomizationPath, ke.ErrorMsg)
}

// KustomizationErrors collects all errors.
type KustomizationErrors struct {
	kErrors []error
}

func (ke *KustomizationErrors) Error() string {
	errormsg := ""
	for _, e := range ke.kErrors {
		errormsg += e.Error() + "\n"
	}
	return errormsg
}

// Append adds error to a collection of errors.
func (ke *KustomizationErrors) Append(e error) {
	ke.kErrors = append(ke.kErrors, e)
}

// Get returns all collected errors.
func (ke *KustomizationErrors) Get() []error {
	return ke.kErrors
}

// BatchAppend adds all errors from another KustomizationErrors
func (ke *KustomizationErrors) BatchAppend(e KustomizationErrors) {
	for _, err := range e.Get() {
		ke.kErrors = append(ke.kErrors, err)
	}
}
