package renewal

import (
	"crypto/x509"
	"testing"

	certutil "k8s.io/client-go/util/cert"
)

func TestFileRenewer(t *testing.T) {
	// creates a File renewer using a test Certificate authority
	fr := NewFileRenewer(testCACert, testCAKey)

	// renews a certificate
	certCfg := &certutil.Config{
		CommonName: "test-certs",
		AltNames: certutil.AltNames{
			DNSNames: []string{"test-domain.space"},
		},
		Usages: []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
	}

	cert, _, err := fr.Renew(certCfg)
	if err != nil {
		t.Fatalf("unexpected error renewing cert: %v", err)
	}

	// verify the renewed certificate
	pool := x509.NewCertPool()
	pool.AddCert(testCACert)

	_, err = cert.Verify(x509.VerifyOptions{
		DNSName:   "test-domain.space",
		Roots:     pool,
		KeyUsages: []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
	})
	if err != nil {
		t.Errorf("couldn't verify new cert: %v", err)
	}

}
