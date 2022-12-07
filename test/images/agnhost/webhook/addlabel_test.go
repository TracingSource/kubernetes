package webhook

import (
	"encoding/json"
	"reflect"
	"testing"

	jsonpatch "github.com/evanphx/json-patch"
	"k8s.io/api/admission/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

func TestAddLabel(t *testing.T) {
	testCases := []struct {
		name           string
		initialLabels  map[string]string
		expectedLabels map[string]string
	}{
		{
			name:           "add first label",
			initialLabels:  nil,
			expectedLabels: map[string]string{"added-label": "yes"},
		},
		{
			name:           "add second label",
			initialLabels:  map[string]string{"other-label": "yes"},
			expectedLabels: map[string]string{"other-label": "yes", "added-label": "yes"},
		},
		{
			name:           "idempotent update label",
			initialLabels:  map[string]string{"added-label": "yes"},
			expectedLabels: map[string]string{"added-label": "yes"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			request := corev1.ConfigMap{ObjectMeta: metav1.ObjectMeta{Labels: tc.initialLabels}}
			raw, err := json.Marshal(request)
			if err != nil {
				t.Fatal(err)
			}
			review := v1.AdmissionReview{Request: &v1.AdmissionRequest{Object: runtime.RawExtension{Raw: raw}}}
			response := addLabel(review)
			if response.Patch != nil {
				patchObj, err := jsonpatch.DecodePatch([]byte(response.Patch))
				if err != nil {
					t.Fatal(err)
				}
				raw, err = patchObj.Apply(raw)
				if err != nil {
					t.Fatal(err)
				}
			}

			objType := reflect.TypeOf(request)
			objTest := reflect.New(objType).Interface()
			err = json.Unmarshal(raw, objTest)
			if err != nil {
				t.Fatal(err)
			}
			actual := objTest.(*corev1.ConfigMap)
			if !reflect.DeepEqual(actual.Labels, tc.expectedLabels) {
				t.Errorf("\nexpected %#v, got %#v, patch: %v", actual.Labels, tc.expectedLabels, string(response.Patch))
			}
		})
	}
}
