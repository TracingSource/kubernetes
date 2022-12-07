// Package policybased implements a standard storage for ClusterRole that prevents privilege escalation.
package policybased

import (
	"context"
	"errors"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apiserver/pkg/authorization/authorizer"
	"k8s.io/apiserver/pkg/registry/rest"
	kapihelper "k8s.io/kubernetes/pkg/apis/core/helper"
	"k8s.io/kubernetes/pkg/apis/rbac"
	rbacregistry "k8s.io/kubernetes/pkg/registry/rbac"
	rbacregistryvalidation "k8s.io/kubernetes/pkg/registry/rbac/validation"
)

var groupResource = rbac.Resource("clusterroles")

type Storage struct {
	rest.StandardStorage

	authorizer authorizer.Authorizer

	ruleResolver rbacregistryvalidation.AuthorizationRuleResolver
}

func NewStorage(s rest.StandardStorage, authorizer authorizer.Authorizer, ruleResolver rbacregistryvalidation.AuthorizationRuleResolver) *Storage {
	return &Storage{s, authorizer, ruleResolver}
}

func (r *Storage) NamespaceScoped() bool {
	return false
}

func (r *Storage) StorageVersion() runtime.GroupVersioner {
	svp, ok := r.StandardStorage.(rest.StorageVersionProvider)
	if !ok {
		return nil
	}
	return svp.StorageVersion()
}

var _ rest.StorageVersionProvider = &Storage{}

var fullAuthority = []rbac.PolicyRule{
	rbac.NewRule("*").Groups("*").Resources("*").RuleOrDie(),
	rbac.NewRule("*").URLs("*").RuleOrDie(),
}

func (s *Storage) Create(ctx context.Context, obj runtime.Object, createValidatingAdmission rest.ValidateObjectFunc, options *metav1.CreateOptions) (runtime.Object, error) {
	if rbacregistry.EscalationAllowed(ctx) || rbacregistry.RoleEscalationAuthorized(ctx, s.authorizer) {
		return s.StandardStorage.Create(ctx, obj, createValidatingAdmission, options)
	}

	clusterRole := obj.(*rbac.ClusterRole)
	rules := clusterRole.Rules
	if err := rbacregistryvalidation.ConfirmNoEscalationInternal(ctx, s.ruleResolver, rules); err != nil {
		return nil, apierrors.NewForbidden(groupResource, clusterRole.Name, err)
	}
	// to set the aggregation rule, since it can gather anything, requires * on *.*
	if hasAggregationRule(clusterRole) {
		if err := rbacregistryvalidation.ConfirmNoEscalationInternal(ctx, s.ruleResolver, fullAuthority); err != nil {
			return nil, apierrors.NewForbidden(groupResource, clusterRole.Name, errors.New("must have cluster-admin privileges to use the aggregationRule"))
		}
	}

	return s.StandardStorage.Create(ctx, obj, createValidatingAdmission, options)
}

func (s *Storage) Update(ctx context.Context, name string, obj rest.UpdatedObjectInfo, createValidation rest.ValidateObjectFunc, updateValidation rest.ValidateObjectUpdateFunc, forceAllowCreate bool, options *metav1.UpdateOptions) (runtime.Object, bool, error) {
	if rbacregistry.EscalationAllowed(ctx) || rbacregistry.RoleEscalationAuthorized(ctx, s.authorizer) {
		return s.StandardStorage.Update(ctx, name, obj, createValidation, updateValidation, forceAllowCreate, options)
	}

	nonEscalatingInfo := rest.WrapUpdatedObjectInfo(obj, func(ctx context.Context, obj runtime.Object, oldObj runtime.Object) (runtime.Object, error) {
		clusterRole := obj.(*rbac.ClusterRole)
		oldClusterRole := oldObj.(*rbac.ClusterRole)

		// if we're only mutating fields needed for the GC to eventually delete this obj, return
		if rbacregistry.IsOnlyMutatingGCFields(obj, oldObj, kapihelper.Semantic) {
			return obj, nil
		}

		rules := clusterRole.Rules
		if err := rbacregistryvalidation.ConfirmNoEscalationInternal(ctx, s.ruleResolver, rules); err != nil {
			return nil, apierrors.NewForbidden(groupResource, clusterRole.Name, err)
		}
		// to change the aggregation rule, since it can gather anything and prevent tightening, requires * on *.*
		if hasAggregationRule(clusterRole) || hasAggregationRule(oldClusterRole) {
			if err := rbacregistryvalidation.ConfirmNoEscalationInternal(ctx, s.ruleResolver, fullAuthority); err != nil {
				return nil, apierrors.NewForbidden(groupResource, clusterRole.Name, errors.New("must have cluster-admin privileges to use the aggregationRule"))
			}
		}

		return obj, nil
	})

	return s.StandardStorage.Update(ctx, name, nonEscalatingInfo, createValidation, updateValidation, forceAllowCreate, options)
}

func hasAggregationRule(clusterRole *rbac.ClusterRole) bool {
	return clusterRole.AggregationRule != nil && len(clusterRole.AggregationRule.ClusterRoleSelectors) > 0
}
