// +build !linux

package volumepathhandler

import (
	"fmt"

	"k8s.io/apimachinery/pkg/types"
)

// AttachFileDevice takes a path to a regular file and makes it available as an
// attached block device.
func (v VolumePathHandler) AttachFileDevice(path string) (string, error) {
	return "", fmt.Errorf("AttachFileDevice not supported for this build.")
}

// DetachFileDevice takes a path to the attached block device and
// detach it from block device.
func (v VolumePathHandler) DetachFileDevice(path string) error {
	return fmt.Errorf("DetachFileDevice not supported for this build.")
}

// GetLoopDevice returns the full path to the loop device associated with the given path.
func (v VolumePathHandler) GetLoopDevice(path string) (string, error) {
	return "", fmt.Errorf("GetLoopDevice not supported for this build.")
}

// FindGlobalMapPathUUIDFromPod finds {pod uuid} bind mount under globalMapPath
// corresponding to map path symlink, and then return global map path with pod uuid.
func (v VolumePathHandler) FindGlobalMapPathUUIDFromPod(pluginDir, mapPath string, podUID types.UID) (string, error) {
	return "", fmt.Errorf("FindGlobalMapPathUUIDFromPod not supported for this build.")
}
