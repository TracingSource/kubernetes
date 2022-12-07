package resource

import (
	"reflect"
	"strings"
	"testing"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func TestDynamicCodecDecode(t *testing.T) {
	testcases := []struct {
		name string
		data []byte
		gvk  *schema.GroupVersionKind
		obj  runtime.Object

		expectErr string
		expectGVK *schema.GroupVersionKind
		expectObj runtime.Object
	}{
		{
			name:      "v1 Status",
			data:      []byte(`{"apiVersion":"v1","kind":"Status"}`),
			expectGVK: &schema.GroupVersionKind{"", "v1", "Status"},
			expectObj: &metav1.Status{TypeMeta: metav1.TypeMeta{APIVersion: "v1", Kind: "Status"}},
		},
		{
			name:      "meta.k8s.io/v1 Status",
			data:      []byte(`{"apiVersion":"meta.k8s.io/v1","kind":"Status"}`),
			expectGVK: &schema.GroupVersionKind{"meta.k8s.io", "v1", "Status"},
			expectObj: &metav1.Status{TypeMeta: metav1.TypeMeta{APIVersion: "meta.k8s.io/v1", Kind: "Status"}},
		},
		{
			name:      "example.com/v1 Status",
			data:      []byte(`{"apiVersion":"example.com/v1","kind":"Status"}`),
			expectGVK: &schema.GroupVersionKind{"example.com", "v1", "Status"},
			expectObj: &unstructured.Unstructured{Object: map[string]interface{}{"apiVersion": "example.com/v1", "kind": "Status"}},
		},
		{
			name:      "example.com/v1 Foo",
			data:      []byte(`{"apiVersion":"example.com/v1","kind":"Foo"}`),
			expectGVK: &schema.GroupVersionKind{"example.com", "v1", "Foo"},
			expectObj: &unstructured.Unstructured{Object: map[string]interface{}{"apiVersion": "example.com/v1", "kind": "Foo"}},
		},
	}

	for _, test := range testcases {
		t.Run(test.name, func(t *testing.T) {
			obj, gvk, err := dynamicCodec{}.Decode(test.data, test.gvk, test.obj)
			if (err == nil) != (test.expectErr == "") {
				t.Fatalf("expected err=%v, got %v", test.expectErr, err)
			}
			if err != nil && !strings.Contains(err.Error(), test.expectErr) {
				t.Fatalf("expected err=%v, got %v", test.expectErr, err)
			}
			if err != nil {
				return
			}

			if !reflect.DeepEqual(test.expectGVK, gvk) {
				t.Errorf("expected\n\tgvk=%#v\ngot\n\t%#v", test.expectGVK, gvk)
			}
			if !reflect.DeepEqual(test.expectObj, obj) {
				t.Errorf("expected\n\t%#v\n\t%#v", test.expectObj, obj)
			}
		})
	}
}
