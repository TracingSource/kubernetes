package v2beta1

import (
	autoscalingv2beta1 "k8s.io/api/autoscaling/v2beta1"
	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/kubernetes/pkg/apis/autoscaling"
)

func addDefaultingFuncs(scheme *runtime.Scheme) error {
	return RegisterDefaults(scheme)
}

func SetDefaults_HorizontalPodAutoscaler(obj *autoscalingv2beta1.HorizontalPodAutoscaler) {
	if obj.Spec.MinReplicas == nil {
		minReplicas := int32(1)
		obj.Spec.MinReplicas = &minReplicas
	}

	if len(obj.Spec.Metrics) == 0 {
		utilizationDefaultVal := int32(autoscaling.DefaultCPUUtilization)
		obj.Spec.Metrics = []autoscalingv2beta1.MetricSpec{
			{
				Type: autoscalingv2beta1.ResourceMetricSourceType,
				Resource: &autoscalingv2beta1.ResourceMetricSource{
					Name:                     v1.ResourceCPU,
					TargetAverageUtilization: &utilizationDefaultVal,
				},
			},
		}
	}
}
