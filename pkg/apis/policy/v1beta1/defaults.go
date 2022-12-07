package v1beta1

import (
	policyv1beta1 "k8s.io/api/policy/v1beta1"
)

func SetDefaults_PodSecurityPolicySpec(obj *policyv1beta1.PodSecurityPolicySpec) {
	// This field was added after PodSecurityPolicy was released.
	// Policies that do not include this field must remain as permissive as they were prior to the introduction of this field.
	if obj.AllowPrivilegeEscalation == nil {
		t := true
		obj.AllowPrivilegeEscalation = &t
	}
}
