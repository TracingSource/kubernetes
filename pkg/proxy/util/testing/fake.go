package testing

import "net"

// FakeNetwork implements the NetworkInterfacer interface for test purpose.
type FakeNetwork struct {
	NetworkInterfaces []net.Interface
	// The key of map Addrs is the network interface name
	Address map[string][]net.Addr
}

// NewFakeNetwork initializes a FakeNetwork.
func NewFakeNetwork() *FakeNetwork {
	return &FakeNetwork{
		NetworkInterfaces: make([]net.Interface, 0),
		Address:           make(map[string][]net.Addr),
	}
}

// AddInterfaceAddr create an interface and its associated addresses for FakeNetwork implementation.
func (f *FakeNetwork) AddInterfaceAddr(intf *net.Interface, addrs []net.Addr) {
	f.NetworkInterfaces = append(f.NetworkInterfaces, *intf)
	f.Address[intf.Name] = addrs
}

// Addrs is part of NetworkInterfacer interface.
func (f *FakeNetwork) Addrs(intf *net.Interface) ([]net.Addr, error) {
	return f.Address[intf.Name], nil
}

// Interfaces is part of NetworkInterfacer interface.
func (f *FakeNetwork) Interfaces() ([]net.Interface, error) {
	return f.NetworkInterfaces, nil
}

// AddrStruct implements the net.Addr for test purpose.
type AddrStruct struct{ Val string }

// Network is part of net.Addr interface.
func (a AddrStruct) Network() string {
	return a.Val
}

// String is part of net.Addr interface.
func (a AddrStruct) String() string {
	return a.Val
}

var _ net.Addr = &AddrStruct{}
