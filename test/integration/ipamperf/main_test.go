package ipamperf

import (
	"flag"
	"testing"

	"k8s.io/klog"
	"k8s.io/kubernetes/pkg/controller/nodeipam/ipam"
	"k8s.io/kubernetes/test/integration/framework"
)

var (
	resultsLogFile string
	isCustom       bool
	customConfig   = &Config{
		NumNodes:      10,
		KubeQPS:       30,
		CloudQPS:      30,
		CreateQPS:     100,
		AllocatorType: ipam.RangeAllocatorType,
	}
)

func TestMain(m *testing.M) {
	allocator := string(ipam.RangeAllocatorType)

	flag.StringVar(&resultsLogFile, "log", "", "log file to write JSON results to")
	flag.BoolVar(&isCustom, "custom", false, "enable custom test configuration")
	flag.StringVar(&allocator, "allocator", allocator, "allocator to use")
	flag.IntVar(&customConfig.KubeQPS, "kube-qps", customConfig.KubeQPS, "API server qps for allocations")
	flag.IntVar(&customConfig.NumNodes, "num-nodes", 10, "number of nodes to simulate")
	flag.IntVar(&customConfig.CreateQPS, "create-qps", customConfig.CreateQPS, "API server qps for node creation")
	flag.IntVar(&customConfig.CloudQPS, "cloud-qps", customConfig.CloudQPS, "GCE Cloud qps limit")
	flag.Parse()

	switch allocator {
	case string(ipam.RangeAllocatorType):
		customConfig.AllocatorType = ipam.RangeAllocatorType
	case string(ipam.CloudAllocatorType):
		customConfig.AllocatorType = ipam.CloudAllocatorType
	case string(ipam.IPAMFromCloudAllocatorType):
		customConfig.AllocatorType = ipam.IPAMFromCloudAllocatorType
	case string(ipam.IPAMFromClusterAllocatorType):
		customConfig.AllocatorType = ipam.IPAMFromClusterAllocatorType
	default:
		klog.Fatalf("Unknown allocator type: %s", allocator)
	}

	framework.EtcdMain(m.Run)
}
