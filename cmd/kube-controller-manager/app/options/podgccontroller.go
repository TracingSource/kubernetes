package options

import (
	"github.com/spf13/pflag"

	podgcconfig "k8s.io/kubernetes/pkg/controller/podgc/config"
)

// PodGCControllerOptions holds the PodGCController options.
type PodGCControllerOptions struct {
	*podgcconfig.PodGCControllerConfiguration
}

// AddFlags adds flags related to PodGCController for controller manager to the specified FlagSet.
func (o *PodGCControllerOptions) AddFlags(fs *pflag.FlagSet) {
	if o == nil {
		return
	}

	fs.Int32Var(&o.TerminatedPodGCThreshold, "terminated-pod-gc-threshold", o.TerminatedPodGCThreshold, "Number of terminated pods that can exist before the terminated pod garbage collector starts deleting terminated pods. If <= 0, the terminated pod garbage collector is disabled.")
}

// ApplyTo fills up PodGCController config with options.
func (o *PodGCControllerOptions) ApplyTo(cfg *podgcconfig.PodGCControllerConfiguration) error {
	if o == nil {
		return nil
	}

	cfg.TerminatedPodGCThreshold = o.TerminatedPodGCThreshold

	return nil
}

// Validate checks validation of PodGCControllerOptions.
func (o *PodGCControllerOptions) Validate() []error {
	if o == nil {
		return nil
	}

	errs := []error{}
	return errs
}
