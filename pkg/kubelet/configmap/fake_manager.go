package configmap

import v1 "k8s.io/api/core/v1"

// fakeManager implements Manager interface for testing purposes.
// simple operations to apiserver.
type fakeManager struct {
}

// NewFakeManager creates empty/fake ConfigMap manager
func NewFakeManager() Manager {
	return &fakeManager{}
}

func (s *fakeManager) GetConfigMap(namespace, name string) (*v1.ConfigMap, error) {
	return nil, nil
}

func (s *fakeManager) RegisterPod(pod *v1.Pod) {
}

func (s *fakeManager) UnregisterPod(pod *v1.Pod) {
}
