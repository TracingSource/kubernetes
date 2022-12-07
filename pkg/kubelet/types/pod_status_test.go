package types

import (
	"k8s.io/api/core/v1"
	"testing"
)

func TestPodConditionByKubelet(t *testing.T) {
	trueCases := []v1.PodConditionType{
		v1.PodScheduled,
		v1.PodReady,
		v1.PodInitialized,
		v1.PodReasonUnschedulable,
	}

	for _, tc := range trueCases {
		if !PodConditionByKubelet(tc) {
			t.Errorf("Expect %q to be condition owned by kubelet.", tc)
		}
	}

	falseCases := []v1.PodConditionType{
		v1.PodConditionType("abcd"),
	}

	for _, tc := range falseCases {
		if PodConditionByKubelet(tc) {
			t.Errorf("Expect %q NOT to be condition owned by kubelet.", tc)
		}
	}
}
