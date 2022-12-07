package v2beta1_test

import (
	"reflect"
	"testing"

	autoscalingv2beta1 "k8s.io/api/autoscaling/v2beta1"
	"k8s.io/api/core/v1"
	apiequality "k8s.io/apimachinery/pkg/api/equality"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/kubernetes/pkg/api/legacyscheme"
	"k8s.io/kubernetes/pkg/apis/autoscaling"
	_ "k8s.io/kubernetes/pkg/apis/autoscaling/install"
	. "k8s.io/kubernetes/pkg/apis/autoscaling/v2beta1"
	_ "k8s.io/kubernetes/pkg/apis/core/install"
	utilpointer "k8s.io/utils/pointer"
)

func TestSetDefaultHPA(t *testing.T) {
	utilizationDefaultVal := int32(autoscaling.DefaultCPUUtilization)
	defaultReplicas := utilpointer.Int32Ptr(1)
	defaultTemplate := []autoscalingv2beta1.MetricSpec{
		{
			Type: autoscalingv2beta1.ResourceMetricSourceType,
			Resource: &autoscalingv2beta1.ResourceMetricSource{
				Name:                     v1.ResourceCPU,
				TargetAverageUtilization: &utilizationDefaultVal,
			},
		},
	}

	tests := []struct {
		original *autoscalingv2beta1.HorizontalPodAutoscaler
		expected *autoscalingv2beta1.HorizontalPodAutoscaler
	}{
		{ // MinReplicas default value
			original: &autoscalingv2beta1.HorizontalPodAutoscaler{
				Spec: autoscalingv2beta1.HorizontalPodAutoscalerSpec{
					Metrics: defaultTemplate,
				},
			},
			expected: &autoscalingv2beta1.HorizontalPodAutoscaler{
				Spec: autoscalingv2beta1.HorizontalPodAutoscalerSpec{
					MinReplicas: defaultReplicas,
					Metrics:     defaultTemplate,
				},
			},
		},
		{ // MinReplicas update
			original: &autoscalingv2beta1.HorizontalPodAutoscaler{
				Spec: autoscalingv2beta1.HorizontalPodAutoscalerSpec{
					MinReplicas: utilpointer.Int32Ptr(3),
					Metrics:     defaultTemplate,
				},
			},
			expected: &autoscalingv2beta1.HorizontalPodAutoscaler{
				Spec: autoscalingv2beta1.HorizontalPodAutoscalerSpec{
					MinReplicas: utilpointer.Int32Ptr(3),
					Metrics:     defaultTemplate,
				},
			},
		},
		{ // Metrics default value
			original: &autoscalingv2beta1.HorizontalPodAutoscaler{
				Spec: autoscalingv2beta1.HorizontalPodAutoscalerSpec{
					MinReplicas: defaultReplicas,
				},
			},
			expected: &autoscalingv2beta1.HorizontalPodAutoscaler{
				Spec: autoscalingv2beta1.HorizontalPodAutoscalerSpec{
					MinReplicas: defaultReplicas,
					Metrics:     defaultTemplate,
				},
			},
		},
	}

	for i, test := range tests {
		original := test.original
		expected := test.expected
		obj2 := roundTrip(t, runtime.Object(original))
		got, ok := obj2.(*autoscalingv2beta1.HorizontalPodAutoscaler)
		if !ok {
			t.Fatalf("(%d) unexpected object: %v", i, obj2)
		}
		if !apiequality.Semantic.DeepEqual(got.Spec, expected.Spec) {
			t.Errorf("(%d) got different than expected\ngot:\n\t%+v\nexpected:\n\t%+v", i, got.Spec, expected.Spec)
		}
	}
}

func roundTrip(t *testing.T, obj runtime.Object) runtime.Object {
	data, err := runtime.Encode(legacyscheme.Codecs.LegacyCodec(SchemeGroupVersion), obj)
	if err != nil {
		t.Errorf("%v\n %#v", err, obj)
		return nil
	}
	obj2, err := runtime.Decode(legacyscheme.Codecs.UniversalDecoder(), data)
	if err != nil {
		t.Errorf("%v\nData: %s\nSource: %#v", err, string(data), obj)
		return nil
	}
	obj3 := reflect.New(reflect.TypeOf(obj).Elem()).Interface().(runtime.Object)
	err = legacyscheme.Scheme.Convert(obj2, obj3, nil)
	if err != nil {
		t.Errorf("%v\nSource: %#v", err, obj2)
		return nil
	}
	return obj3
}
