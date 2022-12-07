package kubemark

import (
	"fmt"
	"net"
	"time"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	clientset "k8s.io/client-go/kubernetes"
	v1core "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/tools/record"
	proxyapp "k8s.io/kubernetes/cmd/kube-proxy/app"
	"k8s.io/kubernetes/pkg/proxy"
	proxyconfig "k8s.io/kubernetes/pkg/proxy/config"
	"k8s.io/kubernetes/pkg/proxy/iptables"
	utiliptables "k8s.io/kubernetes/pkg/util/iptables"
	utilnode "k8s.io/kubernetes/pkg/util/node"
	utilsysctl "k8s.io/kubernetes/pkg/util/sysctl"
	utilexec "k8s.io/utils/exec"
	utilpointer "k8s.io/utils/pointer"

	"k8s.io/klog"
)

type HollowProxy struct {
	ProxyServer *proxyapp.ProxyServer
}

type FakeProxier struct {
	proxyconfig.NoopEndpointSliceHandler
	proxyconfig.NoopNodeHandler
}

func (*FakeProxier) Sync() {}
func (*FakeProxier) SyncLoop() {
	select {}
}
func (*FakeProxier) OnServiceAdd(service *v1.Service)                        {}
func (*FakeProxier) OnServiceUpdate(oldService, service *v1.Service)         {}
func (*FakeProxier) OnServiceDelete(service *v1.Service)                     {}
func (*FakeProxier) OnServiceSynced()                                        {}
func (*FakeProxier) OnEndpointsAdd(endpoints *v1.Endpoints)                  {}
func (*FakeProxier) OnEndpointsUpdate(oldEndpoints, endpoints *v1.Endpoints) {}
func (*FakeProxier) OnEndpointsDelete(endpoints *v1.Endpoints)               {}
func (*FakeProxier) OnEndpointsSynced()                                      {}

func NewHollowProxyOrDie(
	nodeName string,
	client clientset.Interface,
	eventClient v1core.EventsGetter,
	iptInterface utiliptables.Interface,
	sysctl utilsysctl.Interface,
	execer utilexec.Interface,
	broadcaster record.EventBroadcaster,
	recorder record.EventRecorder,
	useRealProxier bool,
	proxierSyncPeriod time.Duration,
	proxierMinSyncPeriod time.Duration,
) (*HollowProxy, error) {
	// Create proxier and service/endpoint handlers.
	var proxier proxy.Provider
	var err error

	if useRealProxier {
		nodeIP := utilnode.GetNodeIP(client, nodeName)
		if nodeIP == nil {
			klog.V(0).Infof("can't determine this node's IP, assuming 127.0.0.1")
			nodeIP = net.ParseIP("127.0.0.1")
		}
		// Real proxier with fake iptables, sysctl, etc underneath it.
		//var err error
		proxier, err = iptables.NewProxier(
			iptInterface,
			sysctl,
			execer,
			proxierSyncPeriod,
			proxierMinSyncPeriod,
			false,
			0,
			"10.0.0.0/8",
			nodeName,
			nodeIP,
			recorder,
			nil,
			[]string{},
		)
		if err != nil {
			return nil, fmt.Errorf("unable to create proxier: %v", err)
		}
	} else {
		proxier = &FakeProxier{}
	}

	// Create a Hollow Proxy instance.
	nodeRef := &v1.ObjectReference{
		Kind:      "Node",
		Name:      nodeName,
		UID:       types.UID(nodeName),
		Namespace: "",
	}
	return &HollowProxy{
		ProxyServer: &proxyapp.ProxyServer{
			Client:           client,
			EventClient:      eventClient,
			IptInterface:     iptInterface,
			Proxier:          proxier,
			Broadcaster:      broadcaster,
			Recorder:         recorder,
			ProxyMode:        "fake",
			NodeRef:          nodeRef,
			OOMScoreAdj:      utilpointer.Int32Ptr(0),
			ConfigSyncPeriod: 30 * time.Second,
		},
	}, nil
}

func (hp *HollowProxy) Run() {
	if err := hp.ProxyServer.Run(); err != nil {
		klog.Fatalf("Error while running proxy: %v\n", err)
	}
}
