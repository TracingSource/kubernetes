// +build !linux

package util

// GetPids gets pids of the desired cgroup
func GetPids(cgroupPath string) ([]int, error) {
	return nil, nil
}
