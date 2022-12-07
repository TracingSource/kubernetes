package webhook

import (
	"encoding/json"
	"reflect"
	"testing"

	jsonpatch "github.com/evanphx/json-patch"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func TestPatches(t *testing.T) {
	testCases := []struct {
		patch    string
		initial  interface{}
		expected interface{}
		toTest   interface{}
	}{
		{
			patch: configMapPatch1,
			initial: corev1.ConfigMap{
				Data: map[string]string{
					"mutation-start": "yes",
				},
			},
			expected: &corev1.ConfigMap{
				Data: map[string]string{
					"mutation-start":   "yes",
					"mutation-stage-1": "yes",
				},
			},
		},
		{
			patch: configMapPatch2,
			initial: corev1.ConfigMap{
				Data: map[string]string{
					"mutation-start": "yes",
				},
			},
			expected: &corev1.ConfigMap{
				Data: map[string]string{
					"mutation-start":   "yes",
					"mutation-stage-2": "yes",
				},
			},
		},

		{
			patch: podsInitContainerPatch,
			initial: corev1.Pod{
				Spec: corev1.PodSpec{
					InitContainers: []corev1.Container{},
				},
			},
			expected: &corev1.Pod{
				Spec: corev1.PodSpec{
					InitContainers: []corev1.Container{
						{
							Image:     "webhook-added-image",
							Name:      "webhook-added-init-container",
							Resources: corev1.ResourceRequirements{},
						},
					},
				},
			},
		},
	}
	for _, testcase := range testCases {
		objJS, err := json.Marshal(testcase.initial)
		if err != nil {
			t.Fatal(err)
		}
		patchObj, err := jsonpatch.DecodePatch([]byte(testcase.patch))
		if err != nil {
			t.Fatal(err)
		}

		patchedJS, err := patchObj.Apply(objJS)
		if err != nil {
			t.Fatal(err)
		}
		objType := reflect.TypeOf(testcase.initial)
		objTest := reflect.New(objType).Interface()
		err = json.Unmarshal(patchedJS, objTest)
		if err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(objTest, testcase.expected) {
			t.Errorf("\nexpected %#v\n, got %#v", testcase.expected, objTest)
		}
	}

}

func TestJSONPatchForUnstructured(t *testing.T) {
	cr := &unstructured.Unstructured{
		Object: map[string]interface{}{
			"kind":       "Something",
			"apiVersion": "somegroup/v1",
			"data": map[string]interface{}{
				"mutation-start": "yes",
			},
		},
	}
	crJS, err := json.Marshal(cr)
	if err != nil {
		t.Fatal(err)
	}

	patchObj, err := jsonpatch.DecodePatch([]byte(configMapPatch1))
	if err != nil {
		t.Fatal(err)
	}
	patchedJS, err := patchObj.Apply(crJS)
	patchedObj := unstructured.Unstructured{}
	err = json.Unmarshal(patchedJS, &patchedObj)
	if err != nil {
		t.Fatal(err)
	}
	expectedData := map[string]interface{}{
		"mutation-start":   "yes",
		"mutation-stage-1": "yes",
	}

	if !reflect.DeepEqual(patchedObj.Object["data"], expectedData) {
		t.Errorf("\nexpected %#v\n, got %#v", expectedData, patchedObj.Object["data"])
	}
}
