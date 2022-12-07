package config

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// CSRSigningControllerConfiguration contains elements describing CSRSigningController.
type CSRSigningControllerConfiguration struct {
	// clusterSigningCertFile is the filename containing a PEM-encoded
	// X509 CA certificate used to issue cluster-scoped certificates
	ClusterSigningCertFile string
	// clusterSigningCertFile is the filename containing a PEM-encoded
	// RSA or ECDSA private key used to issue cluster-scoped certificates
	ClusterSigningKeyFile string
	// clusterSigningDuration is the length of duration signed certificates
	// will be given.
	ClusterSigningDuration metav1.Duration
}
