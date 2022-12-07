// +build !linux

package emptydir

import (
	"k8s.io/utils/mount"

	v1 "k8s.io/api/core/v1"
)

// realMountDetector pretends to implement mediumer.
type realMountDetector struct {
	mounter mount.Interface
}

func (m *realMountDetector) GetMountMedium(path string) (v1.StorageMedium, bool, error) {
	return v1.StorageMediumDefault, false, nil
}
