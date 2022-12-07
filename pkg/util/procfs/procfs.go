package procfs

type ProcFSInterface interface {
	// GetFullContainerName gets the container name given the root process id of the container.
	GetFullContainerName(pid int) (string, error)
}
