package v1alpha1

import (
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kubectrlmgrconfigv1alpha1 "k8s.io/kube-controller-manager/config/v1alpha1"
)

// RecommendedDefaultCSRSigningControllerConfiguration defaults a pointer to a
// CSRSigningControllerConfiguration struct. This will set the recommended default
// values, but they may be subject to change between API versions. This function
// is intentionally not registered in the scheme as a "normal" `SetDefaults_Foo`
// function to allow consumers of this type to set whatever defaults for their
// embedded configs. Forcing consumers to use these defaults would be problematic
// as defaulting in the scheme is done as part of the conversion, and there would
// be no easy way to opt-out. Instead, if you want to use this defaulting method
// run it in your wrapper struct of this type in its `SetDefaults_` method.
func RecommendedDefaultCSRSigningControllerConfiguration(obj *kubectrlmgrconfigv1alpha1.CSRSigningControllerConfiguration) {
	zero := metav1.Duration{}
	if obj.ClusterSigningCertFile == "" {
		obj.ClusterSigningCertFile = "/etc/kubernetes/ca/ca.pem"
	}
	if obj.ClusterSigningKeyFile == "" {
		obj.ClusterSigningKeyFile = "/etc/kubernetes/ca/ca.key"
	}
	if obj.ClusterSigningDuration == zero {
		obj.ClusterSigningDuration = metav1.Duration{Duration: 365 * 24 * time.Hour}
	}
}
