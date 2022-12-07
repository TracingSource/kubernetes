package v1alpha1

import (
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kubectrlmgrconfigv1alpha1 "k8s.io/kube-controller-manager/config/v1alpha1"
	utilpointer "k8s.io/utils/pointer"
)

// RecommendedDefaultNodeLifecycleControllerConfiguration defaults a pointer to a
// NodeLifecycleControllerConfiguration struct. This will set the recommended default
// values, but they may be subject to change between API versions. This function
// is intentionally not registered in the scheme as a "normal" `SetDefaults_Foo`
// function to allow consumers of this type to set whatever defaults for their
// embedded configs. Forcing consumers to use these defaults would be problematic
// as defaulting in the scheme is done as part of the conversion, and there would
// be no easy way to opt-out. Instead, if you want to use this defaulting method
// run it in your wrapper struct of this type in its `SetDefaults_` method.
func RecommendedDefaultNodeLifecycleControllerConfiguration(obj *kubectrlmgrconfigv1alpha1.NodeLifecycleControllerConfiguration) {
	zero := metav1.Duration{}
	if obj.PodEvictionTimeout == zero {
		obj.PodEvictionTimeout = metav1.Duration{Duration: 5 * time.Minute}
	}
	if obj.NodeMonitorGracePeriod == zero {
		obj.NodeMonitorGracePeriod = metav1.Duration{Duration: 40 * time.Second}
	}
	if obj.NodeStartupGracePeriod == zero {
		obj.NodeStartupGracePeriod = metav1.Duration{Duration: 60 * time.Second}
	}
	if obj.EnableTaintManager == nil {
		obj.EnableTaintManager = utilpointer.BoolPtr(true)
	}
}
