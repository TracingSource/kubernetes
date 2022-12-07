package validation

import (
	"testing"

	resourcequotaapi "k8s.io/kubernetes/plugin/pkg/admission/resourcequota/apis/resourcequota"
)

func TestValidateConfiguration(t *testing.T) {
	successCases := []resourcequotaapi.Configuration{
		{
			LimitedResources: []resourcequotaapi.LimitedResource{
				{
					Resource:      "pods",
					MatchContains: []string{"requests.cpu"},
				},
			},
		},
		{
			LimitedResources: []resourcequotaapi.LimitedResource{
				{
					Resource:      "persistentvolumeclaims",
					MatchContains: []string{"requests.storage"},
				},
			},
		},
	}
	for i := range successCases {
		configuration := successCases[i]
		if errs := ValidateConfiguration(&configuration); len(errs) != 0 {
			t.Errorf("expected success: %v", errs)
		}
	}
	errorCases := map[string]resourcequotaapi.Configuration{
		"missing apiGroupResource": {LimitedResources: []resourcequotaapi.LimitedResource{
			{MatchContains: []string{"requests.cpu"}},
		}},
	}
	for k, v := range errorCases {
		if errs := ValidateConfiguration(&v); len(errs) == 0 {
			t.Errorf("expected failure for %s", k)
		}
	}
}
