package v1

import (
	autoscalingv1 "k8s.io/api/autoscaling/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

func addDefaultingFuncs(scheme *runtime.Scheme) error {
	return RegisterDefaults(scheme)
}

func SetDefaults_HorizontalPodAutoscaler(obj *autoscalingv1.HorizontalPodAutoscaler) {
	if obj.Spec.MinReplicas == nil {
		minReplicas := int32(1)
		obj.Spec.MinReplicas = &minReplicas
	}

	// NB: we apply a default for CPU utilization in conversion because
	// we need access to the annotations to properly apply the default.
}
