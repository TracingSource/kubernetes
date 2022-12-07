// +build linux

package phases

import (
	"io/ioutil"
	"strings"
	"syscall"

	"k8s.io/klog"
)

// unmountKubeletDirectory unmounts all paths that contain KubeletRunDirectory
func unmountKubeletDirectory(absoluteKubeletRunDirectory string) error {
	raw, err := ioutil.ReadFile("/proc/mounts")
	if err != nil {
		return err
	}
	mounts := strings.Split(string(raw), "\n")
	for _, mount := range mounts {
		m := strings.Split(mount, " ")
		if len(m) < 2 || !strings.HasPrefix(m[1], absoluteKubeletRunDirectory) {
			continue
		}
		if err := syscall.Unmount(m[1], 0); err != nil {
			klog.Warningf("[reset] Failed to unmount mounted directory in %s: %s", absoluteKubeletRunDirectory, m[1])
		}
	}
	return nil
}
