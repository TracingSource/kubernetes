package util

import (
	"os"
	"strings"

	"github.com/pkg/errors"
)

// GetHostname returns OS's hostname if 'hostnameOverride' is empty; otherwise, return 'hostnameOverride'
// NOTE: This function copied from pkg/util/node package to avoid external kubeadm dependency
func GetHostname(hostnameOverride string) (string, error) {
	hostName := hostnameOverride
	if len(hostName) == 0 {
		nodeName, err := os.Hostname()
		if err != nil {
			return "", errors.Wrap(err, "couldn't determine hostname")
		}
		hostName = nodeName
	}

	// Trim whitespaces first to avoid getting an empty hostname
	// For linux, the hostname is read from file /proc/sys/kernel/hostname directly
	hostName = strings.TrimSpace(hostName)
	if len(hostName) == 0 {
		return "", errors.New("empty hostname is invalid")
	}

	return strings.ToLower(hostName), nil
}
