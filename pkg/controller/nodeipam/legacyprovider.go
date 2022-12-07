// +build !providerless

package nodeipam

import (
	"net"

	coreinformers "k8s.io/client-go/informers/core/v1"
	clientset "k8s.io/client-go/kubernetes"
	cloudprovider "k8s.io/cloud-provider"
	"k8s.io/klog"

	"k8s.io/kubernetes/pkg/controller/nodeipam/ipam"
	nodesync "k8s.io/kubernetes/pkg/controller/nodeipam/ipam/sync"
)

func startLegacyIPAM(
	ic *Controller,
	nodeInformer coreinformers.NodeInformer,
	cloud cloudprovider.Interface,
	kubeClient clientset.Interface,
	clusterCIDRs []*net.IPNet,
	serviceCIDR *net.IPNet,
	nodeCIDRMaskSizes []int,
) {
	cfg := &ipam.Config{
		Resync:       ipamResyncInterval,
		MaxBackoff:   ipamMaxBackoff,
		InitialRetry: ipamInitialBackoff,
	}
	switch ic.allocatorType {
	case ipam.IPAMFromClusterAllocatorType:
		cfg.Mode = nodesync.SyncFromCluster
	case ipam.IPAMFromCloudAllocatorType:
		cfg.Mode = nodesync.SyncFromCloud
	}

	// we may end up here with no cidr at all in case of FromCloud/FromCluster
	var cidr *net.IPNet
	if len(clusterCIDRs) > 0 {
		cidr = clusterCIDRs[0]
	}
	if len(clusterCIDRs) > 1 {
		klog.Warningf("Multiple cidrs were configured with FromCluster or FromCloud. cidrs except first one were discarded")
	}
	ipamc, err := ipam.NewController(cfg, kubeClient, cloud, cidr, serviceCIDR, nodeCIDRMaskSizes[0])
	if err != nil {
		klog.Fatalf("Error creating ipam controller: %v", err)
	}
	if err := ipamc.Start(nodeInformer); err != nil {
		klog.Fatalf("Error trying to Init(): %v", err)
	}
}
