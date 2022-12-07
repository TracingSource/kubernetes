package cpumanager

import (
	"k8s.io/api/core/v1"
	"k8s.io/klog"
	"k8s.io/kubernetes/pkg/kubelet/cm/cpumanager/state"
	"k8s.io/kubernetes/pkg/kubelet/cm/topologymanager"
	"k8s.io/kubernetes/pkg/kubelet/config"
	"k8s.io/kubernetes/pkg/kubelet/status"
)

type fakeManager struct {
	state state.State
}

func (m *fakeManager) Start(activePods ActivePodsFunc, sourcesReady config.SourcesReady, podStatusProvider status.PodStatusProvider, containerRuntime runtimeService) {
	klog.Info("[fake cpumanager] Start()")
}

func (m *fakeManager) Policy() Policy {
	klog.Info("[fake cpumanager] Policy()")
	return NewNonePolicy()
}

func (m *fakeManager) AddContainer(pod *v1.Pod, container *v1.Container, containerID string) error {
	klog.Infof("[fake cpumanager] AddContainer (pod: %s, container: %s, container id: %s)", pod.Name, container.Name, containerID)
	return nil
}

func (m *fakeManager) RemoveContainer(containerID string) error {
	klog.Infof("[fake cpumanager] RemoveContainer (container id: %s)", containerID)
	return nil
}

func (m *fakeManager) GetTopologyHints(pod v1.Pod, container v1.Container) map[string][]topologymanager.TopologyHint {
	klog.Infof("[fake cpumanager] Get Topology Hints")
	return map[string][]topologymanager.TopologyHint{}
}

func (m *fakeManager) State() state.Reader {
	return m.state
}

// NewFakeManager creates empty/fake cpu manager
func NewFakeManager() Manager {
	return &fakeManager{
		state: state.NewMemoryState(),
	}
}
