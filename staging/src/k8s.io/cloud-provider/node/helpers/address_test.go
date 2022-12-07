package helpers

import (
	"testing"

	"k8s.io/api/core/v1"
	apiequality "k8s.io/apimachinery/pkg/api/equality"
)

func TestAddToNodeAddresses(t *testing.T) {
	testCases := []struct {
		name     string
		existing []v1.NodeAddress
		toAdd    []v1.NodeAddress
		expected []v1.NodeAddress
	}{
		{
			name:     "add no addresses to empty node addresses",
			existing: []v1.NodeAddress{},
			toAdd:    []v1.NodeAddress{},
			expected: []v1.NodeAddress{},
		},
		{
			name:     "add new node addresses to empty existing addresses",
			existing: []v1.NodeAddress{},
			toAdd: []v1.NodeAddress{
				{Type: v1.NodeExternalIP, Address: "1.1.1.1"},
				{Type: v1.NodeHostName, Address: "localhost"},
			},
			expected: []v1.NodeAddress{
				{Type: v1.NodeExternalIP, Address: "1.1.1.1"},
				{Type: v1.NodeHostName, Address: "localhost"},
			},
		},
		{
			name:     "add a duplicate node address",
			existing: []v1.NodeAddress{},
			toAdd: []v1.NodeAddress{
				{Type: v1.NodeExternalIP, Address: "1.1.1.1"},
				{Type: v1.NodeExternalIP, Address: "1.1.1.1"},
			},
			expected: []v1.NodeAddress{
				{Type: v1.NodeExternalIP, Address: "1.1.1.1"},
			},
		},
		{
			name: "add 1 new and 1 duplicate address",
			existing: []v1.NodeAddress{
				{Type: v1.NodeExternalIP, Address: "1.1.1.1"},
				{Type: v1.NodeInternalIP, Address: "10.1.1.1"},
			},
			toAdd: []v1.NodeAddress{
				{Type: v1.NodeExternalIP, Address: "1.1.1.1"},
				{Type: v1.NodeHostName, Address: "localhost"},
			},
			expected: []v1.NodeAddress{
				{Type: v1.NodeExternalIP, Address: "1.1.1.1"},
				{Type: v1.NodeInternalIP, Address: "10.1.1.1"},
				{Type: v1.NodeHostName, Address: "localhost"},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			AddToNodeAddresses(&tc.existing, tc.toAdd...)
			if !apiequality.Semantic.DeepEqual(tc.expected, tc.existing) {
				t.Errorf("expected: %v, got: %v", tc.expected, tc.existing)
			}
		})
	}
}
