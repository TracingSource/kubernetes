package testing

import (
	"net"

	"k8s.io/kubernetes/pkg/util/netsh"
)

// FakeNetsh is a no-op implementation of the netsh Interface
type FakeNetsh struct {
}

// NewFake returns a fakenetsh no-op implementation of the netsh Interface
func NewFake() *FakeNetsh {
	return &FakeNetsh{}
}

// EnsurePortProxyRule function implementing the netsh interface and always returns true and nil without any error
func (*FakeNetsh) EnsurePortProxyRule(args []string) (bool, error) {
	// Do Nothing
	return true, nil
}

// DeletePortProxyRule deletes the specified portproxy rule.  If the rule did not exist, return error.
func (*FakeNetsh) DeletePortProxyRule(args []string) error {
	// Do Nothing
	return nil
}

// EnsureIPAddress checks if the specified IP Address is added to vEthernet (HNSTransparent) interface, if not, add it.  If the address existed, return true.
func (*FakeNetsh) EnsureIPAddress(args []string, ip net.IP) (bool, error) {
	return true, nil
}

// DeleteIPAddress checks if the specified IP address is present and, if so, deletes it.
func (*FakeNetsh) DeleteIPAddress(args []string) error {
	// Do Nothing
	return nil
}

// Restore runs `netsh exec` to restore portproxy or addresses using a file.
// TODO Check if this is required, most likely not
func (*FakeNetsh) Restore(args []string) error {
	// Do Nothing
	return nil
}

// GetInterfaceToAddIP returns the interface name where Service IP needs to be added
// IP Address needs to be added for netsh portproxy to redirect traffic
// Reads Environment variable INTERFACE_TO_ADD_SERVICE_IP, if it is not defined then "vEthernet (HNSTransparent)" is returned
func (*FakeNetsh) GetInterfaceToAddIP() string {
	return "Interface 1"
}

var _ = netsh.Interface(&FakeNetsh{})
