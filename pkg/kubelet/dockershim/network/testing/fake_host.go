package testing

// helper for testing plugins
// a fake host is created here that can be used by plugins for testing

import (
	"k8s.io/api/core/v1"
	clientset "k8s.io/client-go/kubernetes"
	"k8s.io/kubernetes/pkg/kubelet/dockershim/network/hostport"
)

type fakeNetworkHost struct {
	fakeNamespaceGetter
	FakePortMappingGetter
	kubeClient clientset.Interface
	Legacy     bool
}

func NewFakeHost(kubeClient clientset.Interface) *fakeNetworkHost {
	host := &fakeNetworkHost{kubeClient: kubeClient, Legacy: true}
	return host
}

func (fnh *fakeNetworkHost) GetPodByName(name, namespace string) (*v1.Pod, bool) {
	return nil, false
}

func (fnh *fakeNetworkHost) GetKubeClient() clientset.Interface {
	return nil
}

func (nh *fakeNetworkHost) SupportsLegacyFeatures() bool {
	return nh.Legacy
}

type fakeNamespaceGetter struct {
	ns string
}

func (nh *fakeNamespaceGetter) GetNetNS(containerID string) (string, error) {
	return nh.ns, nil
}

type FakePortMappingGetter struct {
	PortMaps map[string][]*hostport.PortMapping
}

func (pm *FakePortMappingGetter) GetPodPortMappings(containerID string) ([]*hostport.PortMapping, error) {
	return pm.PortMaps[containerID], nil
}
