// +build !freebsd,!linux,!windows,!darwin

package util

import (
	"fmt"
	"net"
	"time"
)

// CreateListener creates a listener on the specified endpoint.
func CreateListener(endpoint string) (net.Listener, error) {
	return nil, fmt.Errorf("CreateListener is unsupported in this build")
}

// GetAddressAndDialer returns the address parsed from the given endpoint and a dialer.
func GetAddressAndDialer(endpoint string) (string, func(addr string, timeout time.Duration) (net.Conn, error), error) {
	return "", nil, fmt.Errorf("GetAddressAndDialer is unsupported in this build")
}

// LockAndCheckSubPath empty implementation
func LockAndCheckSubPath(volumePath, subPath string) ([]uintptr, error) {
	return []uintptr{}, nil
}

// UnlockPath empty implementation
func UnlockPath(fileHandles []uintptr) {
}

// LocalEndpoint empty implementation
func LocalEndpoint(path, file string) (string, error) {
	return "", fmt.Errorf("LocalEndpoints are unsupported in this build")
}

// GetBootTime empty implementation
func GetBootTime() (time.Time, error) {
	return time.Time{}, fmt.Errorf("GetBootTime is unsupported in this build")
}
