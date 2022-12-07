package options

import (
	"github.com/spf13/pflag"

	ttlafterfinishedconfig "k8s.io/kubernetes/pkg/controller/ttlafterfinished/config"
)

// TTLAfterFinishedControllerOptions holds the TTLAfterFinishedController options.
type TTLAfterFinishedControllerOptions struct {
	*ttlafterfinishedconfig.TTLAfterFinishedControllerConfiguration
}

// AddFlags adds flags related to TTLAfterFinishedController for controller manager to the specified FlagSet.
func (o *TTLAfterFinishedControllerOptions) AddFlags(fs *pflag.FlagSet) {
	if o == nil {
		return
	}

	fs.Int32Var(&o.ConcurrentTTLSyncs, "concurrent-ttl-after-finished-syncs", o.ConcurrentTTLSyncs, "The number of TTL-after-finished controller workers that are allowed to sync concurrently.")
}

// ApplyTo fills up TTLAfterFinishedController config with options.
func (o *TTLAfterFinishedControllerOptions) ApplyTo(cfg *ttlafterfinishedconfig.TTLAfterFinishedControllerConfiguration) error {
	if o == nil {
		return nil
	}

	cfg.ConcurrentTTLSyncs = o.ConcurrentTTLSyncs

	return nil
}

// Validate checks validation of TTLAfterFinishedControllerOptions.
func (o *TTLAfterFinishedControllerOptions) Validate() []error {
	if o == nil {
		return nil
	}

	errs := []error{}
	return errs
}
