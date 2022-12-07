package config

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ResourceQuotaControllerConfiguration contains elements describing ResourceQuotaController.
type ResourceQuotaControllerConfiguration struct {
	// resourceQuotaSyncPeriod is the period for syncing quota usage status
	// in the system.
	ResourceQuotaSyncPeriod metav1.Duration
	// concurrentResourceQuotaSyncs is the number of resource quotas that are
	// allowed to sync concurrently. Larger number = more responsive quota
	// management, but more CPU (and network) load.
	ConcurrentResourceQuotaSyncs int32
}
