package flexvolume

import (
	"time"

	"k8s.io/klog"
	"k8s.io/utils/mount"

	"k8s.io/apimachinery/pkg/types"
)

type detacherDefaults flexVolumeDetacher

// Detach is part of the volume.Detacher interface.
func (d *detacherDefaults) Detach(volumeName string, hostName types.NodeName) error {
	klog.Warning(logPrefix(d.plugin.flexVolumePlugin), "using default Detach for volume ", volumeName, ", host ", hostName)
	return nil
}

// WaitForDetach is part of the volume.Detacher interface.
func (d *detacherDefaults) WaitForDetach(devicePath string, timeout time.Duration) error {
	klog.Warning(logPrefix(d.plugin.flexVolumePlugin), "using default WaitForDetach for device ", devicePath)
	return nil
}

// UnmountDevice is part of the volume.Detacher interface.
func (d *detacherDefaults) UnmountDevice(deviceMountPath string) error {
	klog.Warning(logPrefix(d.plugin.flexVolumePlugin), "using default UnmountDevice for device mount path ", deviceMountPath)
	return mount.CleanupMountPoint(deviceMountPath, d.plugin.host.GetMounter(d.plugin.GetPluginName()), false)
}
