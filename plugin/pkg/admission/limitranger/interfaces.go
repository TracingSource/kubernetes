package limitranger

import (
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apiserver/pkg/admission"
)

// LimitRangerActions is an interface defining actions to be carried over ranges to identify and manipulate their limits
type LimitRangerActions interface {
	// MutateLimit is a pluggable function to set limits on the object.
	MutateLimit(limitRange *corev1.LimitRange, kind string, obj runtime.Object) error
	// ValidateLimits is a pluggable function to enforce limits on the object.
	ValidateLimit(limitRange *corev1.LimitRange, kind string, obj runtime.Object) error
	// SupportsAttributes is a pluggable function to allow overridding what resources the limitranger
	// supports.
	SupportsAttributes(attr admission.Attributes) bool
	// SupportsLimit is a pluggable function to allow ignoring limits that should not be applied
	// for any reason.
	SupportsLimit(limitRange *corev1.LimitRange) bool
}
