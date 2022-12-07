package state

import (
	"k8s.io/kubernetes/pkg/kubelet/cm/cpuset"
)

// ContainerCPUAssignments type used in cpu manager state
type ContainerCPUAssignments map[string]cpuset.CPUSet

// Clone returns a copy of ContainerCPUAssignments
func (as ContainerCPUAssignments) Clone() ContainerCPUAssignments {
	ret := make(ContainerCPUAssignments)
	for key, val := range as {
		ret[key] = val
	}
	return ret
}

// Reader interface used to read current cpu/pod assignment state
type Reader interface {
	GetCPUSet(containerID string) (cpuset.CPUSet, bool)
	GetDefaultCPUSet() cpuset.CPUSet
	GetCPUSetOrDefault(containerID string) cpuset.CPUSet
	GetCPUAssignments() ContainerCPUAssignments
}

type writer interface {
	SetCPUSet(containerID string, cpuset cpuset.CPUSet)
	SetDefaultCPUSet(cpuset cpuset.CPUSet)
	SetCPUAssignments(ContainerCPUAssignments)
	Delete(containerID string)
	ClearState()
}

// State interface provides methods for tracking and setting cpu/pod assignment
type State interface {
	Reader
	writer
}
