package validation

import (
	api "k8s.io/kubernetes/pkg/apis/core"
	internalapi "k8s.io/kubernetes/plugin/pkg/admission/podtolerationrestriction/apis/podtolerationrestriction"
	"testing"
)

func TestValidateConfiguration(t *testing.T) {

	tests := []struct {
		config     internalapi.Configuration
		testName   string
		testStatus bool
	}{
		{
			config: internalapi.Configuration{
				Default: []api.Toleration{
					{Key: "foo", Operator: "Exists", Value: "", Effect: "NoExecute", TolerationSeconds: &[]int64{60}[0]},
					{Key: "foo", Operator: "Equal", Value: "bar", Effect: "NoExecute", TolerationSeconds: &[]int64{60}[0]},
					{Key: "foo", Operator: "Equal", Value: "bar", Effect: "NoSchedule"},
					{Operator: "Exists", Effect: "NoSchedule"},
				},
				Whitelist: []api.Toleration{
					{Key: "foo", Value: "bar", Effect: "NoSchedule"},
					{Key: "foo", Operator: "Equal", Value: "bar"},
				},
			},
			testName:   "Valid cases",
			testStatus: true,
		},
		{
			config: internalapi.Configuration{
				Whitelist: []api.Toleration{{Key: "foo", Operator: "Exists", Value: "bar", Effect: "NoSchedule"}},
			},
			testName:   "Invalid case",
			testStatus: false,
		},
		{
			config: internalapi.Configuration{
				Default: []api.Toleration{{Operator: "Equal", Value: "bar", Effect: "NoSchedule"}},
			},
			testName:   "Invalid case",
			testStatus: false,
		},
	}

	for i := range tests {
		errs := ValidateConfiguration(&tests[i].config)
		if tests[i].testStatus && errs != nil {
			t.Errorf("Test: %s, expected success: %v", tests[i].testName, errs)
		}
		if !tests[i].testStatus && errs == nil {
			t.Errorf("Test: %s, expected errors: %v", tests[i].testName, errs)
		}
	}
}
