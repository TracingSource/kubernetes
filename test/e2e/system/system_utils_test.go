package system

import (
	"testing"

	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestIsMasterNode(t *testing.T) {
	testCases := []struct {
		input  string
		result bool
	}{
		{"foo-master", true},
		{"foo-master-", false},
		{"foo-master-a", false},
		{"foo-master-ab", false},
		{"foo-master-abc", true},
		{"foo-master-abdc", false},
		{"foo-bar", false},
	}

	for _, tc := range testCases {
		node := v1.Node{ObjectMeta: metav1.ObjectMeta{Name: tc.input}}
		res := DeprecatedMightBeMasterNode(node.Name)
		if res != tc.result {
			t.Errorf("case \"%s\": expected %t, got %t", tc.input, tc.result, res)
		}
	}
}
