package certificates

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"k8s.io/api/certificates/v1beta1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestIsCertificateRequestApproved(t *testing.T) {
	testCases := []struct {
		name               string
		conditions         []v1beta1.CertificateSigningRequestCondition
		expectedIsApproved bool
	}{
		{
			"Not any conditions exist",
			nil,
			false,
		}, {
			"Approved not exist and Denied exist",
			[]v1beta1.CertificateSigningRequestCondition{
				{
					Type: v1beta1.CertificateDenied,
				},
			},
			false,
		}, {
			"Approved exist and Denied not exist",
			[]v1beta1.CertificateSigningRequestCondition{
				{
					Type: v1beta1.CertificateApproved,
				},
			},
			true,
		}, {
			"Both of Approved and Denied exist",
			[]v1beta1.CertificateSigningRequestCondition{
				{
					Type: v1beta1.CertificateApproved,
				},
				{
					Type: v1beta1.CertificateDenied,
				},
			},
			false,
		},
	}

	for _, tc := range testCases {
		csr := &v1beta1.CertificateSigningRequest{
			ObjectMeta: v1.ObjectMeta{
				Name: "fake-csr",
			},
			Status: v1beta1.CertificateSigningRequestStatus{
				Conditions: tc.conditions,
			},
		}

		assert.Equalf(t, tc.expectedIsApproved, IsCertificateRequestApproved(csr), "Failed to test: %s", tc.name)
	}
}
