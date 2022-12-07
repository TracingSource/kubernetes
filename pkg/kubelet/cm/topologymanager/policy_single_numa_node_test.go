package topologymanager

import (
	"testing"
)

func TestPolicySingleNumaNodeCanAdmitPodResult(t *testing.T) {
	tcases := []struct {
		name     string
		hint     TopologyHint
		expected bool
	}{
		{
			name:     "Preferred is set to false in topology hints",
			hint:     TopologyHint{nil, false},
			expected: false,
		},
	}

	for _, tc := range tcases {
		numaNodes := []int{0, 1}
		policy := NewSingleNumaNodePolicy(numaNodes)
		result := policy.(*singleNumaNodePolicy).canAdmitPodResult(&tc.hint)

		if result.Admit != tc.expected {
			t.Errorf("Expected Admit field in result to be %t, got %t", tc.expected, result.Admit)
		}

		if tc.expected == false {
			if len(result.Reason) == 0 {
				t.Errorf("Expected Reason field to be not empty")
			}
			if len(result.Message) == 0 {
				t.Errorf("Expected Message field to be not empty")
			}
		}
	}
}

func TestSingleNumaNodePolicyMerge(t *testing.T) {
	numaNodes := []int{0, 1}
	policy := NewSingleNumaNodePolicy(numaNodes)

	tcases := commonPolicyMergeTestCases(numaNodes)
	tcases = append(tcases, policy.(*singleNumaNodePolicy).mergeTestCases(numaNodes)...)

	testPolicyMerge(policy, tcases, t)
}
