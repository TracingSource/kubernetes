// +build !linux

package procfs

import (
	"fmt"
	"syscall"
)

type ProcFS struct{}

func NewProcFS() ProcFSInterface {
	return &ProcFS{}
}

// GetFullContainerName gets the container name given the root process id of the container.
func (pfs *ProcFS) GetFullContainerName(pid int) (string, error) {
	return "", fmt.Errorf("GetFullContainerName is unsupported in this build")
}

// Find process(es) using a regular expression and send a specified
// signal to each process
func PKill(name string, sig syscall.Signal) error {
	return fmt.Errorf("PKill is unsupported in this build")
}

// Find process(es) with a specified name (exact match)
// and return their pid(s)
func PidOf(name string) ([]int, error) {
	return []int{}, fmt.Errorf("PidOf is unsupported in this build")
}
