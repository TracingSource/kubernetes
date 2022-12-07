// +build !windows

package cni

import (
	"fmt"

	"github.com/containernetworking/cni/libcni"
	runtimeapi "k8s.io/cri-api/pkg/apis/runtime/v1alpha2"
	kubecontainer "k8s.io/kubernetes/pkg/kubelet/container"
	"k8s.io/kubernetes/pkg/kubelet/dockershim/network"
)

func getLoNetwork(binDirs []string) *cniNetwork {
	loConfig, err := libcni.ConfListFromBytes([]byte(`{
  "cniVersion": "0.2.0",
  "name": "cni-loopback",
  "plugins":[{
    "type": "loopback"
  }]
}`))
	if err != nil {
		// The hardcoded config above should always be valid and unit tests will
		// catch this
		panic(err)
	}
	loNetwork := &cniNetwork{
		name:          "lo",
		NetworkConfig: loConfig,
		CNIConfig:     &libcni.CNIConfig{Path: binDirs},
	}

	return loNetwork
}

func (plugin *cniNetworkPlugin) platformInit() error {
	var err error
	plugin.nsenterPath, err = plugin.execer.LookPath("nsenter")
	if err != nil {
		return err
	}
	return nil
}

// TODO: Use the addToNetwork function to obtain the IP of the Pod. That will assume idempotent ADD call to the plugin.
// Also fix the runtime's call to Status function to be done only in the case that the IP is lost, no need to do periodic calls
func (plugin *cniNetworkPlugin) GetPodNetworkStatus(namespace string, name string, id kubecontainer.ContainerID) (*network.PodNetworkStatus, error) {
	netnsPath, err := plugin.host.GetNetNS(id.ID)
	if err != nil {
		return nil, fmt.Errorf("CNI failed to retrieve network namespace path: %v", err)
	}
	if netnsPath == "" {
		return nil, fmt.Errorf("cannot find the network namespace, skipping pod network status for container %q", id)
	}

	ips, err := network.GetPodIPs(plugin.execer, plugin.nsenterPath, netnsPath, network.DefaultInterfaceName)
	if err != nil {
		return nil, err
	}

	if len(ips) == 0 {
		return nil, fmt.Errorf("cannot find pod IPs in the network namespace, skipping pod network status for container %q", id)
	}

	return &network.PodNetworkStatus{
		IP:  ips[0],
		IPs: ips,
	}, nil
}

// buildDNSCapabilities builds cniDNSConfig from runtimeapi.DNSConfig.
func buildDNSCapabilities(dnsConfig *runtimeapi.DNSConfig) *cniDNSConfig {
	return nil
}
