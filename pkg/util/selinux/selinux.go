package selinux

// Note: the libcontainer SELinux package is only built for Linux, so it is
// necessary to have a NOP wrapper which is built for non-Linux platforms to
// allow code that links to this package not to differentiate its own methods
// for Linux and non-Linux platforms.
//
// SELinuxRunner wraps certain libcontainer SELinux calls. For more
// information, see:
//
// https://github.com/opencontainers/runc/blob/master/libcontainer/selinux/selinux.go
type SELinuxRunner interface {
	// Getfilecon returns the SELinux context for the given path or returns an
	// error.
	Getfilecon(path string) (string, error)
}

// NewSELinuxRunner returns a new SELinuxRunner appropriate for the platform.
// On Linux, all methods short-circuit and return NOP values if SELinux is
// disabled. On non-Linux platforms, a NOP implementation is returned.
func NewSELinuxRunner() SELinuxRunner {
	return &realSELinuxRunner{}
}
