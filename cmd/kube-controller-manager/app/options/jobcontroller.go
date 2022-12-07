package options

import (
	"github.com/spf13/pflag"

	jobconfig "k8s.io/kubernetes/pkg/controller/job/config"
)

// JobControllerOptions holds the JobController options.
type JobControllerOptions struct {
	*jobconfig.JobControllerConfiguration
}

// AddFlags adds flags related to JobController for controller manager to the specified FlagSet.
func (o *JobControllerOptions) AddFlags(fs *pflag.FlagSet) {
	if o == nil {
		return
	}
}

// ApplyTo fills up JobController config with options.
func (o *JobControllerOptions) ApplyTo(cfg *jobconfig.JobControllerConfiguration) error {
	if o == nil {
		return nil
	}

	cfg.ConcurrentJobSyncs = o.ConcurrentJobSyncs

	return nil
}

// Validate checks validation of JobControllerOptions.
func (o *JobControllerOptions) Validate() []error {
	if o == nil {
		return nil
	}

	errs := []error{}
	return errs
}
