// Package options contains flags and options for initializing kube-apiserver
package options

import (
	"fmt"
	"net"

	utilnet "k8s.io/apimachinery/pkg/util/net"
	genericoptions "k8s.io/apiserver/pkg/server/options"
)

// NewSecureServingOptions gives default values for the kube-apiserver which are not the options wanted by
// "normal" API servers running on the platform
func NewSecureServingOptions() *genericoptions.SecureServingOptionsWithLoopback {
	o := genericoptions.SecureServingOptions{
		BindAddress: net.ParseIP("0.0.0.0"),
		BindPort:    6443,
		Required:    true,
		ServerCert: genericoptions.GeneratableKeyCert{
			PairName:      "apiserver",
			CertDirectory: "/var/run/kubernetes",
		},
	}
	return o.WithLoopback()
}

// NewInsecureServingOptions gives default values for the kube-apiserver.
// TODO: switch insecure serving off by default
func NewInsecureServingOptions() *genericoptions.DeprecatedInsecureServingOptionsWithLoopback {
	o := genericoptions.DeprecatedInsecureServingOptions{
		BindAddress: net.ParseIP("127.0.0.1"),
		BindPort:    8080,
	}
	return o.WithLoopback()
}

// DefaultAdvertiseAddress sets the field AdvertiseAddress if
// unset. The field will be set based on the SecureServingOptions. If
// the SecureServingOptions is not present, DefaultExternalAddress
// will fall back to the insecure ServingOptions.
func DefaultAdvertiseAddress(s *genericoptions.ServerRunOptions, insecure *genericoptions.DeprecatedInsecureServingOptions) error {
	if insecure == nil {
		return nil
	}

	if s.AdvertiseAddress == nil || s.AdvertiseAddress.IsUnspecified() {
		hostIP, err := utilnet.ResolveBindAddress(insecure.BindAddress)
		if err != nil {
			return fmt.Errorf("unable to find suitable network address.error='%v'. "+
				"Try to set the AdvertiseAddress directly or provide a valid BindAddress to fix this", err)
		}
		s.AdvertiseAddress = hostIP
	}

	return nil
}
