package options

import (
	"github.com/spf13/pflag"

	replicasetconfig "k8s.io/kubernetes/pkg/controller/replicaset/config"
)

// ReplicaSetControllerOptions holds the ReplicaSetController options.
type ReplicaSetControllerOptions struct {
	*replicasetconfig.ReplicaSetControllerConfiguration
}

// AddFlags adds flags related to ReplicaSetController for controller manager to the specified FlagSet.
func (o *ReplicaSetControllerOptions) AddFlags(fs *pflag.FlagSet) {
	if o == nil {
		return
	}

	fs.Int32Var(&o.ConcurrentRSSyncs, "concurrent-replicaset-syncs", o.ConcurrentRSSyncs, "The number of replica sets that are allowed to sync concurrently. Larger number = more responsive replica management, but more CPU (and network) load")
}

// ApplyTo fills up ReplicaSetController config with options.
func (o *ReplicaSetControllerOptions) ApplyTo(cfg *replicasetconfig.ReplicaSetControllerConfiguration) error {
	if o == nil {
		return nil
	}

	cfg.ConcurrentRSSyncs = o.ConcurrentRSSyncs

	return nil
}

// Validate checks validation of ReplicaSetControllerOptions.
func (o *ReplicaSetControllerOptions) Validate() []error {
	if o == nil {
		return nil
	}

	errs := []error{}
	return errs
}
