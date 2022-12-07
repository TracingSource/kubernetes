package testing

import (
	"fmt"

	"k8s.io/kubernetes/pkg/kubelet/dockershim/network/hostport"
)

type fakeSyncer struct{}

func NewFakeHostportSyncer() hostport.HostportSyncer {
	return &fakeSyncer{}
}

func (h *fakeSyncer) OpenPodHostportsAndSync(newPortMapping *hostport.PodPortMapping, natInterfaceName string, activePortMapping []*hostport.PodPortMapping) error {
	return h.SyncHostports(natInterfaceName, activePortMapping)
}

func (h *fakeSyncer) SyncHostports(natInterfaceName string, activePortMapping []*hostport.PodPortMapping) error {
	for _, r := range activePortMapping {
		if r.IP.To4() == nil {
			return fmt.Errorf("invalid or missing pod %s/%s IP", r.Namespace, r.Name)
		}
	}

	return nil
}
