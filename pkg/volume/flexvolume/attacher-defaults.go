package flexvolume

import (
	"time"

	"k8s.io/klog"
	"k8s.io/utils/mount"

	"k8s.io/apimachinery/pkg/types"
	"k8s.io/kubernetes/pkg/volume"
)

type attacherDefaults flexVolumeAttacher

// Attach is part of the volume.Attacher interface
func (a *attacherDefaults) Attach(spec *volume.Spec, hostName types.NodeName) (string, error) {
	klog.Warning(logPrefix(a.plugin.flexVolumePlugin), "using default Attach for volume ", spec.Name(), ", host ", hostName)
	return "", nil
}

// WaitForAttach is part of the volume.Attacher interface
func (a *attacherDefaults) WaitForAttach(spec *volume.Spec, devicePath string, timeout time.Duration) (string, error) {
	klog.Warning(logPrefix(a.plugin.flexVolumePlugin), "using default WaitForAttach for volume ", spec.Name(), ", device ", devicePath)
	return devicePath, nil
}

// GetDeviceMountPath is part of the volume.Attacher interface
func (a *attacherDefaults) GetDeviceMountPath(spec *volume.Spec, mountsDir string) (string, error) {
	return a.plugin.getDeviceMountPath(spec)
}

// MountDevice is part of the volume.Attacher interface
func (a *attacherDefaults) MountDevice(spec *volume.Spec, devicePath string, deviceMountPath string, mounter mount.Interface) error {
	klog.Warning(logPrefix(a.plugin.flexVolumePlugin), "using default MountDevice for volume ", spec.Name(), ", device ", devicePath, ", deviceMountPath ", deviceMountPath)

	volSourceFSType, err := getFSType(spec)
	if err != nil {
		return err
	}

	readOnly, err := getReadOnly(spec)
	if err != nil {
		return err
	}

	options := make([]string, 0)

	if readOnly {
		options = append(options, "ro")
	} else {
		options = append(options, "rw")
	}

	diskMounter := &mount.SafeFormatAndMount{Interface: mounter, Exec: a.plugin.host.GetExec(a.plugin.GetPluginName())}

	return diskMounter.FormatAndMount(devicePath, deviceMountPath, volSourceFSType, options)
}
