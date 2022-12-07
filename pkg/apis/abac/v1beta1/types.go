// +k8s:openapi-gen=true

package v1beta1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Policy contains a single ABAC policy rule
type Policy struct {
	metav1.TypeMeta `json:",inline"`

	// Spec describes the policy rule
	Spec PolicySpec `json:"spec"`
}

// PolicySpec contains the attributes for a policy rule
type PolicySpec struct {
	// User is the username this rule applies to.
	// Either user or group is required to match the request.
	// "*" matches all users.
	// +optional
	User string `json:"user,omitempty"`

	// Group is the group this rule applies to.
	// Either user or group is required to match the request.
	// "*" matches all groups.
	// +optional
	Group string `json:"group,omitempty"`

	// Readonly matches readonly requests when true, and all requests when false
	// +optional
	Readonly bool `json:"readonly,omitempty"`

	// APIGroup is the name of an API group. APIGroup, Resource, and Namespace are required to match resource requests.
	// "*" matches all API groups
	// +optional
	APIGroup string `json:"apiGroup,omitempty"`

	// Resource is the name of a resource. APIGroup, Resource, and Namespace are required to match resource requests.
	// "*" matches all resources
	// +optional
	Resource string `json:"resource,omitempty"`

	// Namespace is the name of a namespace. APIGroup, Resource, and Namespace are required to match resource requests.
	// "*" matches all namespaces (including unnamespaced requests)
	// +optional
	Namespace string `json:"namespace,omitempty"`

	// NonResourcePath matches non-resource request paths.
	// "*" matches all paths
	// "/foo/*" matches all subpaths of foo
	// +optional
	NonResourcePath string `json:"nonResourcePath,omitempty"`
}
