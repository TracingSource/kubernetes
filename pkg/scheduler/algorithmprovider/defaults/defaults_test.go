package defaults

import (
	"testing"

	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/kubernetes/pkg/scheduler/algorithm/predicates"
	"k8s.io/kubernetes/pkg/scheduler/algorithm/priorities"
)

func TestCopyAndReplace(t *testing.T) {
	testCases := []struct {
		set         sets.String
		replaceWhat string
		replaceWith string
		expected    sets.String
	}{
		{
			set:         sets.String{"A": sets.Empty{}, "B": sets.Empty{}},
			replaceWhat: "A",
			replaceWith: "C",
			expected:    sets.String{"B": sets.Empty{}, "C": sets.Empty{}},
		},
		{
			set:         sets.String{"A": sets.Empty{}, "B": sets.Empty{}},
			replaceWhat: "D",
			replaceWith: "C",
			expected:    sets.String{"A": sets.Empty{}, "B": sets.Empty{}},
		},
	}
	for _, testCase := range testCases {
		result := copyAndReplace(testCase.set, testCase.replaceWhat, testCase.replaceWith)
		if !result.Equal(testCase.expected) {
			t.Errorf("expected %v got %v", testCase.expected, result)
		}
	}
}

func TestDefaultPriorities(t *testing.T) {
	result := sets.NewString(
		priorities.SelectorSpreadPriority,
		priorities.InterPodAffinityPriority,
		priorities.LeastRequestedPriority,
		priorities.BalancedResourceAllocation,
		priorities.NodePreferAvoidPodsPriority,
		priorities.NodeAffinityPriority,
		priorities.TaintTolerationPriority,
		priorities.ImageLocalityPriority,
	)
	if expected := defaultPriorities(); !result.Equal(expected) {
		t.Errorf("expected %v got %v", expected, result)
	}
}

func TestDefaultPredicates(t *testing.T) {
	result := sets.NewString(
		predicates.NoVolumeZoneConflictPred,
		predicates.MaxEBSVolumeCountPred,
		predicates.MaxGCEPDVolumeCountPred,
		predicates.MaxAzureDiskVolumeCountPred,
		predicates.MaxCSIVolumeCountPred,
		predicates.MatchInterPodAffinityPred,
		predicates.NoDiskConflictPred,
		predicates.GeneralPred,
		predicates.PodToleratesNodeTaintsPred,
		predicates.CheckVolumeBindingPred,
		predicates.CheckNodeUnschedulablePred,
	)

	if expected := defaultPredicates(); !result.Equal(expected) {
		t.Errorf("expected %v got %v", expected, result)
	}
}
