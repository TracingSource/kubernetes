// +build !providerless
// +build linux

package awsebs

import (
	"k8s.io/legacy-cloud-providers/aws"
)

func (attacher *awsElasticBlockStoreAttacher) getDevicePath(volumeID, partition, devicePath string) (string, error) {
	devicePaths := getDiskByIDPaths(aws.KubernetesVolumeID(volumeID), partition, devicePath)
	return verifyDevicePath(devicePaths)
}
