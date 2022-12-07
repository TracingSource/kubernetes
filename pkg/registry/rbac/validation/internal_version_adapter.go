package validation

import (
	"context"

	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/kubernetes/pkg/apis/rbac"
	rbacv1helpers "k8s.io/kubernetes/pkg/apis/rbac/v1"
)

func ConfirmNoEscalationInternal(ctx context.Context, ruleResolver AuthorizationRuleResolver, inRules []rbac.PolicyRule) error {
	rules := []rbacv1.PolicyRule{}
	for i := range inRules {
		v1Rule := rbacv1.PolicyRule{}
		err := rbacv1helpers.Convert_rbac_PolicyRule_To_v1_PolicyRule(&inRules[i], &v1Rule, nil)
		if err != nil {
			return err
		}
		rules = append(rules, v1Rule)
	}

	return ConfirmNoEscalation(ctx, ruleResolver, rules)
}
