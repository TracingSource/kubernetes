package ipvs

import (
	"k8s.io/apimachinery/pkg/util/sets"
)

// 	@implementBy: pkg/proxy/ipvs/netlink_linux.go -> netlinkHandle{}
//
// NetLinkHandle for revoke netlink interface
type NetLinkHandle interface {
	// EnsureAddressBind checks if address is bound to the interface and, if not, binds it. 
	// If the address is already bound, return true.
	EnsureAddressBind(address, devName string) (exist bool, err error)
	// UnbindAddress unbind address from the interface
	UnbindAddress(address, devName string) error
	// EnsureDummyDevice checks if dummy device is exist and, if not, create one. 
	// If the dummy device is already exist, return true.
	EnsureDummyDevice(devName string) (exist bool, err error)
	// DeleteDummyDevice deletes the given dummy device by name.
	DeleteDummyDevice(devName string) error
	// ListBindAddress will list all IP addresses which are bound in a given interface
	ListBindAddress(devName string) ([]string, error)
	// GetLocalAddresses returns all unique local type IP addresses based on
	// specified device and filter device
	// If device is not specified, it will list all unique local type addresses
	// except filter device addresses
	GetLocalAddresses(dev, filterDev string) (sets.String, error)
}
