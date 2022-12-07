package secret

import v1 "k8s.io/api/core/v1"

// fakeManager implements Manager interface for testing purposes.
// simple operations to apiserver.
type fakeManager struct {
}

// NewFakeManager creates empty/fake secret manager
func NewFakeManager() Manager {
	return &fakeManager{}
}

// GetSecret returns a nil secret for testing
func (s *fakeManager) GetSecret(namespace, name string) (*v1.Secret, error) {
	return nil, nil
}

// RegisterPod implements the RegisterPod method for testing purposes.
func (s *fakeManager) RegisterPod(pod *v1.Pod) {
}

// UnregisterPod implements the UnregisterPod method for testing purposes.
func (s *fakeManager) UnregisterPod(pod *v1.Pod) {
}
