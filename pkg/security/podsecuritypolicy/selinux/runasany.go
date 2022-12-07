package selinux

import (
	policy "k8s.io/api/policy/v1beta1"
	"k8s.io/apimachinery/pkg/util/validation/field"
	api "k8s.io/kubernetes/pkg/apis/core"
)

// runAsAny implements the SELinuxStrategy interface.
type runAsAny struct{}

var _ SELinuxStrategy = &runAsAny{}

// NewRunAsAny provides a strategy that will return the configured se linux context or nil.
func NewRunAsAny(options *policy.SELinuxStrategyOptions) (SELinuxStrategy, error) {
	return &runAsAny{}, nil
}

// Generate creates the SELinuxOptions based on constraint rules.
func (s *runAsAny) Generate(pod *api.Pod, container *api.Container) (*api.SELinuxOptions, error) {
	return nil, nil
}

// Validate ensures that the specified values fall within the range of the strategy.
func (s *runAsAny) Validate(fldPath *field.Path, _ *api.Pod, _ *api.Container, options *api.SELinuxOptions) field.ErrorList {
	return field.ErrorList{}
}
