package dockershim

import (
	"reflect"
	"testing"
)

func TestModifySecurityOptions(t *testing.T) {
	testCases := []struct {
		name     string
		config   []string
		optName  string
		optVal   string
		expected []string
	}{
		{
			name:     "Empty val",
			config:   []string{"a:b", "c:d"},
			optName:  "optA",
			optVal:   "",
			expected: []string{"a:b", "c:d"},
		},
		{
			name:     "Valid",
			config:   []string{"a:b", "c:d"},
			optName:  "e",
			optVal:   "f",
			expected: []string{"a:b", "c:d", "e:f"},
		},
	}

	for _, tc := range testCases {
		actual := modifySecurityOption(tc.config, tc.optName, tc.optVal)
		if !reflect.DeepEqual(tc.expected, actual) {
			t.Errorf("Failed to apply options correctly for tc: %s.  Expected: %v but got %v", tc.name, tc.expected, actual)
		}
	}
}
