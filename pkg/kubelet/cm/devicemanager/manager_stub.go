package devicemanager

import (
	v1 "k8s.io/api/core/v1"
	podresourcesapi "k8s.io/kubernetes/pkg/kubelet/apis/podresources/v1alpha1"
	"k8s.io/kubernetes/pkg/kubelet/cm/topologymanager"
	"k8s.io/kubernetes/pkg/kubelet/config"
	"k8s.io/kubernetes/pkg/kubelet/lifecycle"
	"k8s.io/kubernetes/pkg/kubelet/pluginmanager/cache"
	schedulernodeinfo "k8s.io/kubernetes/pkg/scheduler/nodeinfo"
)

// ManagerStub provides a simple stub implementation for the Device Manager.
type ManagerStub struct{}

// NewManagerStub creates a ManagerStub.
func NewManagerStub() (*ManagerStub, error) {
	return &ManagerStub{}, nil
}

// Start simply returns nil.
func (h *ManagerStub) Start(activePods ActivePodsFunc, sourcesReady config.SourcesReady) error {
	return nil
}

// Stop simply returns nil.
func (h *ManagerStub) Stop() error {
	return nil
}

// Allocate simply returns nil.
func (h *ManagerStub) Allocate(node *schedulernodeinfo.NodeInfo, attrs *lifecycle.PodAdmitAttributes) error {
	return nil
}

// GetDeviceRunContainerOptions simply returns nil.
func (h *ManagerStub) GetDeviceRunContainerOptions(pod *v1.Pod, container *v1.Container) (*DeviceRunContainerOptions, error) {
	return nil, nil
}

// GetCapacity simply returns nil capacity and empty removed resource list.
func (h *ManagerStub) GetCapacity() (v1.ResourceList, v1.ResourceList, []string) {
	return nil, nil, []string{}
}

// GetWatcherHandler returns plugin watcher interface
func (h *ManagerStub) GetWatcherHandler() cache.PluginHandler {
	return nil
}

// GetTopologyHints returns an empty TopologyHint map
func (h *ManagerStub) GetTopologyHints(pod v1.Pod, container v1.Container) map[string][]topologymanager.TopologyHint {
	return map[string][]topologymanager.TopologyHint{}
}

// GetDevices returns nil
func (h *ManagerStub) GetDevices(_, _ string) []*podresourcesapi.ContainerDevices {
	return nil
}

// ShouldResetExtendedResourceCapacity returns false
func (h *ManagerStub) ShouldResetExtendedResourceCapacity() bool {
	return false
}
