// +build !providerless
// +build !linux,!windows

package awsebs

import "errors"

func (attacher *awsElasticBlockStoreAttacher) getDevicePath(volumeID, partition, devicePath string) (string, error) {
	return "", errors.New("unsupported")
}
