package options

import (
	"github.com/spf13/pflag"

	deploymentconfig "k8s.io/kubernetes/pkg/controller/deployment/config"
)

// DeploymentControllerOptions holds the DeploymentController options.
type DeploymentControllerOptions struct {
	*deploymentconfig.DeploymentControllerConfiguration
}

// AddFlags adds flags related to DeploymentController for controller manager to the specified FlagSet.
func (o *DeploymentControllerOptions) AddFlags(fs *pflag.FlagSet) {
	if o == nil {
		return
	}

	fs.Int32Var(&o.ConcurrentDeploymentSyncs, "concurrent-deployment-syncs", o.ConcurrentDeploymentSyncs, "The number of deployment objects that are allowed to sync concurrently. Larger number = more responsive deployments, but more CPU (and network) load")
	fs.DurationVar(&o.DeploymentControllerSyncPeriod.Duration, "deployment-controller-sync-period", o.DeploymentControllerSyncPeriod.Duration, "Period for syncing the deployments.")
}

// ApplyTo fills up DeploymentController config with options.
func (o *DeploymentControllerOptions) ApplyTo(cfg *deploymentconfig.DeploymentControllerConfiguration) error {
	if o == nil {
		return nil
	}

	cfg.ConcurrentDeploymentSyncs = o.ConcurrentDeploymentSyncs
	cfg.DeploymentControllerSyncPeriod = o.DeploymentControllerSyncPeriod

	return nil
}

// Validate checks validation of DeploymentControllerOptions.
func (o *DeploymentControllerOptions) Validate() []error {
	if o == nil {
		return nil
	}

	errs := []error{}
	return errs
}
