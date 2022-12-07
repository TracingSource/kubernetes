package cm

import (
	"k8s.io/api/core/v1"

	utilfeature "k8s.io/apiserver/pkg/util/feature"
	kubefeatures "k8s.io/kubernetes/pkg/features"
	"k8s.io/kubernetes/pkg/kubelet/cm/cpumanager"
	"k8s.io/kubernetes/pkg/kubelet/cm/topologymanager"
)

type InternalContainerLifecycle interface {
	PreStartContainer(pod *v1.Pod, container *v1.Container, containerID string) error
	PreStopContainer(containerID string) error
	PostStopContainer(containerID string) error
}

// Implements InternalContainerLifecycle interface.
type internalContainerLifecycleImpl struct {
	cpuManager      cpumanager.Manager
	topologyManager topologymanager.Manager
}

func (i *internalContainerLifecycleImpl) PreStartContainer(
	pod *v1.Pod, container *v1.Container, containerID string,
) error {
	if utilfeature.DefaultFeatureGate.Enabled(kubefeatures.CPUManager) {
		err := i.cpuManager.AddContainer(pod, container, containerID)
		if err != nil {
			return err
		}
	}
	if utilfeature.DefaultFeatureGate.Enabled(kubefeatures.TopologyManager) {
		err := i.topologyManager.AddContainer(pod, containerID)
		if err != nil {
			return err
		}
	}
	return nil
}

func (i *internalContainerLifecycleImpl) PreStopContainer(containerID string) error {
	if utilfeature.DefaultFeatureGate.Enabled(kubefeatures.CPUManager) {
		return i.cpuManager.RemoveContainer(containerID)
	}
	return nil
}

func (i *internalContainerLifecycleImpl) PostStopContainer(containerID string) error {
	if utilfeature.DefaultFeatureGate.Enabled(kubefeatures.CPUManager) {
		err := i.cpuManager.RemoveContainer(containerID)
		if err != nil {
			return err
		}
	}
	if utilfeature.DefaultFeatureGate.Enabled(kubefeatures.TopologyManager) {
		err := i.topologyManager.RemoveContainer(containerID)
		if err != nil {
			return err
		}
	}
	return nil
}
