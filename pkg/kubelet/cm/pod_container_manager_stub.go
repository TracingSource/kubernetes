package cm

import (
	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
)

type podContainerManagerStub struct {
}

var _ PodContainerManager = &podContainerManagerStub{}

func (m *podContainerManagerStub) Exists(_ *v1.Pod) bool {
	return true
}

func (m *podContainerManagerStub) EnsureExists(_ *v1.Pod) error {
	return nil
}

func (m *podContainerManagerStub) GetPodContainerName(_ *v1.Pod) (CgroupName, string) {
	return nil, ""
}

func (m *podContainerManagerStub) Destroy(_ CgroupName) error {
	return nil
}

func (m *podContainerManagerStub) ReduceCPULimits(_ CgroupName) error {
	return nil
}

func (m *podContainerManagerStub) GetAllPodsFromCgroups() (map[types.UID]CgroupName, error) {
	return nil, nil
}

func (m *podContainerManagerStub) IsPodCgroup(cgroupfs string) (bool, types.UID) {
	return false, types.UID("")
}
