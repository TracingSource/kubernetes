// +build windows

package v1beta2

const (
	// DefaultCACertPath defines default location of CA certificate on Windows
	DefaultCACertPath = "C:/etc/kubernetes/pki/ca.crt"
	// DefaultUrlScheme defines default socket url prefix
	DefaultUrlScheme = "npipe"
)
