package options

import (
	"github.com/spf13/pflag"

	namespaceconfig "k8s.io/kubernetes/pkg/controller/namespace/config"
)

// NamespaceControllerOptions holds the NamespaceController options.
type NamespaceControllerOptions struct {
	*namespaceconfig.NamespaceControllerConfiguration
}

// AddFlags adds flags related to NamespaceController for controller manager to the specified FlagSet.
func (o *NamespaceControllerOptions) AddFlags(fs *pflag.FlagSet) {
	if o == nil {
		return
	}

	fs.DurationVar(&o.NamespaceSyncPeriod.Duration, "namespace-sync-period", o.NamespaceSyncPeriod.Duration, "The period for syncing namespace life-cycle updates")
	fs.Int32Var(&o.ConcurrentNamespaceSyncs, "concurrent-namespace-syncs", o.ConcurrentNamespaceSyncs, "The number of namespace objects that are allowed to sync concurrently. Larger number = more responsive namespace termination, but more CPU (and network) load")
}

// ApplyTo fills up NamespaceController config with options.
func (o *NamespaceControllerOptions) ApplyTo(cfg *namespaceconfig.NamespaceControllerConfiguration) error {
	if o == nil {
		return nil
	}

	cfg.NamespaceSyncPeriod = o.NamespaceSyncPeriod
	cfg.ConcurrentNamespaceSyncs = o.ConcurrentNamespaceSyncs

	return nil
}

// Validate checks validation of NamespaceControllerOptions.
func (o *NamespaceControllerOptions) Validate() []error {
	if o == nil {
		return nil
	}

	errs := []error{}
	return errs
}
