// +build linux

package e2e_node

import (
	"path/filepath"
	"syscall"

	"k8s.io/utils/mount"
)

func detectMountpoint(m mount.Interface, path string) string {
	path, err := filepath.Abs(path)
	if err == nil {
		path, err = filepath.EvalSymlinks(path)
	}
	if err != nil {
		return ""
	}
	for path != "" && path != "/" {
		isNotMount, err := m.IsLikelyNotMountPoint(path)
		if err != nil {
			return ""
		}
		if !isNotMount {
			return path
		}
		path = filepath.Dir(path)
	}
	return "/"
}

const (
	xfsMagic = 0x58465342
)

// XFS over-allocates and then eventually removes that excess allocation.
// That can lead to a file growing beyond its eventual size, causing
// an unnecessary eviction:
//
// % ls -ls
// total 32704
// 32704 -rw-r--r-- 1 rkrawitz rkrawitz 20971520 Jan 15 13:16 foo.bin
//
// This issue can be hit regardless of the means used to count storage.
// It is not present in ext4fs.
func isXfs(dir string) bool {
	mountpoint := detectMountpoint(mount.New(""), dir)
	if mountpoint == "" {
		return false
	}
	var buf syscall.Statfs_t
	err := syscall.Statfs(mountpoint, &buf)
	if err != nil {
		return false
	}
	return buf.Type == xfsMagic
}
