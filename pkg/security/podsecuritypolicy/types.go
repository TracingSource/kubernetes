package podsecuritypolicy

import (
	policy "k8s.io/api/policy/v1beta1"
	"k8s.io/apimachinery/pkg/util/validation/field"
	api "k8s.io/kubernetes/pkg/apis/core"
	"k8s.io/kubernetes/pkg/security/podsecuritypolicy/apparmor"
	"k8s.io/kubernetes/pkg/security/podsecuritypolicy/capabilities"
	"k8s.io/kubernetes/pkg/security/podsecuritypolicy/group"
	"k8s.io/kubernetes/pkg/security/podsecuritypolicy/seccomp"
	"k8s.io/kubernetes/pkg/security/podsecuritypolicy/selinux"
	"k8s.io/kubernetes/pkg/security/podsecuritypolicy/sysctl"
	"k8s.io/kubernetes/pkg/security/podsecuritypolicy/user"
)

// Provider provides the implementation to generate a new security
// context based on constraints or validate an existing security context against constraints.
type Provider interface {
	// MutatePod sets the default values of the required but not filled fields of the pod and all
	// containers in the pod.
	MutatePod(pod *api.Pod) error
	// ValidatePod ensures a pod and all its containers are in compliance with the given constraints.
	// ValidatePod MUST NOT mutate the pod.
	ValidatePod(pod *api.Pod) field.ErrorList
	// Get the name of the PSP that this provider was initialized with.
	GetPSPName() string
}

// StrategyFactory abstracts how the strategies are created from the provider so that you may
// implement your own custom strategies that may pull information from other resources as necessary.
// For example, if you would like to populate the strategies with values from namespace annotations
// you may create a factory with a client that can pull the namespace and populate the appropriate
// values.
type StrategyFactory interface {
	// CreateStrategies creates the strategies that a provider will use.  The namespace argument
	// should be the namespace of the object being checked (the pod's namespace).
	CreateStrategies(psp *policy.PodSecurityPolicy, namespace string) (*ProviderStrategies, error)
}

// ProviderStrategies is a holder for all strategies that the provider requires to be populated.
type ProviderStrategies struct {
	RunAsUserStrategy         user.RunAsUserStrategy
	RunAsGroupStrategy        group.GroupStrategy
	SELinuxStrategy           selinux.SELinuxStrategy
	AppArmorStrategy          apparmor.Strategy
	FSGroupStrategy           group.GroupStrategy
	SupplementalGroupStrategy group.GroupStrategy
	CapabilitiesStrategy      capabilities.Strategy
	SysctlsStrategy           sysctl.SysctlsStrategy
	SeccompStrategy           seccomp.Strategy
}
