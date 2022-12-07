package state

import (
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/require"
	"k8s.io/kubernetes/pkg/kubelet/checkpointmanager"
	"k8s.io/kubernetes/pkg/kubelet/cm/cpuset"
)

const compatibilityTestingCheckpoint = "cpumanager_state_compatibility_test"

var state = &stateMemory{
	assignments: ContainerCPUAssignments{
		"container1": cpuset.NewCPUSet(4, 5, 6),
		"container2": cpuset.NewCPUSet(1, 2, 3),
	},
	defaultCPUSet: cpuset.NewCPUSet(1, 2, 3),
}

func TestFileToCheckpointCompatibility(t *testing.T) {
	statePath := path.Join(testingDir, compatibilityTestingCheckpoint)

	// ensure there is no previous state saved at testing path
	os.Remove(statePath)
	// ensure testing state is removed after testing
	defer os.Remove(statePath)

	fileState := NewFileState(statePath, "none")

	fileState.SetDefaultCPUSet(state.defaultCPUSet)
	fileState.SetCPUAssignments(state.assignments)

	restoredState, err := NewCheckpointState(testingDir, compatibilityTestingCheckpoint, "none")
	if err != nil {
		t.Fatalf("could not restore file state: %v", err)
	}

	AssertStateEqual(t, restoredState, state)
}

func TestCheckpointToFileCompatibility(t *testing.T) {
	cpm, err := checkpointmanager.NewCheckpointManager(testingDir)
	if err != nil {
		t.Fatalf("could not create testing checkpoint manager: %v", err)
	}

	// ensure there is no previous checkpoint
	cpm.RemoveCheckpoint(compatibilityTestingCheckpoint)
	// ensure testing checkpoint is removed after testing
	defer cpm.RemoveCheckpoint(compatibilityTestingCheckpoint)

	checkpointState, err := NewCheckpointState(testingDir, compatibilityTestingCheckpoint, "none")
	require.NoError(t, err)

	checkpointState.SetDefaultCPUSet(state.defaultCPUSet)
	checkpointState.SetCPUAssignments(state.assignments)

	restoredState := NewFileState(path.Join(testingDir, compatibilityTestingCheckpoint), "none")

	AssertStateEqual(t, restoredState, state)
}
