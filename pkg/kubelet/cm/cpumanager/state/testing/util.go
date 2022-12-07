package testing

import "k8s.io/kubernetes/pkg/kubelet/checkpointmanager"

var _ checkpointmanager.Checkpoint = &MockCheckpoint{}

// MockCheckpoint struct is used for mocking checkpoint values in testing
type MockCheckpoint struct {
	Content string
}

// MarshalCheckpoint returns fake content
func (mc *MockCheckpoint) MarshalCheckpoint() ([]byte, error) {
	return []byte(mc.Content), nil
}

// UnmarshalCheckpoint fakes unmarshaling
func (mc *MockCheckpoint) UnmarshalCheckpoint(blob []byte) error {
	return nil
}

// VerifyChecksum fakes verifying checksum
func (mc *MockCheckpoint) VerifyChecksum() error {
	return nil
}
