package v1beta1_test

import (
	"reflect"
	"testing"

	"k8s.io/apiserver/pkg/authentication/user"
	"k8s.io/kubernetes/pkg/apis/abac"
	"k8s.io/kubernetes/pkg/apis/abac/v1beta1"
)

func TestV1Beta1Conversion(t *testing.T) {
	testcases := map[string]struct {
		old      *v1beta1.Policy
		expected *abac.Policy
	}{
		// specifying a user is preserved
		"user": {
			old:      &v1beta1.Policy{Spec: v1beta1.PolicySpec{User: "bob"}},
			expected: &abac.Policy{Spec: abac.PolicySpec{User: "bob"}},
		},

		// specifying a group is preserved
		"group": {
			old:      &v1beta1.Policy{Spec: v1beta1.PolicySpec{Group: "mygroup"}},
			expected: &abac.Policy{Spec: abac.PolicySpec{Group: "mygroup"}},
		},

		// specifying * for user or group maps to all authenticated subjects
		"* user": {
			old:      &v1beta1.Policy{Spec: v1beta1.PolicySpec{User: "*"}},
			expected: &abac.Policy{Spec: abac.PolicySpec{Group: user.AllAuthenticated}},
		},
		"* group": {
			old:      &v1beta1.Policy{Spec: v1beta1.PolicySpec{Group: "*"}},
			expected: &abac.Policy{Spec: abac.PolicySpec{Group: user.AllAuthenticated}},
		},
	}
	for k, tc := range testcases {
		internal := &abac.Policy{}
		if err := abac.Scheme.Convert(tc.old, internal, nil); err != nil {
			t.Errorf("%s: unexpected error: %v", k, err)
		}
		if !reflect.DeepEqual(internal, tc.expected) {
			t.Errorf("%s: expected\n\t%#v, got \n\t%#v", k, tc.expected, internal)
		}
	}
}
