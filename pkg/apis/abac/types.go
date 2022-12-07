package abac

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Policy contains a single ABAC policy rule
type Policy struct {
	metav1.TypeMeta

	// Spec describes the policy rule
	Spec PolicySpec
}

// PolicySpec contains the attributes for a policy rule
type PolicySpec struct {

	// User is the username this rule applies to.
	// Either user or group is required to match the request.
	// "*" matches all users.
	User string

	// Group is the group this rule applies to.
	// Either user or group is required to match the request.
	// "*" matches all groups.
	Group string

	// Readonly matches readonly requests when true, and all requests when false
	Readonly bool

	// APIGroup is the name of an API group. APIGroup, Resource, and Namespace are required to match resource requests.
	// "*" matches all API groups
	APIGroup string

	// Resource is the name of a resource. APIGroup, Resource, and Namespace are required to match resource requests.
	// "*" matches all resources
	Resource string

	// Namespace is the name of a namespace. APIGroup, Resource, and Namespace are required to match resource requests.
	// "*" matches all namespaces (including unnamespaced requests)
	Namespace string

	// NonResourcePath matches non-resource request paths.
	// "*" matches all paths
	// "/foo/*" matches all subpaths of foo
	NonResourcePath string

	// TODO: "expires" string in RFC3339 format.

	// TODO: want a way to allow some users to restart containers of a pod but
	// not delete or modify it.

	// TODO: want a way to allow a controller to create a pod based only on a
	// certain podTemplates.

}
