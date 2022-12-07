package options

import (
	"github.com/spf13/pflag"

	kubectrlmgrconfig "k8s.io/kubernetes/pkg/controller/apis/config"
)

// CloudProviderOptions holds the cloudprovider options.
type CloudProviderOptions struct {
	*kubectrlmgrconfig.CloudProviderConfiguration
}

// Validate checks validation of cloudprovider options.
func (s *CloudProviderOptions) Validate() []error {
	allErrors := []error{}
	return allErrors
}

// AddFlags adds flags related to cloudprovider for controller manager to the specified FlagSet.
func (s *CloudProviderOptions) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&s.Name, "cloud-provider", s.Name,
		"The provider for cloud services. Empty string for no provider.")

	fs.StringVar(&s.CloudConfigFile, "cloud-config", s.CloudConfigFile,
		"The path to the cloud provider configuration file. Empty string for no configuration file.")
}

// ApplyTo fills up cloudprovider config with options.
func (s *CloudProviderOptions) ApplyTo(cfg *kubectrlmgrconfig.CloudProviderConfiguration) error {
	if s == nil {
		return nil
	}

	cfg.Name = s.Name
	cfg.CloudConfigFile = s.CloudConfigFile

	return nil
}
