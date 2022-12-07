// +build !linux

package phases

import (
	"k8s.io/klog"
)

// unmountKubeletDirectory is a NOOP on all but linux.
func unmountKubeletDirectory(absoluteKubeletRunDirectory string) error {
	klog.Warning("Cannot unmount filesystems on current OS, all mounted file systems will need to be manually unmounted")
	return nil
}
