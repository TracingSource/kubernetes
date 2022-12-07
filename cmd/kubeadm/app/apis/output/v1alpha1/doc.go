// +k8s:defaulter-gen=TypeMeta
// +groupName=output.kubeadm.k8s.io
// +k8s:deepcopy-gen=package
// +k8s:conversion-gen=k8s.io/kubernetes/cmd/kubeadm/app/apis/output

// Package v1alpha1 defines the v1alpha1 version of the kubeadm data structures
// related to structured output
// The purpose of the kubeadm structured output is to have a well
// defined versioned output format that other software that uses
// kubeadm for cluster deployments can use and rely on.
package v1alpha1 // import "k8s.io/kubernetes/cmd/kubeadm/app/apis/output/v1alpha1"
