package config

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// DeploymentControllerConfiguration contains elements describing DeploymentController.
type DeploymentControllerConfiguration struct {
	// concurrentDeploymentSyncs is the number of deployment objects that are
	// allowed to sync concurrently. Larger number = more responsive deployments,
	// but more CPU (and network) load.
	ConcurrentDeploymentSyncs int32
	// deploymentControllerSyncPeriod is the period for syncing the deployments.
	DeploymentControllerSyncPeriod metav1.Duration
}
