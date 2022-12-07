package v1beta1

import (
	"crypto/x509"
	"encoding/pem"
	"errors"

	certificatesv1beta1 "k8s.io/api/certificates/v1beta1"
)

// ParseCSR extracts the CSR from the API object and decodes it.
func ParseCSR(obj *certificatesv1beta1.CertificateSigningRequest) (*x509.CertificateRequest, error) {
	// extract PEM from request object
	pemBytes := obj.Spec.Request
	block, _ := pem.Decode(pemBytes)
	if block == nil || block.Type != "CERTIFICATE REQUEST" {
		return nil, errors.New("PEM block type must be CERTIFICATE REQUEST")
	}
	csr, err := x509.ParseCertificateRequest(block.Bytes)
	if err != nil {
		return nil, err
	}
	return csr, nil
}
