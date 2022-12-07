package cloud

import (
	"strings"
)

// NetworkTier represents the Network Service Tier used by a resource
type NetworkTier string

// LbScheme represents the possible types of load balancers
type LbScheme string

const (
	NetworkTierStandard NetworkTier = "Standard"
	NetworkTierPremium  NetworkTier = "Premium"
	NetworkTierDefault  NetworkTier = NetworkTierPremium

	SchemeExternal LbScheme = "EXTERNAL"
	SchemeInternal LbScheme = "INTERNAL"
)

// ToGCEValue converts NetworkTier to a string that we can populate the
// NetworkTier field of GCE objects, including ForwardingRules and Addresses.
func (n NetworkTier) ToGCEValue() string {
	return strings.ToUpper(string(n))
}

// NetworkTierGCEValueToType converts the value of the NetworkTier field of a
// GCE object to the NetworkTier type.
func NetworkTierGCEValueToType(s string) NetworkTier {
	switch s {
	case NetworkTierStandard.ToGCEValue():
		return NetworkTierStandard
	case NetworkTierPremium.ToGCEValue():
		return NetworkTierPremium
	default:
		return NetworkTier(s)
	}
}
