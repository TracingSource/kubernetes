package system

// KernelValidatorHelper is an interface intended to help with os specific kernel validation
type KernelValidatorHelper interface {
	// GetKernelReleaseVersion gets the current kernel release version of the system
	GetKernelReleaseVersion() (string, error)
}
