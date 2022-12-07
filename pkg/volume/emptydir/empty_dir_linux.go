// +build linux

package emptydir

import (
	"fmt"

	"golang.org/x/sys/unix"
	"k8s.io/klog"
	"k8s.io/utils/mount"

	v1 "k8s.io/api/core/v1"
)

// Defined by Linux - the type number for tmpfs mounts.
const (
	linuxTmpfsMagic     = 0x01021994
	linuxHugetlbfsMagic = 0x958458f6
)

// realMountDetector implements mountDetector in terms of syscalls.
type realMountDetector struct {
	mounter mount.Interface
}

func (m *realMountDetector) GetMountMedium(path string) (v1.StorageMedium, bool, error) {
	klog.V(5).Infof("Determining mount medium of %v", path)
	notMnt, err := m.mounter.IsLikelyNotMountPoint(path)
	if err != nil {
		return v1.StorageMediumDefault, false, fmt.Errorf("IsLikelyNotMountPoint(%q): %v", path, err)
	}
	buf := unix.Statfs_t{}
	if err := unix.Statfs(path, &buf); err != nil {
		return v1.StorageMediumDefault, false, fmt.Errorf("statfs(%q): %v", path, err)
	}

	klog.V(5).Infof("Statfs_t of %v: %+v", path, buf)
	if buf.Type == linuxTmpfsMagic {
		return v1.StorageMediumMemory, !notMnt, nil
	} else if int64(buf.Type) == linuxHugetlbfsMagic {
		return v1.StorageMediumHugePages, !notMnt, nil
	}
	return v1.StorageMediumDefault, !notMnt, nil
}
