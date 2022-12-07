package settings

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	api "k8s.io/kubernetes/pkg/apis/core"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// PodPreset is a policy resource that defines additional runtime
// requirements for a Pod.
type PodPreset struct {
	metav1.TypeMeta
	// +optional
	metav1.ObjectMeta

	// +optional
	Spec PodPresetSpec
}

// PodPresetSpec is a description of a pod preset.
type PodPresetSpec struct {
	// Selector is a label query over a set of resources, in this case pods.
	// Required.
	Selector metav1.LabelSelector
	// Env defines the collection of EnvVar to inject into containers.
	// +optional
	Env []api.EnvVar
	// EnvFrom defines the collection of EnvFromSource to inject into containers.
	// +optional
	EnvFrom []api.EnvFromSource
	// Volumes defines the collection of Volume to inject into the pod.
	// +optional
	Volumes []api.Volume
	// VolumeMounts defines the collection of VolumeMount to inject into containers.
	// +optional
	VolumeMounts []api.VolumeMount
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// PodPresetList is a list of PodPreset objects.
type PodPresetList struct {
	metav1.TypeMeta
	// +optional
	metav1.ListMeta

	Items []PodPreset
}
