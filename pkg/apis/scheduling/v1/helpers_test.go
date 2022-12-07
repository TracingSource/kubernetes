package v1

import (
	"testing"

	v1 "k8s.io/api/scheduling/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/kubernetes/pkg/apis/scheduling"
)

func TestIsKnownSystemPriorityClass(t *testing.T) {
	tests := []struct {
		name     string
		pc       *v1.PriorityClass
		expected bool
	}{
		{
			name:     "system priority class",
			pc:       SystemPriorityClasses()[0],
			expected: true,
		},
		{
			name: "non-system priority class",
			pc: &v1.PriorityClass{
				ObjectMeta: metav1.ObjectMeta{
					Name: scheduling.SystemNodeCritical,
				},
				Value:       scheduling.SystemCriticalPriority, // This is the value of system cluster critical
				Description: "Used for system critical pods that must not be moved from their current node.",
			},
			expected: false,
		},
	}

	for _, test := range tests {
		if is, err := IsKnownSystemPriorityClass(test.pc.Name, test.pc.Value, test.pc.GlobalDefault); test.expected != is {
			t.Errorf("Test [%v]: Expected %v, but got %v. Error: %v", test.name, test.expected, is, err)
		}
	}
}
