package certificates

import (
	"context"
	"reflect"
	"testing"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/diff"
	"k8s.io/apiserver/pkg/authentication/user"
	genericapirequest "k8s.io/apiserver/pkg/endpoints/request"
	certapi "k8s.io/kubernetes/pkg/apis/certificates"
)

func TestStrategyCreate(t *testing.T) {
	tests := map[string]struct {
		ctx         context.Context
		obj         runtime.Object
		expectedObj runtime.Object
	}{
		"no user in context, no user in obj": {
			ctx: genericapirequest.NewContext(),
			obj: &certapi.CertificateSigningRequest{},
			expectedObj: &certapi.CertificateSigningRequest{
				Status: certapi.CertificateSigningRequestStatus{Conditions: []certapi.CertificateSigningRequestCondition{}},
			},
		},
		"user in context, no user in obj": {
			ctx: genericapirequest.WithUser(
				genericapirequest.NewContext(),
				&user.DefaultInfo{
					Name:   "bob",
					UID:    "123",
					Groups: []string{"group1"},
					Extra:  map[string][]string{"foo": {"bar"}},
				},
			),
			obj: &certapi.CertificateSigningRequest{},
			expectedObj: &certapi.CertificateSigningRequest{
				Spec: certapi.CertificateSigningRequestSpec{
					Username: "bob",
					UID:      "123",
					Groups:   []string{"group1"},
					Extra:    map[string]certapi.ExtraValue{"foo": {"bar"}},
				},
				Status: certapi.CertificateSigningRequestStatus{Conditions: []certapi.CertificateSigningRequestCondition{}},
			},
		},
		"no user in context, user in obj": {
			ctx: genericapirequest.NewContext(),
			obj: &certapi.CertificateSigningRequest{
				Spec: certapi.CertificateSigningRequestSpec{
					Username: "bob",
					UID:      "123",
					Groups:   []string{"group1"},
				},
			},
			expectedObj: &certapi.CertificateSigningRequest{
				Status: certapi.CertificateSigningRequestStatus{Conditions: []certapi.CertificateSigningRequestCondition{}},
			},
		},
		"user in context, user in obj": {
			ctx: genericapirequest.WithUser(
				genericapirequest.NewContext(),
				&user.DefaultInfo{
					Name: "alice",
					UID:  "234",
				},
			),
			obj: &certapi.CertificateSigningRequest{
				Spec: certapi.CertificateSigningRequestSpec{
					Username: "bob",
					UID:      "123",
					Groups:   []string{"group1"},
				},
			},
			expectedObj: &certapi.CertificateSigningRequest{
				Spec: certapi.CertificateSigningRequestSpec{
					Username: "alice",
					UID:      "234",
					Groups:   nil,
				},
				Status: certapi.CertificateSigningRequestStatus{Conditions: []certapi.CertificateSigningRequestCondition{}},
			},
		},
		"pre-approved status": {
			ctx: genericapirequest.NewContext(),
			obj: &certapi.CertificateSigningRequest{
				Status: certapi.CertificateSigningRequestStatus{
					Conditions: []certapi.CertificateSigningRequestCondition{
						{Type: certapi.CertificateApproved},
					},
				},
			},
			expectedObj: &certapi.CertificateSigningRequest{
				Status: certapi.CertificateSigningRequestStatus{Conditions: []certapi.CertificateSigningRequestCondition{}},
			}},
	}

	for k, tc := range tests {
		obj := tc.obj
		Strategy.PrepareForCreate(tc.ctx, obj)
		if !reflect.DeepEqual(obj, tc.expectedObj) {
			t.Errorf("%s: object diff: %s", k, diff.ObjectDiff(obj, tc.expectedObj))
		}
	}
}
