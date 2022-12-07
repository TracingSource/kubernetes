package options

import (
	"fmt"

	"github.com/spf13/pflag"

	statefulsetconfig "k8s.io/kubernetes/pkg/controller/statefulset/config"
)

// StatefulSetControllerOptions holds the StatefulSetController options.
type StatefulSetControllerOptions struct {
	*statefulsetconfig.StatefulSetControllerConfiguration
}

// AddFlags adds flags related to StatefulSetController for controller manager to the specified FlagSet.
func (o *StatefulSetControllerOptions) AddFlags(fs *pflag.FlagSet) {
	if o == nil {
		return
	}

	fs.Int32Var(&o.ConcurrentStatefulSetSyncs, "concurrent-statefulset-syncs", o.ConcurrentStatefulSetSyncs, "The number of statefulset objects that are allowed to sync concurrently. Larger number = more responsive statefulsets, but more CPU (and network) load")
}

// ApplyTo fills up StatefulSetController config with options.
func (o *StatefulSetControllerOptions) ApplyTo(cfg *statefulsetconfig.StatefulSetControllerConfiguration) error {
	if o == nil {
		return nil
	}

	cfg.ConcurrentStatefulSetSyncs = o.ConcurrentStatefulSetSyncs

	return nil
}

// Validate checks validation of StatefulSetControllerOptions.
func (o *StatefulSetControllerOptions) Validate() []error {
	if o == nil {
		return nil
	}

	errs := []error{}
	if o.ConcurrentStatefulSetSyncs < 1 {
		errs = append(errs, fmt.Errorf("concurrent-statefulset-syncs must be greater than 0, but got %d", o.ConcurrentStatefulSetSyncs))
	}
	return errs
}
