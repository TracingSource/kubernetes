package flexvolume

import (
	"fmt"
	"os"

	"k8s.io/klog"
	"k8s.io/utils/mount"

	"k8s.io/apimachinery/pkg/types"
	"k8s.io/kubernetes/pkg/volume"
)

type flexVolumeDetacher struct {
	plugin *flexVolumeAttachablePlugin
}

var _ volume.Detacher = &flexVolumeDetacher{}

var _ volume.DeviceUnmounter = &flexVolumeDetacher{}

// Detach is part of the volume.Detacher interface.
func (d *flexVolumeDetacher) Detach(volumeName string, hostName types.NodeName) error {

	call := d.plugin.NewDriverCall(detachCmd)
	call.Append(volumeName)
	call.Append(string(hostName))

	_, err := call.Run()
	if isCmdNotSupportedErr(err) {
		return (*detacherDefaults)(d).Detach(volumeName, hostName)
	}
	return err
}

// UnmountDevice is part of the volume.Detacher interface.
func (d *flexVolumeDetacher) UnmountDevice(deviceMountPath string) error {

	pathExists, pathErr := mount.PathExists(deviceMountPath)
	if !pathExists {
		klog.Warningf("Warning: Unmount skipped because path does not exist: %v", deviceMountPath)
		return nil
	}
	if pathErr != nil && !mount.IsCorruptedMnt(pathErr) {
		return fmt.Errorf("Error checking path: %v", pathErr)
	}

	notmnt, err := isNotMounted(d.plugin.host.GetMounter(d.plugin.GetPluginName()), deviceMountPath)
	if err != nil {
		if mount.IsCorruptedMnt(err) {
			notmnt = false // Corrupted error is assumed to be mounted.
		} else {
			return err
		}
	}

	if notmnt {
		klog.Warningf("Warning: Path: %v already unmounted", deviceMountPath)
	} else {
		call := d.plugin.NewDriverCall(unmountDeviceCmd)
		call.Append(deviceMountPath)

		_, err := call.Run()
		if isCmdNotSupportedErr(err) {
			err = (*detacherDefaults)(d).UnmountDevice(deviceMountPath)
		}
		if err != nil {
			return err
		}
	}

	// Flexvolume driver may remove the directory. Ignore if it does.
	if pathExists, pathErr := mount.PathExists(deviceMountPath); pathErr != nil {
		return fmt.Errorf("Error checking if path exists: %v", pathErr)
	} else if !pathExists {
		return nil
	}
	return os.Remove(deviceMountPath)
}
