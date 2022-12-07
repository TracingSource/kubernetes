// +build !providerless

package gce

const (
	// AlphaFeatureNetworkTiers allows Services backed by a GCP load balancer to choose
	// what network tier to use. Currently supports "Standard" and "Premium" (default).
	//
	// alpha: v1.8 (for Services)
	AlphaFeatureNetworkTiers = "NetworkTiers"
	// AlphaFeatureILBSubsets allows InternalLoadBalancer services to include a subset
	// of cluster nodes as backends instead of all nodes.
	AlphaFeatureILBSubsets = "ILBSubsets"
	// AlphaFeatureILBCustomSubnet allows InternalLoadBalancer services to specify a
	// network subnet to allocate ip addresses from.
	AlphaFeatureILBCustomSubnet = "ILBCustomSubnet"
)

// AlphaFeatureGate contains a mapping of alpha features to whether they are enabled
type AlphaFeatureGate struct {
	features map[string]bool
}

// Enabled returns true if the provided alpha feature is enabled
func (af *AlphaFeatureGate) Enabled(key string) bool {
	if af == nil || af.features == nil {
		return false
	}
	return af.features[key]
}

// NewAlphaFeatureGate marks the provided alpha features as enabled
func NewAlphaFeatureGate(features []string) *AlphaFeatureGate {
	featureMap := make(map[string]bool)
	for _, name := range features {
		featureMap[name] = true
	}
	return &AlphaFeatureGate{featureMap}
}
