package options

import (
	"github.com/spf13/pflag"

	attachdetachconfig "k8s.io/kubernetes/pkg/controller/volume/attachdetach/config"
)

// AttachDetachControllerOptions holds the AttachDetachController options.
type AttachDetachControllerOptions struct {
	*attachdetachconfig.AttachDetachControllerConfiguration
}

// AddFlags adds flags related to AttachDetachController for controller manager to the specified FlagSet.
func (o *AttachDetachControllerOptions) AddFlags(fs *pflag.FlagSet) {
	if o == nil {
		return
	}

	fs.BoolVar(&o.DisableAttachDetachReconcilerSync, "disable-attach-detach-reconcile-sync", false, "Disable volume attach detach reconciler sync. Disabling this may cause volumes to be mismatched with pods. Use wisely.")
	fs.DurationVar(&o.ReconcilerSyncLoopPeriod.Duration, "attach-detach-reconcile-sync-period", o.ReconcilerSyncLoopPeriod.Duration, "The reconciler sync wait time between volume attach detach. This duration must be larger than one second, and increasing this value from the default may allow for volumes to be mismatched with pods.")
}

// ApplyTo fills up AttachDetachController config with options.
func (o *AttachDetachControllerOptions) ApplyTo(cfg *attachdetachconfig.AttachDetachControllerConfiguration) error {
	if o == nil {
		return nil
	}

	cfg.DisableAttachDetachReconcilerSync = o.DisableAttachDetachReconcilerSync
	cfg.ReconcilerSyncLoopPeriod = o.ReconcilerSyncLoopPeriod

	return nil
}

// Validate checks validation of AttachDetachControllerOptions.
func (o *AttachDetachControllerOptions) Validate() []error {
	if o == nil {
		return nil
	}

	errs := []error{}
	return errs
}
