package cpumanager

import (
	"k8s.io/api/core/v1"
	"k8s.io/klog"
	"k8s.io/kubernetes/pkg/kubelet/cm/cpumanager/state"
	"k8s.io/kubernetes/pkg/kubelet/cm/topologymanager"
)

type nonePolicy struct{}

var _ Policy = &nonePolicy{}

// PolicyNone name of none policy
const PolicyNone policyName = "none"

// NewNonePolicy returns a cupset manager policy that does nothing
func NewNonePolicy() Policy {
	return &nonePolicy{}
}

func (p *nonePolicy) Name() string {
	return string(PolicyNone)
}

func (p *nonePolicy) Start(s state.State) {
	klog.Info("[cpumanager] none policy: Start")
}

func (p *nonePolicy) AddContainer(s state.State, pod *v1.Pod, container *v1.Container, containerID string) error {
	return nil
}

func (p *nonePolicy) RemoveContainer(s state.State, containerID string) error {
	return nil
}

func (p *nonePolicy) GetTopologyHints(s state.State, pod v1.Pod, container v1.Container) map[string][]topologymanager.TopologyHint {
	return nil
}
