package app

import (
	"testing"
)

func TestGetServiceIPAndRanges(t *testing.T) {
	tests := []struct {
		body                    string
		apiServerServiceIP      string
		primaryServiceIPRange   string
		secondaryServiceIPRange string
		expectedError           bool
	}{
		{"", "10.0.0.1", "10.0.0.0/24", "<nil>", false},
		{"192.0.2.1/24", "192.0.2.1", "192.0.2.0/24", "<nil>", false},
		{"192.0.2.1/24,192.168.128.0/17", "192.0.2.1", "192.0.2.0/24", "192.168.128.0/17", false},
		{"192.0.2.1/30,192.168.128.0/17", "<nil>", "<nil>", "<nil>", true},
	}

	for _, test := range tests {
		apiServerServiceIP, primaryServiceIPRange, secondaryServiceIPRange, err := getServiceIPAndRanges(test.body)

		if apiServerServiceIP.String() != test.apiServerServiceIP {
			t.Errorf("expected apiServerServiceIP: %s, got: %s", test.apiServerServiceIP, apiServerServiceIP.String())
		}

		if primaryServiceIPRange.String() != test.primaryServiceIPRange {
			t.Errorf("expected primaryServiceIPRange: %s, got: %s", test.primaryServiceIPRange, primaryServiceIPRange.String())
		}

		if secondaryServiceIPRange.String() != test.secondaryServiceIPRange {
			t.Errorf("expected secondaryServiceIPRange: %s, got: %s", test.secondaryServiceIPRange, secondaryServiceIPRange.String())
		}

		if (err == nil) == test.expectedError {
			t.Errorf("expected err to be: %t, but it was %t", test.expectedError, !test.expectedError)
		}
	}
}
