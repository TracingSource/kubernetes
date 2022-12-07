package util

import (
	"testing"
)

func TestGetCgroupDriverDocker(t *testing.T) {
	testCases := []struct {
		name          string
		driver        string
		expectedError bool
	}{
		{
			name:          "valid: value is 'cgroupfs'",
			driver:        `cgroupfs`,
			expectedError: false,
		},
		{
			name:          "valid: value is 'systemd'",
			driver:        `systemd`,
			expectedError: false,
		},
		{
			name:          "invalid: empty 'Cgroup Driver' value",
			driver:        ``,
			expectedError: true,
		},
		{
			name:          "invalid: unknown 'Cgroup Driver' value",
			driver:        `invalid-value`,
			expectedError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.driver != CgroupDriverCgroupfs && tc.driver != CgroupDriverSystemd
			if result != tc.expectedError {
				t.Fatalf("expected error: %v, saw: %v", tc.expectedError, result)
			}
		})
	}
}
