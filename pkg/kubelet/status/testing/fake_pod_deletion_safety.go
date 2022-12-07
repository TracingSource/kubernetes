package testing

import "k8s.io/api/core/v1"

// FakePodDeletionSafetyProvider is a fake PodDeletionSafetyProvider for test.
type FakePodDeletionSafetyProvider struct{}

// PodResourcesAreReclaimed implements PodDeletionSafetyProvider.
// Always reports that all pod resources are reclaimed.
func (f *FakePodDeletionSafetyProvider) PodResourcesAreReclaimed(pod *v1.Pod, status v1.PodStatus) bool {
	return true
}
