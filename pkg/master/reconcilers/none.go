// Package reconcilers a noop based reconciler
package reconcilers

import (
	"net"

	corev1 "k8s.io/api/core/v1"
)

// NoneEndpointReconciler allows for the endpoint reconciler to be disabled
type noneEndpointReconciler struct{}

// NewNoneEndpointReconciler creates a new EndpointReconciler that reconciles based on a
// nothing. It is a no-op.
func NewNoneEndpointReconciler() EndpointReconciler {
	return &noneEndpointReconciler{}
}

// ReconcileEndpoints noop reconcile
func (r *noneEndpointReconciler) ReconcileEndpoints(serviceName string, ip net.IP, endpointPorts []corev1.EndpointPort, reconcilePorts bool) error {
	return nil
}

// RemoveEndpoints noop reconcile
func (r *noneEndpointReconciler) RemoveEndpoints(serviceName string, ip net.IP, endpointPorts []corev1.EndpointPort) error {
	return nil
}

func (r *noneEndpointReconciler) StopReconciling() {
}
