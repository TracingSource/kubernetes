package event

import (
	"reflect"
	"testing"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/util/diff"
	apitesting "k8s.io/kubernetes/pkg/api/testing"
	api "k8s.io/kubernetes/pkg/apis/core"

	// install all api groups for testing
	_ "k8s.io/kubernetes/pkg/api/testapi"
)

func TestGetAttrs(t *testing.T) {
	eventA := &api.Event{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "f0118",
			Namespace: "default",
		},
		InvolvedObject: api.ObjectReference{
			Kind:            "Pod",
			Name:            "foo",
			Namespace:       "baz",
			UID:             "long uid string",
			APIVersion:      "v1",
			ResourceVersion: "0",
			FieldPath:       "",
		},
		Reason: "ForTesting",
		Source: api.EventSource{Component: "test"},
		Type:   api.EventTypeNormal,
	}
	field := ToSelectableFields(eventA)
	expect := fields.Set{
		"metadata.name":                  "f0118",
		"metadata.namespace":             "default",
		"involvedObject.kind":            "Pod",
		"involvedObject.name":            "foo",
		"involvedObject.namespace":       "baz",
		"involvedObject.uid":             "long uid string",
		"involvedObject.apiVersion":      "v1",
		"involvedObject.resourceVersion": "0",
		"involvedObject.fieldPath":       "",
		"reason":                         "ForTesting",
		"source":                         "test",
		"type":                           api.EventTypeNormal,
	}
	if e, a := expect, field; !reflect.DeepEqual(e, a) {
		t.Errorf("diff: %s", diff.ObjectDiff(e, a))
	}
}

func TestSelectableFieldLabelConversions(t *testing.T) {
	fset := ToSelectableFields(&api.Event{})
	apitesting.TestSelectableFieldLabelConversionsOfKind(t,
		"v1",
		"Event",
		fset,
		nil,
	)
}
