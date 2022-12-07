package cpumanager

import (
	"testing"

	"k8s.io/api/core/v1"
	apimachinery "k8s.io/apimachinery/pkg/types"
)

func TestContainerMap(t *testing.T) {
	testCases := []struct {
		podUID         string
		containerNames []string
		containerIDs   []string
	}{
		{
			"fakePodUID",
			[]string{"fakeContainerName-1", "fakeContainerName-2"},
			[]string{"fakeContainerID-1", "fakeContainerName-2"},
		},
	}

	for _, tc := range testCases {
		pod := v1.Pod{}
		pod.UID = apimachinery.UID(tc.podUID)

		// Build a new containerMap from the testCases, checking proper
		// addition, retrieval along the way.
		cm := newContainerMap()
		for i := range tc.containerNames {
			container := v1.Container{Name: tc.containerNames[i]}

			cm.Add(&pod, &container, tc.containerIDs[i])
			containerID, err := cm.Get(&pod, &container)
			if err != nil {
				t.Errorf("error adding and retrieving container: %v", err)
			}
			if containerID != tc.containerIDs[i] {
				t.Errorf("mismatched containerIDs %v, %v", containerID, tc.containerIDs[i])
			}
		}

		// Remove all entries from the containerMap, checking proper removal of
		// each along the way.
		for i := range tc.containerNames {
			container := v1.Container{Name: tc.containerNames[i]}
			cm.Remove(tc.containerIDs[i])
			containerID, err := cm.Get(&pod, &container)
			if err == nil {
				t.Errorf("unexpected retrieval of containerID after removal: %v", containerID)
			}
		}

		// Verify containerMap now empty.
		if len(cm) != 0 {
			t.Errorf("unexpected entries still in containerMap: %v", cm)
		}

	}
}
