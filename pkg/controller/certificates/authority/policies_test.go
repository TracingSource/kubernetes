package authority

import (
	"crypto/x509"
	"fmt"
	"reflect"
	"testing"

	capi "k8s.io/api/certificates/v1beta1"
)

func TestKeyUsagesFromStrings(t *testing.T) {
	testcases := []struct {
		usages              []capi.KeyUsage
		expectedKeyUsage    x509.KeyUsage
		expectedExtKeyUsage []x509.ExtKeyUsage
		expectErr           bool
	}{
		{
			usages:              []capi.KeyUsage{"signing"},
			expectedKeyUsage:    x509.KeyUsageDigitalSignature,
			expectedExtKeyUsage: nil,
			expectErr:           false,
		},
		{
			usages:              []capi.KeyUsage{"client auth"},
			expectedKeyUsage:    0,
			expectedExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
			expectErr:           false,
		},
		{
			usages:              []capi.KeyUsage{"client auth", "client auth"},
			expectedKeyUsage:    0,
			expectedExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
			expectErr:           false,
		},
		{
			usages:              []capi.KeyUsage{"cert sign", "encipher only"},
			expectedKeyUsage:    x509.KeyUsageCertSign | x509.KeyUsageEncipherOnly,
			expectedExtKeyUsage: nil,
			expectErr:           false,
		},
		{
			usages:              []capi.KeyUsage{"ocsp signing", "crl sign", "s/mime", "content commitment"},
			expectedKeyUsage:    x509.KeyUsageCRLSign | x509.KeyUsageContentCommitment,
			expectedExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageEmailProtection, x509.ExtKeyUsageOCSPSigning},
			expectErr:           false,
		},
		{
			usages:              []capi.KeyUsage{"unsupported string"},
			expectedKeyUsage:    0,
			expectedExtKeyUsage: nil,
			expectErr:           true,
		},
	}

	for _, tc := range testcases {
		t.Run(fmt.Sprint(tc.usages), func(t *testing.T) {
			ku, eku, err := keyUsagesFromStrings(tc.usages)

			if tc.expectErr {
				if err == nil {
					t.Errorf("did not return an error, but expected one")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			if ku != tc.expectedKeyUsage || !reflect.DeepEqual(eku, tc.expectedExtKeyUsage) {
				t.Errorf("got=(%v, %v), want=(%v, %v)", ku, eku, tc.expectedKeyUsage, tc.expectedExtKeyUsage)
			}
		})
	}
}
