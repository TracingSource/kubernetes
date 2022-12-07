// +build windows

package util

import (
	winio "github.com/Microsoft/go-winio"
)

const (
	dockerSocket     = "//./pipe/docker_engine"         // The Docker socket is not CRI compatible
	containerdSocket = "//./pipe/containerd-containerd" // Proposed containerd named pipe for Windows
)

// isExistingSocket checks if path exists and is domain socket
func isExistingSocket(path string) bool {
	_, err := winio.DialPipe(path, nil)
	if err != nil {
		return false
	}

	return true
}
