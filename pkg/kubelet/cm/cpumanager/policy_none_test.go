package cpumanager

import (
	"testing"

	"k8s.io/kubernetes/pkg/kubelet/cm/cpumanager/state"
	"k8s.io/kubernetes/pkg/kubelet/cm/cpuset"
)

func TestNonePolicyName(t *testing.T) {
	policy := &nonePolicy{}

	policyName := policy.Name()
	if policyName != "none" {
		t.Errorf("NonePolicy Name() error. expected: none, returned: %v",
			policyName)
	}
}

func TestNonePolicyAdd(t *testing.T) {
	policy := &nonePolicy{}

	st := &mockState{
		assignments:   state.ContainerCPUAssignments{},
		defaultCPUSet: cpuset.NewCPUSet(1, 2, 3, 4, 5, 6, 7),
	}

	testPod := makePod("1000m", "1000m")

	container := &testPod.Spec.Containers[0]
	err := policy.AddContainer(st, testPod, container, "fakeID")
	if err != nil {
		t.Errorf("NonePolicy AddContainer() error. expected no error but got: %v", err)
	}
}

func TestNonePolicyRemove(t *testing.T) {
	policy := &nonePolicy{}

	st := &mockState{
		assignments:   state.ContainerCPUAssignments{},
		defaultCPUSet: cpuset.NewCPUSet(1, 2, 3, 4, 5, 6, 7),
	}

	err := policy.RemoveContainer(st, "fakeID")
	if err != nil {
		t.Errorf("NonePolicy RemoveContainer() error. expected no error but got %v", err)
	}
}
