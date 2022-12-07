// +build !providerless
// +build linux

package vsphere_volume

import (
	"fmt"

	"k8s.io/utils/mount"
)

func verifyDevicePath(path string) (string, error) {
	if pathExists, err := mount.PathExists(path); err != nil {
		return "", fmt.Errorf("Error checking if path exists: %v", err)
	} else if pathExists {
		return path, nil
	}

	return "", nil
}
