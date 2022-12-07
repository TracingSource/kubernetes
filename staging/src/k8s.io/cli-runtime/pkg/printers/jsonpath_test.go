package printers

import (
	"bytes"
	"testing"

	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/sets"
)

func TestPrinters(t *testing.T) {
	om := func(name string) metav1.ObjectMeta { return metav1.ObjectMeta{Name: name} }

	jsonpathPrinter, err := NewJSONPathPrinter("{.metadata.name}")
	if err != nil {
		t.Fatal(err)
	}

	objects := map[string]runtime.Object{
		"pod":             &v1.Pod{ObjectMeta: om("pod")},
		"emptyPodList":    &v1.PodList{},
		"nonEmptyPodList": &v1.PodList{Items: []v1.Pod{{}}},
		"endpoints": &v1.Endpoints{
			Subsets: []v1.EndpointSubset{{
				Addresses: []v1.EndpointAddress{{IP: "127.0.0.1"}, {IP: "localhost"}},
				Ports:     []v1.EndpointPort{{Port: 8080}},
			}}},
	}

	// Set of strings representing objects that should produce an error.
	expectedErrors := sets.NewString("emptyPodList", "nonEmptyPodList", "endpoints")

	for oName, obj := range objects {
		b := &bytes.Buffer{}
		if err := jsonpathPrinter.PrintObj(obj, b); err != nil {
			if expectedErrors.Has(oName) {
				// expected error
				continue
			}
			t.Errorf("JSONPathPrinter error object '%v'; error: '%v'", oName, err)
		}
	}
}
