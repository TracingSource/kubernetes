package sysctl

import (
	"k8s.io/apimachinery/pkg/util/validation/field"
	api "k8s.io/kubernetes/pkg/apis/core"
)

// SysctlsStrategy defines the interface for all sysctl strategies.
type SysctlsStrategy interface {
	// Validate ensures that the specified values fall within the range of the strategy.
	Validate(pod *api.Pod) field.ErrorList
}
