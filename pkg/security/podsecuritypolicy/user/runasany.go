package user

import (
	policy "k8s.io/api/policy/v1beta1"
	"k8s.io/apimachinery/pkg/util/validation/field"
	api "k8s.io/kubernetes/pkg/apis/core"
)

// runAsAny implements the interface RunAsUserStrategy.
type runAsAny struct{}

var _ RunAsUserStrategy = &runAsAny{}

// NewRunAsAny provides a strategy that will return nil.
func NewRunAsAny(options *policy.RunAsUserStrategyOptions) (RunAsUserStrategy, error) {
	return &runAsAny{}, nil
}

// Generate creates the uid based on policy rules.
func (s *runAsAny) Generate(pod *api.Pod, container *api.Container) (*int64, error) {
	return nil, nil
}

// Validate ensures that the specified values fall within the range of the strategy.
func (s *runAsAny) Validate(_ *field.Path, _ *api.Pod, _ *api.Container, runAsNonRoot *bool, runAsUser *int64) field.ErrorList {
	return field.ErrorList{}
}
