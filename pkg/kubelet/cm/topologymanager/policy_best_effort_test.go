package topologymanager

import (
	"testing"
)

func TestPolicyBestEffortCanAdmitPodResult(t *testing.T) {
	tcases := []struct {
		name     string
		hint     TopologyHint
		expected bool
	}{
		{
			name:     "Preferred is set to false in topology hints",
			hint:     TopologyHint{nil, false},
			expected: true,
		},
		{
			name:     "Preferred is set to true in topology hints",
			hint:     TopologyHint{nil, true},
			expected: true,
		},
	}

	for _, tc := range tcases {
		numaNodes := []int{0, 1}
		policy := NewBestEffortPolicy(numaNodes)
		result := policy.(*bestEffortPolicy).canAdmitPodResult(&tc.hint)

		if result.Admit != tc.expected {
			t.Errorf("Expected Admit field in result to be %t, got %t", tc.expected, result.Admit)
		}
	}
}

func TestBestEffortPolicyMerge(t *testing.T) {
	numaNodes := []int{0, 1}
	policy := NewBestEffortPolicy(numaNodes)

	tcases := commonPolicyMergeTestCases(numaNodes)
	tcases = append(tcases, policy.(*bestEffortPolicy).mergeTestCases(numaNodes)...)

	testPolicyMerge(policy, tcases, t)
}
