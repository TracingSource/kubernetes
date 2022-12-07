package options

import (
	"testing"
)

func TestValidateDeprecatedKubeSchedulerConfiguration(t *testing.T) {
	scenarios := map[string]struct {
		expectedToFail bool
		config         *DeprecatedOptions
	}{
		"good": {
			expectedToFail: false,
			config: &DeprecatedOptions{
				PolicyConfigFile:      "/some/file",
				UseLegacyPolicyConfig: true,
				AlgorithmProvider:     "",
			},
		},
		"bad-policy-config-file-null": {
			expectedToFail: true,
			config: &DeprecatedOptions{
				PolicyConfigFile:      "",
				UseLegacyPolicyConfig: true,
				AlgorithmProvider:     "",
			},
		},
	}

	for name, scenario := range scenarios {
		errs := scenario.config.Validate()
		if len(errs) == 0 && scenario.expectedToFail {
			t.Errorf("Unexpected success for scenario: %s", name)
		}
		if len(errs) > 0 && !scenario.expectedToFail {
			t.Errorf("Unexpected failure for scenario: %s - %+v", name, errs)
		}
	}
}
