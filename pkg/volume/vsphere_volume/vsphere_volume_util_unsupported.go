// +build !providerless
// +build !linux,!windows

package vsphere_volume

import "errors"

func verifyDevicePath(path string) (string, error) {
	return "", errors.New("unsupported")
}
