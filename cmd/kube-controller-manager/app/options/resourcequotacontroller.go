package options

import (
	"github.com/spf13/pflag"

	resourcequotaconfig "k8s.io/kubernetes/pkg/controller/resourcequota/config"
)

// ResourceQuotaControllerOptions holds the ResourceQuotaController options.
type ResourceQuotaControllerOptions struct {
	*resourcequotaconfig.ResourceQuotaControllerConfiguration
}

// AddFlags adds flags related to ResourceQuotaController for controller manager to the specified FlagSet.
func (o *ResourceQuotaControllerOptions) AddFlags(fs *pflag.FlagSet) {
	if o == nil {
		return
	}

	fs.DurationVar(&o.ResourceQuotaSyncPeriod.Duration, "resource-quota-sync-period", o.ResourceQuotaSyncPeriod.Duration, "The period for syncing quota usage status in the system")
	fs.Int32Var(&o.ConcurrentResourceQuotaSyncs, "concurrent-resource-quota-syncs", o.ConcurrentResourceQuotaSyncs, "The number of resource quotas that are allowed to sync concurrently. Larger number = more responsive quota management, but more CPU (and network) load")
}

// ApplyTo fills up ResourceQuotaController config with options.
func (o *ResourceQuotaControllerOptions) ApplyTo(cfg *resourcequotaconfig.ResourceQuotaControllerConfiguration) error {
	if o == nil {
		return nil
	}

	cfg.ResourceQuotaSyncPeriod = o.ResourceQuotaSyncPeriod
	cfg.ConcurrentResourceQuotaSyncs = o.ConcurrentResourceQuotaSyncs

	return nil
}

// Validate checks validation of ResourceQuotaControllerOptions.
func (o *ResourceQuotaControllerOptions) Validate() []error {
	if o == nil {
		return nil
	}

	errs := []error{}
	return errs
}
