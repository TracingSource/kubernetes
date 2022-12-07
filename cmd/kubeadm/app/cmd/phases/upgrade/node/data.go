package node

import (
	clientset "k8s.io/client-go/kubernetes"
	kubeadmapi "k8s.io/kubernetes/cmd/kubeadm/app/apis/kubeadm"
)

// Data is the interface to use for kubeadm upgrade node phases.
// The "nodeData" type from "cmd/upgrade/node.go" must satisfy this interface.
type Data interface {
	EtcdUpgrade() bool
	RenewCerts() bool
	DryRun() bool
	KubeletVersion() string
	Cfg() *kubeadmapi.InitConfiguration
	IsControlPlaneNode() bool
	Client() clientset.Interface
	KustomizeDir() string
}
