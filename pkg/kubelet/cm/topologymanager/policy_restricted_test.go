package topologymanager

import (
	"testing"
)

func TestPolicyRestrictedCanAdmitPodResult(t *testing.T) {
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
		{
			name:     "Preferred is set to true in topology hints",
			hint:     TopologyHint{nil, true},
			expected: true,
		},
	}

	for _, tc := range tcases {
		numaNodes := []int{0, 1}
		policy := NewRestrictedPolicy(numaNodes)
		result := policy.(*restrictedPolicy).canAdmitPodResult(&tc.hint)

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

func TestRestrictedPolicyMerge(t *testing.T) {
	numaNodes := []int{0, 1}
	policy := NewRestrictedPolicy(numaNodes)

	tcases := commonPolicyMergeTestCases(numaNodes)
	tcases = append(tcases, policy.(*restrictedPolicy).mergeTestCases(numaNodes)...)

	testPolicyMerge(policy, tcases, t)
}
