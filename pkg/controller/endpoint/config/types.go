package config

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EndpointControllerConfiguration contains elements describing EndpointController.
type EndpointControllerConfiguration struct {
	// concurrentEndpointSyncs is the number of endpoint syncing operations
	// that will be done concurrently. Larger number = faster endpoint updating,
	// but more CPU (and network) load.
	ConcurrentEndpointSyncs int32

	// EndpointUpdatesBatchPeriod can be used to batch endpoint updates.
	// All updates of endpoint triggered by pod change will be delayed by up to
	// 'EndpointUpdatesBatchPeriod'. If other pods in the same endpoint change
	// in that period, they will be batched to a single endpoint update.
	// Default 0 value means that each pod update triggers an endpoint update.
	EndpointUpdatesBatchPeriod metav1.Duration
}
