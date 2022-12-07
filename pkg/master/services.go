package master

import (
	"fmt"
	"net"

	"k8s.io/klog"
	"k8s.io/utils/integer"
	utilnet "k8s.io/utils/net"

	kubeoptions "k8s.io/kubernetes/pkg/kubeapiserver/options"
)

// ServiceIPRange 从 --service-cluster-ip-range 参数所指定的 service 网段中,
// 取第一个合法的 IP 作为 apiserver service(即 default/kubernetes 的 service) 的 IP.
//
// ServiceIPRange checks if the serviceClusterIPRange flag is nil, raising a warning if so and
// setting service ip range to the default value in kubeoptions.DefaultServiceIPCIDR
// for now until the default is removed per the deprecation timeline guidelines.
// Returns service ip range, api server service IP, and an error
func ServiceIPRange(passedServiceClusterIPRange net.IPNet) (net.IPNet, net.IP, error) {
	serviceClusterIPRange := passedServiceClusterIPRange
	if passedServiceClusterIPRange.IP == nil {
		klog.Warningf(
			"No CIDR for service cluster IPs specified. "+
			"Default value which was %s is deprecated and will be removed in future releases. "+
			"Please specify it using --service-cluster-ip-range on kube-apiserver.", 
			kubeoptions.DefaultServiceIPCIDR.String(),
		)
		serviceClusterIPRange = kubeoptions.DefaultServiceIPCIDR
	}

	size := integer.Int64Min(utilnet.RangeSize(&serviceClusterIPRange), 1<<16)
	if size < 8 {
		return net.IPNet{}, net.IP{}, fmt.Errorf(
			"The service cluster IP range must be at least %d IP addresses", 8,
		)
	}
	// 从指定的 service 网段中, 取出第一个合法的IP.
	// Select the first valid IP from ServiceClusterIPRange to use as the GenericAPIServer service IP.
	apiServerServiceIP, err := utilnet.GetIndexedIP(&serviceClusterIPRange, 1)
	if err != nil {
		return net.IPNet{}, net.IP{}, err
	}
	klog.V(4).Infof("Setting service IP to %q (read-write).", apiServerServiceIP)

	return serviceClusterIPRange, apiServerServiceIP, nil
}
