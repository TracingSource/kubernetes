package util

import (
	"fmt"
	"net"
	"strconv"

	"k8s.io/klog"
)

// IPPart returns just the IP part of an IP or IP:port or endpoint string. If the IP
// part is an IPv6 address enclosed in brackets (e.g. "[fd00:1::5]:9999"),
// then the brackets are stripped as well.
func IPPart(s string) string {
	if ip := net.ParseIP(s); ip != nil {
		// IP address without port
		return s
	}
	// Must be IP:port
	host, _, err := net.SplitHostPort(s)
	if err != nil {
		klog.Errorf("Error parsing '%s': %v", s, err)
		return ""
	}
	// Check if host string is a valid IP address
	ip := net.ParseIP(host)
	if ip == nil {
		klog.Errorf("invalid IP part '%s'", host)
		return ""
	}
	return ip.String()
}

// PortPart returns just the port part of an endpoint string.
func PortPart(s string) (int, error) {
	// Must be IP:port
	_, port, err := net.SplitHostPort(s)
	if err != nil {
		klog.Errorf("Error parsing '%s': %v", s, err)
		return -1, err
	}
	portNumber, err := strconv.Atoi(port)
	if err != nil {
		klog.Errorf("Error parsing '%s': %v", port, err)
		return -1, err
	}
	return portNumber, nil
}

// ToCIDR returns a host address of the form <ip-address>/32 for
// IPv4 and <ip-address>/128 for IPv6
func ToCIDR(ip net.IP) string {
	len := 32
	if ip.To4() == nil {
		len = 128
	}
	return fmt.Sprintf("%s/%d", ip.String(), len)
}
