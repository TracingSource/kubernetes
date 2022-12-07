package secret

import (
	"reflect"
	"testing"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	genericapirequest "k8s.io/apiserver/pkg/endpoints/request"
	apitesting "k8s.io/kubernetes/pkg/api/testing"
	api "k8s.io/kubernetes/pkg/apis/core"

	// install all api groups for testing
	_ "k8s.io/kubernetes/pkg/api/testapi"
)

func TestExportSecret(t *testing.T) {
	tests := []struct {
		objIn     runtime.Object
		objOut    runtime.Object
		exact     bool
		expectErr bool
	}{
		{
			objIn: &api.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "foo",
					Namespace: "bar",
				},
				Data: map[string][]byte{
					"foo": []byte("bar"),
				},
			},
			objOut: &api.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "foo",
					Namespace: "bar",
				},
				Data: map[string][]byte{
					"foo": []byte("bar"),
				},
			},
			exact: true,
		},
		{
			objIn: &api.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "foo",
					Namespace: "bar",
				},
				Type: api.SecretTypeServiceAccountToken,
			},
			expectErr: true,
		},
		{
			objIn: &api.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "foo",
					Namespace: "bar",
					Annotations: map[string]string{
						api.ServiceAccountUIDKey: "true",
					},
				},
			},
			expectErr: true,
		},
		{
			objIn:     &api.Pod{},
			expectErr: true,
		},
	}

	for _, test := range tests {
		err := Strategy.Export(genericapirequest.NewContext(), test.objIn, test.exact)
		if err != nil {
			if !test.expectErr {
				t.Errorf("unexpected error: %v", err)
			}
			continue
		}
		if test.expectErr {
			t.Error("unexpected non-error")
			continue
		}
		if !reflect.DeepEqual(test.objIn, test.objOut) {
			t.Errorf("expected:\n%v\nsaw:\n%v\n", test.objOut, test.objIn)
		}
	}
}

func TestSelectableFieldLabelConversions(t *testing.T) {
	apitesting.TestSelectableFieldLabelConversionsOfKind(t,
		"v1",
		"Secret",
		SelectableFields(&api.Secret{}),
		nil,
	)
}
