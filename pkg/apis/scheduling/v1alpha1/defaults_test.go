package v1alpha1_test

import (
	"reflect"
	"testing"

	"k8s.io/api/scheduling/v1alpha1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/kubernetes/pkg/api/legacyscheme"

	apiv1 "k8s.io/api/core/v1"
	utilfeature "k8s.io/apiserver/pkg/util/feature"
	featuregatetesting "k8s.io/component-base/featuregate/testing"
	// enforce that all types are installed
	_ "k8s.io/kubernetes/pkg/api/testapi"
	"k8s.io/kubernetes/pkg/features"
)

func roundTrip(t *testing.T, obj runtime.Object) runtime.Object {
	codec := legacyscheme.Codecs.LegacyCodec(v1alpha1.SchemeGroupVersion)
	data, err := runtime.Encode(codec, obj)
	if err != nil {
		t.Errorf("%v\n %#v", err, obj)
		return nil
	}
	obj2, err := runtime.Decode(codec, data)
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

func TestSetDefaultPreempting(t *testing.T) {
	priorityClass := &v1alpha1.PriorityClass{}

	// set NonPreemptingPriority true
	defer featuregatetesting.SetFeatureGateDuringTest(t, utilfeature.DefaultFeatureGate, features.NonPreemptingPriority, true)()

	output := roundTrip(t, runtime.Object(priorityClass)).(*v1alpha1.PriorityClass)
	if output.PreemptionPolicy == nil || *output.PreemptionPolicy != apiv1.PreemptLowerPriority {
		t.Errorf("Expected PriorityClass.Preempting value: %+v\ngot: %+v\n", apiv1.PreemptLowerPriority, output.PreemptionPolicy)
	}
}
