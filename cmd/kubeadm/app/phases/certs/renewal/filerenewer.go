package renewal

import (
	"crypto"
	"crypto/x509"

	certutil "k8s.io/client-go/util/cert"
	"k8s.io/kubernetes/cmd/kubeadm/app/util/pkiutil"
)

// FileRenewer define a certificate renewer implementation that uses given CA cert and key for generating new certficiates
type FileRenewer struct {
	caCert *x509.Certificate
	caKey  crypto.Signer
}

// NewFileRenewer returns a new certificate renewer that uses given CA cert and key for generating new certficiates
func NewFileRenewer(caCert *x509.Certificate, caKey crypto.Signer) *FileRenewer {
	return &FileRenewer{
		caCert: caCert,
		caKey:  caKey,
	}
}

// Renew a certificate using a given CA cert and key
func (r *FileRenewer) Renew(cfg *certutil.Config) (*x509.Certificate, crypto.Signer, error) {
	return pkiutil.NewCertAndKey(r.caCert, r.caKey, cfg)
}
