package v1_test

import (
	"reflect"
	"testing"

	autoscalingv1 "k8s.io/api/autoscaling/v1"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/kubernetes/pkg/api/legacyscheme"
	_ "k8s.io/kubernetes/pkg/apis/autoscaling/install"
	. "k8s.io/kubernetes/pkg/apis/autoscaling/v1"
	_ "k8s.io/kubernetes/pkg/apis/core/install"
	utilpointer "k8s.io/utils/pointer"
)

func TestSetDefaultHPA(t *testing.T) {
	tests := []struct {
		hpa            autoscalingv1.HorizontalPodAutoscaler
		expectReplicas int32
		test           string
	}{
		{
			hpa:            autoscalingv1.HorizontalPodAutoscaler{},
			expectReplicas: 1,
			test:           "unspecified min replicas, use the default value",
		},
		{
			hpa: autoscalingv1.HorizontalPodAutoscaler{
				Spec: autoscalingv1.HorizontalPodAutoscalerSpec{
					MinReplicas: utilpointer.Int32Ptr(3),
				},
			},
			expectReplicas: 3,
			test:           "set min replicas to 3",
		},
	}

	for _, test := range tests {
		hpa := &test.hpa
		obj2 := roundTrip(t, runtime.Object(hpa))
		hpa2, ok := obj2.(*autoscalingv1.HorizontalPodAutoscaler)
		if !ok {
			t.Fatalf("unexpected object: %v", obj2)
		}
		if hpa2.Spec.MinReplicas == nil {
			t.Errorf("unexpected nil MinReplicas")
		} else if test.expectReplicas != *hpa2.Spec.MinReplicas {
			t.Errorf("expected: %d MinReplicas, got: %d", test.expectReplicas, *hpa2.Spec.MinReplicas)
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
