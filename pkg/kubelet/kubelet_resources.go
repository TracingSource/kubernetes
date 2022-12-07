package kubelet

import (
	"fmt"

	"k8s.io/klog"

	"k8s.io/api/core/v1"
	"k8s.io/kubernetes/pkg/api/v1/resource"
)

// defaultPodLimitsForDownwardAPI copies the input pod, and optional container,
// and applies default resource limits. it returns a copy of the input pod,
// and a copy of the input container (if specified) with default limits
// applied. if a container has no limit specified, it will default the limit to
// the node allocatable.
// TODO: if/when we have pod level resources, we need to update this function
// to use those limits instead of node allocatable.
func (kl *Kubelet) defaultPodLimitsForDownwardAPI(pod *v1.Pod, container *v1.Container) (*v1.Pod, *v1.Container, error) {
	if pod == nil {
		return nil, nil, fmt.Errorf("invalid input, pod cannot be nil")
	}

	node, err := kl.getNodeAnyWay()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to find node object, expected a node")
	}
	allocatable := node.Status.Allocatable
	klog.Infof("allocatable: %v", allocatable)
	outputPod := pod.DeepCopy()
	for idx := range outputPod.Spec.Containers {
		resource.MergeContainerResourceLimits(&outputPod.Spec.Containers[idx], allocatable)
	}

	var outputContainer *v1.Container
	if container != nil {
		outputContainer = container.DeepCopy()
		resource.MergeContainerResourceLimits(outputContainer, allocatable)
	}
	return outputPod, outputContainer, nil
}
