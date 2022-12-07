package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kubeadmapiv1beta2 "k8s.io/kubernetes/cmd/kubeadm/app/apis/kubeadm/v1beta2"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// BootstrapToken represents information for the bootstrap token output produced by kubeadm
// This is a copy of BoostrapToken struct from ../kubeadm/types.go with 2 additions:
// metav1.TypeMeta and metav1.ObjectMeta
type BootstrapToken struct {
	metav1.TypeMeta `json:",inline"`

	kubeadmapiv1beta2.BootstrapToken
}
