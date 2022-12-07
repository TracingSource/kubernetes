package options

import (
	"github.com/spf13/pflag"

	daemonconfig "k8s.io/kubernetes/pkg/controller/daemon/config"
)

// DaemonSetControllerOptions holds the DaemonSetController options.
type DaemonSetControllerOptions struct {
	*daemonconfig.DaemonSetControllerConfiguration
}

// AddFlags adds flags related to DaemonSetController for controller manager to the specified FlagSet.
func (o *DaemonSetControllerOptions) AddFlags(fs *pflag.FlagSet) {
	if o == nil {
		return
	}
}

// ApplyTo fills up DaemonSetController config with options.
func (o *DaemonSetControllerOptions) ApplyTo(cfg *daemonconfig.DaemonSetControllerConfiguration) error {
	if o == nil {
		return nil
	}

	cfg.ConcurrentDaemonSetSyncs = o.ConcurrentDaemonSetSyncs

	return nil
}

// Validate checks validation of DaemonSetControllerOptions.
func (o *DaemonSetControllerOptions) Validate() []error {
	if o == nil {
		return nil
	}

	errs := []error{}
	return errs
}
