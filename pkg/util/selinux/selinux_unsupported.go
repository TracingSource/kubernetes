// +build !linux

package selinux

// SELinuxEnabled always returns false on non-linux platforms.
func SELinuxEnabled() bool {
	return false
}

// realSELinuxRunner is the NOP implementation of the SELinuxRunner interface.
type realSELinuxRunner struct{}

var _ SELinuxRunner = &realSELinuxRunner{}

func (_ *realSELinuxRunner) Getfilecon(path string) (string, error) {
	return "", nil
}

// FileLabel returns the SELinux label for this path or returns an error.
func SetFileLabel(path string, label string) error {
	return nil
}
