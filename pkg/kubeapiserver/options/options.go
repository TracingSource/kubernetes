package options

import (
	"net"

	utilnet "k8s.io/apimachinery/pkg/util/net"
)

// DefaultServiceNodePortRange is the default port range for NodePort services.
var DefaultServiceNodePortRange = utilnet.PortRange{Base: 30000, Size: 2768}

// DefaultServiceIPCIDR is a CIDR notation of IP range from which to allocate service cluster IPs
var DefaultServiceIPCIDR net.IPNet = net.IPNet{IP: net.ParseIP("10.0.0.0"), Mask: net.CIDRMask(24, 32)}

const DefaultEtcdPathPrefix = "/registry"
