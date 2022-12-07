package v1alpha1

import (
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kubectrlmgrconfigv1alpha1 "k8s.io/kube-controller-manager/config/v1alpha1"
)

// RecommendedDefaultDeploymentControllerConfiguration defaults a pointer to a
// DeploymentControllerConfiguration struct. This will set the recommended default
// values, but they may be subject to change between API versions. This function
// is intentionally not registered in the scheme as a "normal" `SetDefaults_Foo`
// function to allow consumers of this type to set whatever defaults for their
// embedded configs. Forcing consumers to use these defaults would be problematic
// as defaulting in the scheme is done as part of the conversion, and there would
// be no easy way to opt-out. Instead, if you want to use this defaulting method
// run it in your wrapper struct of this type in its `SetDefaults_` method.
func RecommendedDefaultDeploymentControllerConfiguration(obj *kubectrlmgrconfigv1alpha1.DeploymentControllerConfiguration) {
	zero := metav1.Duration{}
	if obj.ConcurrentDeploymentSyncs == 0 {
		obj.ConcurrentDeploymentSyncs = 5
	}
	if obj.DeploymentControllerSyncPeriod == zero {
		obj.DeploymentControllerSyncPeriod = metav1.Duration{Duration: 30 * time.Second}
	}
}
