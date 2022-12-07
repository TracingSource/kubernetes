package bandwidth

import "k8s.io/apimachinery/pkg/api/resource"

// Shaper is designed so that the shaper structs created
// satisfy the Shaper interface.
type Shaper interface {
	// Limit the bandwidth for a particular CIDR on a particular interface
	//   * ingress and egress are in bits/second
	//   * cidr is expected to be a valid network CIDR (e.g. '1.2.3.4/32' or '10.20.0.1/16')
	// 'egress' bandwidth limit applies to all packets on the interface whose source matches 'cidr'
	// 'ingress' bandwidth limit applies to all packets on the interface whose destination matches 'cidr'
	// Limits are aggregate limits for the CIDR, not per IP address.  CIDRs must be unique, but can be overlapping, traffic
	// that matches multiple CIDRs counts against all limits.
	Limit(cidr string, egress, ingress *resource.Quantity) error
	// Remove a bandwidth limit for a particular CIDR on a particular network interface
	Reset(cidr string) error
	// Reconcile the interface managed by this shaper with the state on the ground.
	ReconcileInterface() error
	// Reconcile a CIDR managed by this shaper with the state on the ground
	ReconcileCIDR(cidr string, egress, ingress *resource.Quantity) error
	// GetCIDRs returns the set of CIDRs that are being managed by this shaper
	GetCIDRs() ([]string, error)
}
