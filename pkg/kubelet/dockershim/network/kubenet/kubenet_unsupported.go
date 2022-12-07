// +build !linux

package kubenet

import (
	"fmt"

	kubeletconfig "k8s.io/kubernetes/pkg/kubelet/apis/config"
	kubecontainer "k8s.io/kubernetes/pkg/kubelet/container"
	"k8s.io/kubernetes/pkg/kubelet/dockershim/network"
)

type kubenetNetworkPlugin struct {
	network.NoopNetworkPlugin
}

func NewPlugin(networkPluginDirs []string, cacheDir string) network.NetworkPlugin {
	return &kubenetNetworkPlugin{}
}

func (plugin *kubenetNetworkPlugin) Init(host network.Host, hairpinMode kubeletconfig.HairpinMode, nonMasqueradeCIDR string, mtu int) error {
	return fmt.Errorf("Kubenet is not supported in this build")
}

func (plugin *kubenetNetworkPlugin) Name() string {
	return "kubenet"
}

func (plugin *kubenetNetworkPlugin) SetUpPod(namespace string, name string, id kubecontainer.ContainerID, annotations, options map[string]string) error {
	return fmt.Errorf("Kubenet is not supported in this build")
}

func (plugin *kubenetNetworkPlugin) TearDownPod(namespace string, name string, id kubecontainer.ContainerID) error {
	return fmt.Errorf("Kubenet is not supported in this build")
}

func (plugin *kubenetNetworkPlugin) GetPodNetworkStatus(namespace string, name string, id kubecontainer.ContainerID) (*network.PodNetworkStatus, error) {
	return nil, fmt.Errorf("Kubenet is not supported in this build")
}
