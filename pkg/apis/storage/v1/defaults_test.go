package v1_test

import (
	"reflect"
	"testing"

	storagev1 "k8s.io/api/storage/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/kubernetes/pkg/api/legacyscheme"
	_ "k8s.io/kubernetes/pkg/apis/storage/install"
)

func roundTrip(t *testing.T, obj runtime.Object) runtime.Object {
	codec := legacyscheme.Codecs.LegacyCodec(storagev1.SchemeGroupVersion)
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

func TestSetDefaultVolumeBindingMode(t *testing.T) {
	class := &storagev1.StorageClass{}

	// field should be defaulted
	defaultMode := storagev1.VolumeBindingImmediate
	output := roundTrip(t, runtime.Object(class)).(*storagev1.StorageClass)
	outMode := output.VolumeBindingMode
	if outMode == nil {
		t.Errorf("Expected VolumeBindingMode to be defaulted to: %+v, got: nil", defaultMode)
	} else if *outMode != defaultMode {
		t.Errorf("Expected VolumeBindingMode to be defaulted to: %+v, got: %+v", defaultMode, outMode)
	}
}
