package podtolerationrestriction

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	api "k8s.io/kubernetes/pkg/apis/core"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Configuration provides configuration for the PodTolerationRestriction admission controller.
type Configuration struct {
	metav1.TypeMeta

	// cluster level default tolerations
	Default []api.Toleration

	// cluster level whitelist of tolerations
	Whitelist []api.Toleration
}
