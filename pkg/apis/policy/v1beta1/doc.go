// +k8s:conversion-gen=k8s.io/kubernetes/pkg/apis/policy
// +k8s:conversion-gen-external-types=k8s.io/api/policy/v1beta1
// +k8s:defaulter-gen=TypeMeta
// +k8s:defaulter-gen-input=../../../../vendor/k8s.io/api/policy/v1beta1

// Package policy is for any kind of policy object.  Suitable examples, even if
// they aren't all here, are policyv1beta1.PodDisruptionBudget, PodSecurityPolicy,
// NetworkPolicy, etc.
package v1beta1 // import "k8s.io/kubernetes/pkg/apis/policy/v1beta1"
