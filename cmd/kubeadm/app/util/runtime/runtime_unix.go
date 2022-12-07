// +build !windows

package util

import (
	"os"
)

const (
	dockerSocket     = "/var/run/docker.sock" // The Docker socket is not CRI compatible
	containerdSocket = "/run/containerd/containerd.sock"
)

// isExistingSocket checks if path exists and is domain socket
func isExistingSocket(path string) bool {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false
	}

	return fileInfo.Mode()&os.ModeSocket != 0
}
