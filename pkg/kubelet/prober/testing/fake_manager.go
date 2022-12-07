package testing

import (
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/sets"
)

// FakeManager simulates a prober.Manager for testing.
type FakeManager struct{}

// Unused methods below.

// AddPod simulates adding a Pod.
func (FakeManager) AddPod(_ *v1.Pod) {}

// RemovePod simulates removing a Pod.
func (FakeManager) RemovePod(_ *v1.Pod) {}

// CleanupPods simulates cleaning up Pods.
func (FakeManager) CleanupPods(_ map[types.UID]sets.Empty) {}

// Start simulates start syncing the probe status
func (FakeManager) Start() {}

// UpdatePodStatus simulates updating the Pod Status.
func (FakeManager) UpdatePodStatus(_ types.UID, podStatus *v1.PodStatus) {
	for i := range podStatus.ContainerStatuses {
		podStatus.ContainerStatuses[i].Ready = true
	}
}
