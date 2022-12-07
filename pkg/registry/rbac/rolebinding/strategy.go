package rolebinding

import (
	"context"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"k8s.io/apiserver/pkg/registry/rest"
	"k8s.io/apiserver/pkg/storage/names"
	"k8s.io/kubernetes/pkg/api/legacyscheme"
	"k8s.io/kubernetes/pkg/apis/rbac"
	"k8s.io/kubernetes/pkg/apis/rbac/validation"
)

// strategy implements behavior for RoleBindings
type strategy struct {
	runtime.ObjectTyper
	names.NameGenerator
}

// strategy is the default logic that applies when creating and updating
// RoleBinding objects.
var Strategy = strategy{legacyscheme.Scheme, names.SimpleNameGenerator}

// Strategy should implement rest.RESTCreateStrategy
var _ rest.RESTCreateStrategy = Strategy

// Strategy should implement rest.RESTUpdateStrategy
var _ rest.RESTUpdateStrategy = Strategy

// NamespaceScoped is true for RoleBindings.
func (strategy) NamespaceScoped() bool {
	return true
}

// AllowCreateOnUpdate is true for RoleBindings.
func (strategy) AllowCreateOnUpdate() bool {
	return true
}

// PrepareForCreate clears fields that are not allowed to be set by end users
// on creation.
func (strategy) PrepareForCreate(ctx context.Context, obj runtime.Object) {
	_ = obj.(*rbac.RoleBinding)
}

// PrepareForUpdate clears fields that are not allowed to be set by end users on update.
func (strategy) PrepareForUpdate(ctx context.Context, obj, old runtime.Object) {
	newRoleBinding := obj.(*rbac.RoleBinding)
	oldRoleBinding := old.(*rbac.RoleBinding)

	_, _ = newRoleBinding, oldRoleBinding
}

// Validate validates a new RoleBinding. Validation must check for a correct signature.
func (strategy) Validate(ctx context.Context, obj runtime.Object) field.ErrorList {
	roleBinding := obj.(*rbac.RoleBinding)
	return validation.ValidateRoleBinding(roleBinding)
}

// Canonicalize normalizes the object after validation.
func (strategy) Canonicalize(obj runtime.Object) {
	_ = obj.(*rbac.RoleBinding)
}

// ValidateUpdate is the default update validation for an end user.
func (strategy) ValidateUpdate(ctx context.Context, obj, old runtime.Object) field.ErrorList {
	newObj := obj.(*rbac.RoleBinding)
	errorList := validation.ValidateRoleBinding(newObj)
	return append(errorList, validation.ValidateRoleBindingUpdate(newObj, old.(*rbac.RoleBinding))...)
}

// If AllowUnconditionalUpdate() is true and the object specified by
// the user does not have a resource version, then generic Update()
// populates it with the latest version. Else, it checks that the
// version specified by the user matches the version of latest etcd
// object.
func (strategy) AllowUnconditionalUpdate() bool {
	return true
}
