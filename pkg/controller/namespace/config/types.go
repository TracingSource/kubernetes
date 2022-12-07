package config

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// NamespaceControllerConfiguration contains elements describing NamespaceController.
type NamespaceControllerConfiguration struct {
	// namespaceSyncPeriod is the period for syncing namespace life-cycle
	// updates.
	NamespaceSyncPeriod metav1.Duration
	// concurrentNamespaceSyncs is the number of namespace objects that are
	// allowed to sync concurrently.
	ConcurrentNamespaceSyncs int32
}
