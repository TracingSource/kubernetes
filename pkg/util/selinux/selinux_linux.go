// +build linux

package selinux

import (
	selinux "github.com/opencontainers/selinux/go-selinux"
)

// SELinuxEnabled returns whether SELinux is enabled on the system.  SELinux
// has a tri-state:
//
// 1.  disabled: SELinux Kernel modules not loaded, SELinux policy is not
//     checked during Kernel MAC checks
// 2.  enforcing: Enabled; SELinux policy violations are denied and logged
//     in the audit log
// 3.  permissive: Enabled, but SELinux policy violations are permitted and
//     logged in the audit log
//
// SELinuxEnabled returns true if SELinux is enforcing or permissive, and
// false if it is disabled.
func SELinuxEnabled() bool {
	return selinux.GetEnabled()
}

// realSELinuxRunner is the real implementation of SELinuxRunner interface for
// Linux.
type realSELinuxRunner struct{}

var _ SELinuxRunner = &realSELinuxRunner{}

func (_ *realSELinuxRunner) Getfilecon(path string) (string, error) {
	if !SELinuxEnabled() {
		return "", nil
	}
	return selinux.FileLabel(path)
}

// SetFileLabel applies the SELinux label on the path or returns an error.
func SetFileLabel(path string, label string) error {
	return selinux.SetFileLabel(path, label)
}
