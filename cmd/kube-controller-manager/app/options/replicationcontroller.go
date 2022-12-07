package options

import (
	"github.com/spf13/pflag"

	replicationconfig "k8s.io/kubernetes/pkg/controller/replication/config"
)

// ReplicationControllerOptions holds the ReplicationController options.
type ReplicationControllerOptions struct {
	*replicationconfig.ReplicationControllerConfiguration
}

// AddFlags adds flags related to ReplicationController for controller manager to the specified FlagSet.
func (o *ReplicationControllerOptions) AddFlags(fs *pflag.FlagSet) {
	if o == nil {
		return
	}

	fs.Int32Var(&o.ConcurrentRCSyncs, "concurrent_rc_syncs", o.ConcurrentRCSyncs, "The number of replication controllers that are allowed to sync concurrently. Larger number = more responsive replica management, but more CPU (and network) load")
}

// ApplyTo fills up ReplicationController config with options.
func (o *ReplicationControllerOptions) ApplyTo(cfg *replicationconfig.ReplicationControllerConfiguration) error {
	if o == nil {
		return nil
	}

	cfg.ConcurrentRCSyncs = o.ConcurrentRCSyncs

	return nil
}

// Validate checks validation of ReplicationControllerOptions.
func (o *ReplicationControllerOptions) Validate() []error {
	if o == nil {
		return nil
	}

	errs := []error{}
	return errs
}
