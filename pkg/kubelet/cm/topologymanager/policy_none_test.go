package topologymanager

import (
	"testing"
)

func TestName(t *testing.T) {
	tcases := []struct {
		name     string
		expected string
	}{
		{
			name:     "New None Policy",
			expected: "none",
		},
	}
	for _, tc := range tcases {
		policy := NewNonePolicy()
		if policy.Name() != tc.expected {
			t.Errorf("Expected Policy Name to be %s, got %s", tc.expected, policy.Name())
		}
	}
}

func TestPolicyNoneCanAdmitPodResult(t *testing.T) {
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
		policy := NewNonePolicy()
		result := policy.(*nonePolicy).canAdmitPodResult(&tc.hint)

		if result.Admit != tc.expected {
			t.Errorf("Expected Admit field in result to be %t, got %t", tc.expected, result.Admit)
		}
	}
}
