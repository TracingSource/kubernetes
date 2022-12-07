// Package policybased implements a standard storage for Role that prevents privilege escalation.
package policybased

import (
	"context"

	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apiserver/pkg/authorization/authorizer"
	"k8s.io/apiserver/pkg/registry/rest"
	kapihelper "k8s.io/kubernetes/pkg/apis/core/helper"
	"k8s.io/kubernetes/pkg/apis/rbac"
	rbacregistry "k8s.io/kubernetes/pkg/registry/rbac"
	rbacregistryvalidation "k8s.io/kubernetes/pkg/registry/rbac/validation"
)

var groupResource = rbac.Resource("roles")

type Storage struct {
	rest.StandardStorage

	authorizer authorizer.Authorizer

	ruleResolver rbacregistryvalidation.AuthorizationRuleResolver
}

func NewStorage(s rest.StandardStorage, authorizer authorizer.Authorizer, ruleResolver rbacregistryvalidation.AuthorizationRuleResolver) *Storage {
	return &Storage{s, authorizer, ruleResolver}
}

func (r *Storage) NamespaceScoped() bool {
	return true
}

func (r *Storage) StorageVersion() runtime.GroupVersioner {
	svp, ok := r.StandardStorage.(rest.StorageVersionProvider)
	if !ok {
		return nil
	}
	return svp.StorageVersion()
}

var _ rest.StorageVersionProvider = &Storage{}

func (s *Storage) Create(ctx context.Context, obj runtime.Object, createValidation rest.ValidateObjectFunc, options *metav1.CreateOptions) (runtime.Object, error) {
	if rbacregistry.EscalationAllowed(ctx) || rbacregistry.RoleEscalationAuthorized(ctx, s.authorizer) {
		return s.StandardStorage.Create(ctx, obj, createValidation, options)
	}

	role := obj.(*rbac.Role)
	rules := role.Rules
	if err := rbacregistryvalidation.ConfirmNoEscalationInternal(ctx, s.ruleResolver, rules); err != nil {
		return nil, errors.NewForbidden(groupResource, role.Name, err)
	}
	return s.StandardStorage.Create(ctx, obj, createValidation, options)
}

func (s *Storage) Update(ctx context.Context, name string, obj rest.UpdatedObjectInfo, createValidation rest.ValidateObjectFunc, updateValidation rest.ValidateObjectUpdateFunc, forceAllowCreate bool, options *metav1.UpdateOptions) (runtime.Object, bool, error) {
	if rbacregistry.EscalationAllowed(ctx) || rbacregistry.RoleEscalationAuthorized(ctx, s.authorizer) {
		return s.StandardStorage.Update(ctx, name, obj, createValidation, updateValidation, forceAllowCreate, options)
	}

	nonEscalatingInfo := rest.WrapUpdatedObjectInfo(obj, func(ctx context.Context, obj runtime.Object, oldObj runtime.Object) (runtime.Object, error) {
		role := obj.(*rbac.Role)

		// if we're only mutating fields needed for the GC to eventually delete this obj, return
		if rbacregistry.IsOnlyMutatingGCFields(obj, oldObj, kapihelper.Semantic) {
			return obj, nil
		}

		rules := role.Rules
		if err := rbacregistryvalidation.ConfirmNoEscalationInternal(ctx, s.ruleResolver, rules); err != nil {
			return nil, errors.NewForbidden(groupResource, role.Name, err)
		}
		return obj, nil
	})

	return s.StandardStorage.Update(ctx, name, nonEscalatingInfo, createValidation, updateValidation, forceAllowCreate, options)
}
