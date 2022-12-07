package genericclioptions

import (
	"bytes"
	"testing"

	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/scheme"
)

func TestNamePrinter(t *testing.T) {
	tests := map[string]struct {
		obj    runtime.Object
		expect string
	}{
		"singleObject": {
			&v1.Pod{
				TypeMeta: metav1.TypeMeta{
					Kind: "Pod",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name: "foo",
				},
			},
			"pod/foo\n"},
		"List": {
			&unstructured.UnstructuredList{
				Object: map[string]interface{}{
					"kind":       "List",
					"apiVersion": "v1",
				},
				Items: []unstructured.Unstructured{
					{
						Object: map[string]interface{}{
							"kind":       "Pod",
							"apiVersion": "v1",
							"metadata": map[string]interface{}{
								"name": "bar",
							},
						},
					},
				},
			},
			"pod/bar\n"},
	}

	printFlags := NewPrintFlags("").WithTypeSetter(scheme.Scheme).WithDefaultOutput("name")
	printer, err := printFlags.ToPrinter()
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}

	for name, item := range tests {
		buff := &bytes.Buffer{}
		err := printer.PrintObj(item.obj, buff)
		if err != nil {
			t.Errorf("%v: unexpected err: %v", name, err)
			continue
		}
		got := buff.String()
		if item.expect != got {
			t.Errorf("%v: expected %v, got %v", name, item.expect, got)
		}
	}
}
