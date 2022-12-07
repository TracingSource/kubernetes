package cpumanager

import (
	"k8s.io/api/core/v1"
	"k8s.io/kubernetes/pkg/kubelet/cm/cpumanager/state"
	"k8s.io/kubernetes/pkg/kubelet/cm/topologymanager"
)

// Policy implements logic for pod container to CPU assignment.
type Policy interface {
	Name() string
	Start(s state.State)
	// AddContainer call is idempotent
	AddContainer(s state.State, pod *v1.Pod, container *v1.Container, containerID string) error
	// RemoveContainer call is idempotent
	RemoveContainer(s state.State, containerID string) error
	// GetTopologyHints implements the topologymanager.HintProvider Interface
	// and is consulted to achieve NUMA aware resource alignment among this
	// and other resource controllers.
	GetTopologyHints(s state.State, pod v1.Pod, container v1.Container) map[string][]topologymanager.TopologyHint
}
