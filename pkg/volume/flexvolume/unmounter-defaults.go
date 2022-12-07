package flexvolume

import (
	"k8s.io/klog"
	"k8s.io/utils/mount"
)

type unmounterDefaults flexVolumeUnmounter

func (f *unmounterDefaults) TearDownAt(dir string) error {
	klog.Warning(logPrefix(f.plugin), "using default TearDownAt for ", dir)
	return mount.CleanupMountPoint(dir, f.mounter, false)
}
