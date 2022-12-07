package options

import (
	"github.com/spf13/pflag"

	csrsigningconfig "k8s.io/kubernetes/pkg/controller/certificates/signer/config"
)

const (
	// These defaults are deprecated and exported so that we can warn if
	// they are being used.

	// DefaultClusterSigningCertFile is deprecated. Do not use.
	DefaultClusterSigningCertFile = "/etc/kubernetes/ca/ca.pem"
	// DefaultClusterSigningKeyFile is deprecated. Do not use.
	DefaultClusterSigningKeyFile = "/etc/kubernetes/ca/ca.key"
)

// CSRSigningControllerOptions holds the CSRSigningController options.
type CSRSigningControllerOptions struct {
	*csrsigningconfig.CSRSigningControllerConfiguration
}

// AddFlags adds flags related to CSRSigningController for controller manager to the specified FlagSet.
func (o *CSRSigningControllerOptions) AddFlags(fs *pflag.FlagSet) {
	if o == nil {
		return
	}

	fs.StringVar(&o.ClusterSigningCertFile, "cluster-signing-cert-file", o.ClusterSigningCertFile, "Filename containing a PEM-encoded X509 CA certificate used to issue cluster-scoped certificates")
	fs.StringVar(&o.ClusterSigningKeyFile, "cluster-signing-key-file", o.ClusterSigningKeyFile, "Filename containing a PEM-encoded RSA or ECDSA private key used to sign cluster-scoped certificates")
	fs.DurationVar(&o.ClusterSigningDuration.Duration, "experimental-cluster-signing-duration", o.ClusterSigningDuration.Duration, "The length of duration signed certificates will be given.")
}

// ApplyTo fills up CSRSigningController config with options.
func (o *CSRSigningControllerOptions) ApplyTo(cfg *csrsigningconfig.CSRSigningControllerConfiguration) error {
	if o == nil {
		return nil
	}

	cfg.ClusterSigningCertFile = o.ClusterSigningCertFile
	cfg.ClusterSigningKeyFile = o.ClusterSigningKeyFile
	cfg.ClusterSigningDuration = o.ClusterSigningDuration

	return nil
}

// Validate checks validation of CSRSigningControllerOptions.
func (o *CSRSigningControllerOptions) Validate() []error {
	if o == nil {
		return nil
	}

	errs := []error{}
	return errs
}
