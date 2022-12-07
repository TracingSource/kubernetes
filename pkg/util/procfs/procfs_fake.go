package procfs

type FakeProcFS struct{}

func NewFakeProcFS() ProcFSInterface {
	return &FakeProcFS{}
}

// GetFullContainerName gets the container name given the root process id of the container.
// E.g. if the devices cgroup for the container is stored in /sys/fs/cgroup/devices/docker/nginx,
// return docker/nginx. Assumes that the process is part of exactly one cgroup hierarchy.
func (fakePfs *FakeProcFS) GetFullContainerName(pid int) (string, error) {
	return "", nil
}
