package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kubectrlmgrconfigv1alpha1 "k8s.io/kube-controller-manager/config/v1alpha1"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type CloudControllerManagerConfiguration struct {
	metav1.TypeMeta `json:",inline"`

	// Generic holds configuration for a generic controller-manager
	Generic kubectrlmgrconfigv1alpha1.GenericControllerManagerConfiguration
	// KubeCloudSharedConfiguration holds configuration for shared related features
	// both in cloud controller manager and kube-controller manager.
	KubeCloudShared kubectrlmgrconfigv1alpha1.KubeCloudSharedConfiguration
	// ServiceControllerConfiguration holds configuration for ServiceController
	// related features.
	ServiceController kubectrlmgrconfigv1alpha1.ServiceControllerConfiguration
	// NodeStatusUpdateFrequency is the frequency at which the controller updates nodes' status
	NodeStatusUpdateFrequency metav1.Duration
}
