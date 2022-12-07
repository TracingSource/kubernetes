package fake

import (
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/labels"
	internalcache "k8s.io/kubernetes/pkg/scheduler/internal/cache"
	schedulerlisters "k8s.io/kubernetes/pkg/scheduler/listers"
	nodeinfosnapshot "k8s.io/kubernetes/pkg/scheduler/nodeinfo/snapshot"
)

// Cache is used for testing
type Cache struct {
	AssumeFunc       func(*v1.Pod)
	ForgetFunc       func(*v1.Pod)
	IsAssumedPodFunc func(*v1.Pod) bool
	GetPodFunc       func(*v1.Pod) *v1.Pod
}

// AssumePod is a fake method for testing.
func (c *Cache) AssumePod(pod *v1.Pod) error {
	c.AssumeFunc(pod)
	return nil
}

// FinishBinding is a fake method for testing.
func (c *Cache) FinishBinding(pod *v1.Pod) error { return nil }

// ForgetPod is a fake method for testing.
func (c *Cache) ForgetPod(pod *v1.Pod) error {
	c.ForgetFunc(pod)
	return nil
}

// AddPod is a fake method for testing.
func (c *Cache) AddPod(pod *v1.Pod) error { return nil }

// UpdatePod is a fake method for testing.
func (c *Cache) UpdatePod(oldPod, newPod *v1.Pod) error { return nil }

// RemovePod is a fake method for testing.
func (c *Cache) RemovePod(pod *v1.Pod) error { return nil }

// IsAssumedPod is a fake method for testing.
func (c *Cache) IsAssumedPod(pod *v1.Pod) (bool, error) {
	return c.IsAssumedPodFunc(pod), nil
}

// GetPod is a fake method for testing.
func (c *Cache) GetPod(pod *v1.Pod) (*v1.Pod, error) {
	return c.GetPodFunc(pod), nil
}

// AddNode is a fake method for testing.
func (c *Cache) AddNode(node *v1.Node) error { return nil }

// UpdateNode is a fake method for testing.
func (c *Cache) UpdateNode(oldNode, newNode *v1.Node) error { return nil }

// RemoveNode is a fake method for testing.
func (c *Cache) RemoveNode(node *v1.Node) error { return nil }

// UpdateNodeInfoSnapshot is a fake method for testing.
func (c *Cache) UpdateNodeInfoSnapshot(nodeSnapshot *nodeinfosnapshot.Snapshot) error {
	return nil
}

// List is a fake method for testing.
func (c *Cache) List(s labels.Selector) ([]*v1.Pod, error) { return nil, nil }

// FilteredList is a fake method for testing.
func (c *Cache) FilteredList(filter schedulerlisters.PodFilter, selector labels.Selector) ([]*v1.Pod, error) {
	return nil, nil
}

// Snapshot is a fake method for testing
func (c *Cache) Snapshot() *internalcache.Snapshot {
	return &internalcache.Snapshot{}
}

// GetNodeInfo is a fake method for testing.
func (c *Cache) GetNodeInfo(nodeName string) (*v1.Node, error) {
	return nil, nil
}

// ListNodes is a fake method for testing.
func (c *Cache) ListNodes() []*v1.Node {
	return nil
}
