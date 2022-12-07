package config

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kubectrlmgrconfig "k8s.io/kubernetes/pkg/controller/apis/config"
	serviceconfig "k8s.io/kubernetes/pkg/controller/service/config"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// CloudControllerManagerConfiguration contains elements describing cloud-controller manager.
type CloudControllerManagerConfiguration struct {
	metav1.TypeMeta

	// Generic holds configuration for a generic controller-manager
	Generic kubectrlmgrconfig.GenericControllerManagerConfiguration
	// KubeCloudSharedConfiguration holds configuration for shared related features
	// both in cloud controller manager and kube-controller manager.
	KubeCloudShared kubectrlmgrconfig.KubeCloudSharedConfiguration

	// ServiceControllerConfiguration holds configuration for ServiceController
	// related features.
	ServiceController serviceconfig.ServiceControllerConfiguration

	// NodeStatusUpdateFrequency is the frequency at which the controller updates nodes' status
	NodeStatusUpdateFrequency metav1.Duration
}
