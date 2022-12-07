package validation

import (
	"testing"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"k8s.io/component-base/config"
)

func TestValidateClientConnectionConfiguration(t *testing.T) {
	validConfig := &config.ClientConnectionConfiguration{
		AcceptContentTypes: "application/json",
		ContentType:        "application/json",
		QPS:                10,
		Burst:              10,
	}

	qpsLessThanZero := validConfig.DeepCopy()
	qpsLessThanZero.QPS = -1

	burstLessThanZero := validConfig.DeepCopy()
	burstLessThanZero.Burst = -1

	scenarios := map[string]struct {
		expectedToFail bool
		config         *config.ClientConnectionConfiguration
	}{
		"good": {
			expectedToFail: false,
			config:         validConfig,
		},
		"good-qps-less-than-zero": {
			expectedToFail: false,
			config:         qpsLessThanZero,
		},
		"bad-burst-less-then-zero": {
			expectedToFail: true,
			config:         burstLessThanZero,
		},
	}

	for name, scenario := range scenarios {
		errs := ValidateClientConnectionConfiguration(scenario.config, field.NewPath("clientConnectionConfiguration"))
		if len(errs) == 0 && scenario.expectedToFail {
			t.Errorf("Unexpected success for scenario: %s", name)
		}
		if len(errs) > 0 && !scenario.expectedToFail {
			t.Errorf("Unexpected failure for scenario: %s - %+v", name, errs)
		}
	}
}

func TestValidateLeaderElectionConfiguration(t *testing.T) {
	validConfig := &config.LeaderElectionConfiguration{
		ResourceLock:      "configmap",
		LeaderElect:       true,
		LeaseDuration:     metav1.Duration{Duration: 30 * time.Second},
		RenewDeadline:     metav1.Duration{Duration: 15 * time.Second},
		RetryPeriod:       metav1.Duration{Duration: 5 * time.Second},
		ResourceNamespace: "namespace",
		ResourceName:      "name",
	}

	renewDeadlineExceedsLeaseDuration := validConfig.DeepCopy()
	renewDeadlineExceedsLeaseDuration.RenewDeadline = metav1.Duration{Duration: 45 * time.Second}

	renewDeadlineZero := validConfig.DeepCopy()
	renewDeadlineZero.RenewDeadline = metav1.Duration{Duration: 0 * time.Second}

	leaseDurationZero := validConfig.DeepCopy()
	leaseDurationZero.LeaseDuration = metav1.Duration{Duration: 0 * time.Second}

	negativeValForRetryPeriod := validConfig.DeepCopy()
	negativeValForRetryPeriod.RetryPeriod = metav1.Duration{Duration: -45 * time.Second}

	negativeValForLeaseDuration := validConfig.DeepCopy()
	negativeValForLeaseDuration.LeaseDuration = metav1.Duration{Duration: -45 * time.Second}

	negativeValForRenewDeadline := validConfig.DeepCopy()
	negativeValForRenewDeadline.RenewDeadline = metav1.Duration{Duration: -45 * time.Second}

	LeaderElectButLeaderElectNotEnabled := validConfig.DeepCopy()
	LeaderElectButLeaderElectNotEnabled.LeaderElect = false
	LeaderElectButLeaderElectNotEnabled.LeaseDuration = metav1.Duration{Duration: -45 * time.Second}

	resourceLockNotDefined := validConfig.DeepCopy()
	resourceLockNotDefined.ResourceLock = ""

	resourceNameNotDefined := validConfig.DeepCopy()
	resourceNameNotDefined.ResourceName = ""

	resourceNamespaceNotDefined := validConfig.DeepCopy()
	resourceNamespaceNotDefined.ResourceNamespace = ""

	scenarios := map[string]struct {
		expectedToFail bool
		config         *config.LeaderElectionConfiguration
	}{
		"good": {
			expectedToFail: false,
			config:         validConfig,
		},
		"good-dont-check-leader-config-if-not-enabled": {
			expectedToFail: false,
			config:         LeaderElectButLeaderElectNotEnabled,
		},
		"bad-renew-deadline-exceeds-lease-duration": {
			expectedToFail: true,
			config:         renewDeadlineExceedsLeaseDuration,
		},
		"bad-negative-value-for-retry-period": {
			expectedToFail: true,
			config:         negativeValForRetryPeriod,
		},
		"bad-negative-value-for-lease-duration": {
			expectedToFail: true,
			config:         negativeValForLeaseDuration,
		},
		"bad-negative-value-for-renew-deadline": {
			expectedToFail: true,
			config:         negativeValForRenewDeadline,
		},
		"bad-renew-deadline-zero": {
			expectedToFail: true,
			config:         renewDeadlineZero,
		},
		"bad-lease-duration-zero": {
			expectedToFail: true,
			config:         leaseDurationZero,
		},
		"bad-resource-lock-not-defined": {
			expectedToFail: true,
			config:         resourceLockNotDefined,
		},
		"bad-resource-name-not-defined": {
			expectedToFail: true,
			config:         resourceNameNotDefined,
		},
		"bad-resource-namespace-not-defined": {
			expectedToFail: true,
			config:         resourceNamespaceNotDefined,
		},
	}

	for name, scenario := range scenarios {
		errs := ValidateLeaderElectionConfiguration(scenario.config, field.NewPath("leaderElectionConfiguration"))
		if len(errs) == 0 && scenario.expectedToFail {
			t.Errorf("Unexpected success for scenario: %s", name)
		}
		if len(errs) > 0 && !scenario.expectedToFail {
			t.Errorf("Unexpected failure for scenario: %s - %+v", name, errs)
		}
	}
}
