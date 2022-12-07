package userspace

import (
	"k8s.io/api/core/v1"
	"k8s.io/kubernetes/pkg/proxy"
	proxyconfig "k8s.io/kubernetes/pkg/proxy/config"
	"net"
)

// LoadBalancer is an interface for distributing incoming requests to service endpoints.
type LoadBalancer interface {
	// NextEndpoint returns the endpoint to handle a request for the given
	// service-port and source address.
	NextEndpoint(service proxy.ServicePortName, srcAddr net.Addr, sessionAffinityReset bool) (string, error)
	NewService(service proxy.ServicePortName, sessionAffinityType v1.ServiceAffinity, stickyMaxAgeSeconds int) error
	DeleteService(service proxy.ServicePortName)
	CleanupStaleStickySessions(service proxy.ServicePortName)
	ServiceHasEndpoints(service proxy.ServicePortName) bool

	proxyconfig.EndpointsHandler
}
