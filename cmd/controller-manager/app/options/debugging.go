package options

import (
	"github.com/spf13/pflag"

	componentbaseconfig "k8s.io/component-base/config"
)

// DebuggingOptions holds the Debugging options.
type DebuggingOptions struct {
	*componentbaseconfig.DebuggingConfiguration
}

// AddFlags adds flags related to debugging for controller manager to the specified FlagSet.
func (o *DebuggingOptions) AddFlags(fs *pflag.FlagSet) {
	if o == nil {
		return
	}

	fs.BoolVar(&o.EnableProfiling, "profiling", o.EnableProfiling,
		"Enable profiling via web interface host:port/debug/pprof/")
	fs.BoolVar(&o.EnableContentionProfiling, "contention-profiling", o.EnableContentionProfiling,
		"Enable lock contention profiling, if profiling is enabled")
}

// ApplyTo fills up Debugging config with options.
func (o *DebuggingOptions) ApplyTo(cfg *componentbaseconfig.DebuggingConfiguration) error {
	if o == nil {
		return nil
	}

	cfg.EnableProfiling = o.EnableProfiling
	cfg.EnableContentionProfiling = o.EnableContentionProfiling

	return nil
}

// Validate checks validation of DebuggingOptions.
func (o *DebuggingOptions) Validate() []error {
	if o == nil {
		return nil
	}

	errs := []error{}
	return errs
}
