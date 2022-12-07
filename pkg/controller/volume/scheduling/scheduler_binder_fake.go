package scheduling

import "k8s.io/api/core/v1"

// FakeVolumeBinderConfig holds configurations for fake volume binder.
type FakeVolumeBinderConfig struct {
	AllBound             bool
	FindUnboundSatsified bool
	FindBoundSatsified   bool
	FindErr              error
	AssumeErr            error
	BindErr              error
}

// NewFakeVolumeBinder sets up all the caches needed for the scheduler to make
// topology-aware volume binding decisions.
func NewFakeVolumeBinder(config *FakeVolumeBinderConfig) *FakeVolumeBinder {
	return &FakeVolumeBinder{
		config: config,
	}
}

// FakeVolumeBinder represents a fake volume binder for testing.
type FakeVolumeBinder struct {
	config       *FakeVolumeBinderConfig
	AssumeCalled bool
	BindCalled   bool
}

// FindPodVolumes implements SchedulerVolumeBinder.FindPodVolumes.
func (b *FakeVolumeBinder) FindPodVolumes(pod *v1.Pod, node *v1.Node) (unboundVolumesSatisfied, boundVolumesSatsified bool, err error) {
	return b.config.FindUnboundSatsified, b.config.FindBoundSatsified, b.config.FindErr
}

// AssumePodVolumes implements SchedulerVolumeBinder.AssumePodVolumes.
func (b *FakeVolumeBinder) AssumePodVolumes(assumedPod *v1.Pod, nodeName string) (bool, error) {
	b.AssumeCalled = true
	return b.config.AllBound, b.config.AssumeErr
}

// BindPodVolumes implements SchedulerVolumeBinder.BindPodVolumes.
func (b *FakeVolumeBinder) BindPodVolumes(assumedPod *v1.Pod) error {
	b.BindCalled = true
	return b.config.BindErr
}

// GetBindingsCache implements SchedulerVolumeBinder.GetBindingsCache.
func (b *FakeVolumeBinder) GetBindingsCache() PodBindingCache {
	return nil
}
