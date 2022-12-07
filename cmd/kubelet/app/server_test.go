package app

import (
	"testing"
)

func TestValueOfAllocatableResources(t *testing.T) {
	testCases := []struct {
		kubeReserved   map[string]string
		systemReserved map[string]string
		errorExpected  bool
		name           string
	}{
		{
			kubeReserved:   map[string]string{"cpu": "200m", "memory": "-150G", "ephemeral-storage": "10Gi"},
			systemReserved: map[string]string{"cpu": "200m", "memory": "15Ki"},
			errorExpected:  true,
			name:           "negative quantity value",
		},
		{
			kubeReserved:   map[string]string{"cpu": "200m", "memory": "150Gi", "ephemeral-storage": "10Gi"},
			systemReserved: map[string]string{"cpu": "200m", "memory": "15Ky"},
			errorExpected:  true,
			name:           "invalid quantity unit",
		},
		{
			kubeReserved:   map[string]string{"cpu": "200m", "memory": "15G", "ephemeral-storage": "10Gi"},
			systemReserved: map[string]string{"cpu": "200m", "memory": "15Ki"},
			errorExpected:  false,
			name:           "Valid resource quantity",
		},
	}

	for _, test := range testCases {
		_, err1 := parseResourceList(test.kubeReserved)
		_, err2 := parseResourceList(test.systemReserved)
		if test.errorExpected {
			if err1 == nil && err2 == nil {
				t.Errorf("%s: error expected", test.name)
			}
		} else {
			if err1 != nil || err2 != nil {
				t.Errorf("%s: unexpected error: %v, %v", test.name, err1, err2)
			}
		}
	}
}
