package v1alpha1

import (
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Configuration provides configuration for the PodTolerationRestriction admission controller.
type Configuration struct {
	metav1.TypeMeta `json:",inline"`

	// cluster level default tolerations
	Default []v1.Toleration `json:"default,omitempty"`

	// cluster level whitelist of tolerations
	Whitelist []v1.Toleration `json:"whitelist,omitempty"`
}
