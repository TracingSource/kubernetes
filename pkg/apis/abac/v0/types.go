// +k8s:openapi-gen=true

package v0

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Policy contains a single ABAC policy rule
type Policy struct {
	metav1.TypeMeta `json:",inline"`

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

	// Resource is the name of a resource
	// "*" matches all resources
	// +optional
	Resource string `json:"resource,omitempty"`

	// Namespace is the name of a namespace
	// "*" matches all namespaces (including unnamespaced requests)
	// +optional
	Namespace string `json:"namespace,omitempty"`
}
