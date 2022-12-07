// +build !providerless
// +build !linux,!windows

package azure_dd

import "k8s.io/utils/exec"

func scsiHostRescan(io ioHandler, exec exec.Interface) {
}

func findDiskByLun(lun int, io ioHandler, exec exec.Interface) (string, error) {
	return "", nil
}
