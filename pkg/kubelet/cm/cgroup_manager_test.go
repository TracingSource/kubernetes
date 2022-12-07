// +build linux

package cm

import (
	"reflect"
	"testing"

	"k8s.io/api/core/v1"
)

func Test(t *testing.T) {
	tests := []struct {
		input    map[string]string
		expected *map[v1.ResourceName]int64
	}{
		{
			input:    map[string]string{"memory": ""},
			expected: nil,
		},
		{
			input:    map[string]string{"memory": "a"},
			expected: nil,
		},
		{
			input:    map[string]string{"memory": "a%"},
			expected: nil,
		},
		{
			input:    map[string]string{"memory": "200%"},
			expected: nil,
		},
		{
			input: map[string]string{"memory": "0%"},
			expected: &map[v1.ResourceName]int64{
				v1.ResourceMemory: 0,
			},
		},
		{
			input: map[string]string{"memory": "100%"},
			expected: &map[v1.ResourceName]int64{
				v1.ResourceMemory: 100,
			},
		},
		{
			// need to change this when CPU is added as a supported resource
			input:    map[string]string{"memory": "100%", "cpu": "50%"},
			expected: nil,
		},
	}
	for _, test := range tests {
		actual, err := ParseQOSReserved(test.input)
		if actual != nil && test.expected == nil {
			t.Errorf("Unexpected success, input: %v, expected: %v, actual: %v, err: %v", test.input, test.expected, actual, err)
		}
		if actual == nil && test.expected != nil {
			t.Errorf("Unexpected failure, input: %v, expected: %v, actual: %v, err: %v", test.input, test.expected, actual, err)
		}
		if (actual == nil && test.expected == nil) || reflect.DeepEqual(*actual, *test.expected) {
			continue
		}
		t.Errorf("Unexpected result, input: %v, expected: %v, actual: %v, err: %v", test.input, test.expected, actual, err)
	}
}
