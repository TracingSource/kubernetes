package install

import (
	"encoding/json"
	"reflect"
	"testing"

	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/kubernetes/pkg/api/legacyscheme"
	internal "k8s.io/kubernetes/pkg/apis/core"
)

func TestResourceVersioner(t *testing.T) {
	pod := internal.Pod{ObjectMeta: metav1.ObjectMeta{ResourceVersion: "10"}}
	version, err := meta.NewAccessor().ResourceVersion(&pod)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if version != "10" {
		t.Errorf("unexpected version %v", version)
	}

	podList := internal.PodList{ListMeta: metav1.ListMeta{ResourceVersion: "10"}}
	version, err = meta.NewAccessor().ResourceVersion(&podList)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if version != "10" {
		t.Errorf("unexpected version %v", version)
	}
}

func TestCodec(t *testing.T) {
	pod := internal.Pod{}
	// We do want to use package registered rather than testapi here, because we
	// want to test if the package install and package registered work as expected.
	data, err := runtime.Encode(legacyscheme.Codecs.LegacyCodec(schema.GroupVersion{Group: "", Version: "v1"}), &pod)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	other := internal.Pod{}
	if err := json.Unmarshal(data, &other); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if other.APIVersion != "v1" || other.Kind != "Pod" {
		t.Errorf("unexpected unmarshalled object %#v", other)
	}
}

func TestUnversioned(t *testing.T) {
	for _, obj := range []runtime.Object{
		&metav1.Status{},
	} {
		if unversioned, ok := legacyscheme.Scheme.IsUnversioned(obj); !unversioned || !ok {
			t.Errorf("%v is expected to be unversioned", reflect.TypeOf(obj))
		}
	}
}
