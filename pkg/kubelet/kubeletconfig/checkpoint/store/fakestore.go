package store

import (
	"fmt"
	"time"

	kubeletconfig "k8s.io/kubernetes/pkg/kubelet/apis/config"
	"k8s.io/kubernetes/pkg/kubelet/kubeletconfig/checkpoint"
)

// so far only implements Assigned(), LastKnownGood(), SetAssigned(), and SetLastKnownGood()
type fakeStore struct {
	assigned      checkpoint.RemoteConfigSource
	lastKnownGood checkpoint.RemoteConfigSource
}

var _ Store = (*fakeStore)(nil)

// NewFakeStore constructs a fake Store
func NewFakeStore() Store {
	return &fakeStore{}
}

func (s *fakeStore) Initialize() error {
	return fmt.Errorf("Initialize method not supported")
}

func (s *fakeStore) Exists(source checkpoint.RemoteConfigSource) (bool, error) {
	return false, fmt.Errorf("Exists method not supported")
}

func (s *fakeStore) Save(c checkpoint.Payload) error {
	return fmt.Errorf("Save method not supported")
}

func (s *fakeStore) Load(source checkpoint.RemoteConfigSource) (*kubeletconfig.KubeletConfiguration, error) {
	return nil, fmt.Errorf("Load method not supported")
}

func (s *fakeStore) AssignedModified() (time.Time, error) {
	return time.Time{}, fmt.Errorf("AssignedModified method not supported")
}

func (s *fakeStore) Assigned() (checkpoint.RemoteConfigSource, error) {
	return s.assigned, nil
}

func (s *fakeStore) LastKnownGood() (checkpoint.RemoteConfigSource, error) {
	return s.lastKnownGood, nil
}

func (s *fakeStore) SetAssigned(source checkpoint.RemoteConfigSource) error {
	s.assigned = source
	return nil
}

func (s *fakeStore) SetLastKnownGood(source checkpoint.RemoteConfigSource) error {
	s.lastKnownGood = source
	return nil
}

func (s *fakeStore) Reset() (bool, error) {
	return false, fmt.Errorf("Reset method not supported")
}
