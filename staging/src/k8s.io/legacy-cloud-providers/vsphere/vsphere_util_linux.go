// +build !providerless
// +build linux

package vsphere

import (
	"io/ioutil"
)

const UUIDPath = "/sys/class/dmi/id/product_serial"

func getRawUUID() (string, error) {
	id, err := ioutil.ReadFile(UUIDPath)
	if err != nil {
		return "", err
	}
	return string(id), nil
}
