package options

import (
	"github.com/spf13/pflag"

	serviceconfig "k8s.io/kubernetes/pkg/controller/service/config"
)

// ServiceControllerOptions holds the ServiceController options.
type ServiceControllerOptions struct {
	*serviceconfig.ServiceControllerConfiguration
}

// AddFlags adds flags related to ServiceController for controller manager to the specified FlagSet.
func (o *ServiceControllerOptions) AddFlags(fs *pflag.FlagSet) {
	if o == nil {
		return
	}

	fs.Int32Var(&o.ConcurrentServiceSyncs, "concurrent-service-syncs", o.ConcurrentServiceSyncs, "The number of services that are allowed to sync concurrently. Larger number = more responsive service management, but more CPU (and network) load")
}

// ApplyTo fills up ServiceController config with options.
func (o *ServiceControllerOptions) ApplyTo(cfg *serviceconfig.ServiceControllerConfiguration) error {
	if o == nil {
		return nil
	}

	cfg.ConcurrentServiceSyncs = o.ConcurrentServiceSyncs

	return nil
}

// Validate checks validation of ServiceControllerOptions.
func (o *ServiceControllerOptions) Validate() []error {
	if o == nil {
		return nil
	}

	errs := []error{}
	return errs
}
