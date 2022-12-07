package v1alpha1

import (
	apiv1 "k8s.io/api/core/v1"
	"k8s.io/api/scheduling/v1alpha1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	utilfeature "k8s.io/apiserver/pkg/util/feature"
	"k8s.io/kubernetes/pkg/features"
)

func addDefaultingFuncs(scheme *runtime.Scheme) error {
	return RegisterDefaults(scheme)
}

// SetDefaults_PriorityClass sets additional defaults compared to its counterpart
// in extensions.
func SetDefaults_PriorityClass(obj *v1alpha1.PriorityClass) {
	if utilfeature.DefaultFeatureGate.Enabled(features.NonPreemptingPriority) && obj.PreemptionPolicy == nil {
		preemptLowerPriority := apiv1.PreemptLowerPriority
		obj.PreemptionPolicy = &preemptLowerPriority
	}
}
