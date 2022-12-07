package options

import (
	"github.com/spf13/pflag"

	endpointconfig "k8s.io/kubernetes/pkg/controller/endpoint/config"
)

// EndpointControllerOptions holds the EndPointController options.
type EndpointControllerOptions struct {
	*endpointconfig.EndpointControllerConfiguration
}

// AddFlags adds flags related to EndPointController for controller manager to the specified FlagSet.
func (o *EndpointControllerOptions) AddFlags(fs *pflag.FlagSet) {
	if o == nil {
		return
	}

	fs.Int32Var(&o.ConcurrentEndpointSyncs, "concurrent-endpoint-syncs", o.ConcurrentEndpointSyncs, "The number of endpoint syncing operations that will be done concurrently. Larger number = faster endpoint updating, but more CPU (and network) load")
	fs.DurationVar(&o.EndpointUpdatesBatchPeriod.Duration, "endpoint-updates-batch-period", o.EndpointUpdatesBatchPeriod.Duration, "The length of endpoint updates batching period. Processing of pod changes will be delayed by this duration to join them with potential upcoming updates and reduce the overall number of endpoints updates. Larger number = higher endpoint programming latency, but lower number of endpoints revision generated")
}

// ApplyTo fills up EndPointController config with options.
func (o *EndpointControllerOptions) ApplyTo(cfg *endpointconfig.EndpointControllerConfiguration) error {
	if o == nil {
		return nil
	}

	cfg.ConcurrentEndpointSyncs = o.ConcurrentEndpointSyncs
	cfg.EndpointUpdatesBatchPeriod = o.EndpointUpdatesBatchPeriod

	return nil
}

// Validate checks validation of EndpointControllerOptions.
func (o *EndpointControllerOptions) Validate() []error {
	if o == nil {
		return nil
	}

	errs := []error{}
	return errs
}
