package lease

import (
	"context"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"k8s.io/apiserver/pkg/storage/names"
	"k8s.io/kubernetes/pkg/api/legacyscheme"
	"k8s.io/kubernetes/pkg/apis/coordination"
	"k8s.io/kubernetes/pkg/apis/coordination/validation"
)

// leaseStrategy implements verification logic for Leases.
type leaseStrategy struct {
	runtime.ObjectTyper
	names.NameGenerator
}

// Strategy is the default logic that applies when creating and updating Lease objects.
var Strategy = leaseStrategy{legacyscheme.Scheme, names.SimpleNameGenerator}

// NamespaceScoped returns true because all Lease' need to be within a namespace.
func (leaseStrategy) NamespaceScoped() bool {
	return true
}

// PrepareForCreate prepares Lease for creation.
func (leaseStrategy) PrepareForCreate(ctx context.Context, obj runtime.Object) {
}

// PrepareForUpdate clears fields that are not allowed to be set by end users on update.
func (leaseStrategy) PrepareForUpdate(ctx context.Context, obj, old runtime.Object) {
}

// Validate validates a new Lease.
func (leaseStrategy) Validate(ctx context.Context, obj runtime.Object) field.ErrorList {
	lease := obj.(*coordination.Lease)
	return validation.ValidateLease(lease)
}

// Canonicalize normalizes the object after validation.
func (leaseStrategy) Canonicalize(obj runtime.Object) {
}

// AllowCreateOnUpdate is true for Lease; this means you may create one with a PUT request.
func (leaseStrategy) AllowCreateOnUpdate() bool {
	return true
}

// ValidateUpdate is the default update validation for an end user.
func (leaseStrategy) ValidateUpdate(ctx context.Context, obj, old runtime.Object) field.ErrorList {
	return validation.ValidateLeaseUpdate(obj.(*coordination.Lease), old.(*coordination.Lease))
}

// AllowUnconditionalUpdate is the default update policy for Lease objects.
func (leaseStrategy) AllowUnconditionalUpdate() bool {
	return false
}
