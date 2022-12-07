package v1alpha1

import (
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kubectrlmgrconfigv1alpha1 "k8s.io/kube-controller-manager/config/v1alpha1"
)

// RecommendedDefaultAttachDetachControllerConfiguration defaults a pointer to a
// AttachDetachControllerConfiguration struct. This will set the recommended default
// values, but they may be subject to change between API versions. This function
// is intentionally not registered in the scheme as a "normal" `SetDefaults_Foo`
// function to allow consumers of this type to set whatever defaults for their
// embedded configs. Forcing consumers to use these defaults would be problematic
// as defaulting in the scheme is done as part of the conversion, and there would
// be no easy way to opt-out. Instead, if you want to use this defaulting method
// run it in your wrapper struct of this type in its `SetDefaults_` method.
func RecommendedDefaultAttachDetachControllerConfiguration(obj *kubectrlmgrconfigv1alpha1.AttachDetachControllerConfiguration) {
	zero := metav1.Duration{}
	if obj.ReconcilerSyncLoopPeriod == zero {
		obj.ReconcilerSyncLoopPeriod = metav1.Duration{Duration: 60 * time.Second}
	}
}
