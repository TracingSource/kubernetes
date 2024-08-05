package util

import (
	"path/filepath"

	libcontainercgroups "github.com/opencontainers/runc/libcontainer/cgroups"
	libcontainerutils "github.com/opencontainers/runc/libcontainer/utils"
)

// GetPids gets pids of the desired cgroup
// Forked from opencontainers/runc/libcontainer/cgroup/fs.Manager.GetPids()
func GetPids(cgroupPath string) ([]int, error) {
	dir, err := getCgroupPath(cgroupPath)
	if err != nil {
		return nil, err
	}
	return libcontainercgroups.GetPids(dir)
}

// getCgroupPath gets the file path to the "devices" subsystem of the desired cgroup.
// cgroupPath is the path in the cgroup hierarchy.
func getCgroupPath(cgroupPath string) (string, error) {
	cgroupPath = libcontainerutils.CleanPath(cgroupPath)

	mnt, root, err := libcontainercgroups.FindCgroupMountpointAndRoot(cgroupPath, "devices")
	// If we didn't mount the subsystem, there is no point we make the path.
	if err != nil {
		return "", err
	}

	// If the cgroup name/path is absolute do not look relative to the cgroup
	// of the init process.
	if filepath.IsAbs(cgroupPath) {
		// Sometimes subsystems can be mounted together as 'cpu,cpuacct'.
		return filepath.Join(root, mnt, cgroupPath), nil
	}

	parentPath, err := getCgroupParentPath(mnt, root)
	if err != nil {
		return "", err
	}

	return filepath.Join(parentPath, cgroupPath), nil
}

// getCgroupParentPath gets the parent filepath to this cgroup,
// for resolving relative cgroup paths.
func getCgroupParentPath(mountpoint, root string) (string, error) {
	// Use GetThisCgroupDir instead of GetInitCgroupDir, because the creating
	// process could in container and shared pid namespace with host, and
	// /proc/1/cgroup could point to whole other world of cgroups.
	initPath, err := libcontainercgroups.GetOwnCgroup("devices")
	if err != nil {
		return "", err
	}
	// This is needed for nested containers, because in /proc/self/cgroup we
	// see paths from host, which don't exist in container.
	relDir, err := filepath.Rel(root, initPath)
	if err != nil {
		return "", err
	}
	return filepath.Join(mountpoint, relDir), nil
}
