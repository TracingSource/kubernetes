package config

// NodeIPAMControllerConfiguration contains elements describing NodeIPAMController.
type NodeIPAMControllerConfiguration struct {
	// serviceCIDR is CIDR Range for Services in cluster.
	ServiceCIDR string
	// secondaryServiceCIDR is CIDR Range for Services in cluster. This is used in dual stack clusters. SecondaryServiceCIDR must be of different IP family than ServiceCIDR
	SecondaryServiceCIDR string
	// NodeCIDRMaskSize is the mask size for node cidr in single-stack cluster.
	// This can be used only with single stack clusters and is incompatible with dual stack clusters.
	NodeCIDRMaskSize int32
	// NodeCIDRMaskSizeIPv4 is the mask size for IPv4 node cidr in dual-stack cluster.
	// This can be used only with dual stack clusters and is incompatible with single stack clusters.
	NodeCIDRMaskSizeIPv4 int32
	// NodeCIDRMaskSizeIPv6 is the mask size for IPv6 node cidr in dual-stack cluster.
	// This can be used only with dual stack clusters and is incompatible with single stack clusters.
	NodeCIDRMaskSizeIPv6 int32
}
