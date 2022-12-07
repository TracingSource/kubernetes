package output

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kubeadmapiv1beta2 "k8s.io/kubernetes/cmd/kubeadm/app/apis/kubeadm/v1beta2"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// BootstrapToken represents information for the output produced by 'kubeadm token list'
// This is a copy of BoostrapToken struct from ../kubeadm/types.go with 2 additions:
// metav1.TypeMeta and metav1.ObjectMeta
type BootstrapToken struct {
	metav1.TypeMeta

	kubeadmapiv1beta2.BootstrapToken
}
