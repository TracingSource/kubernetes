// +build !providerless

package openstack

import (
	"sync"

	"k8s.io/component-base/metrics"
	"k8s.io/component-base/metrics/legacyregistry"
)

const (
	openstackSubsystem         = "openstack"
	openstackOperationKey      = "cloudprovider_openstack_api_request_duration_seconds"
	openstackOperationErrorKey = "cloudprovider_openstack_api_request_errors"
)

var (
	openstackOperationsLatency = metrics.NewHistogramVec(
		&metrics.HistogramOpts{
			Subsystem:      openstackSubsystem,
			Name:           openstackOperationKey,
			Help:           "Latency of openstack api call",
			StabilityLevel: metrics.ALPHA,
		},
		[]string{"request"},
	)

	openstackAPIRequestErrors = metrics.NewCounterVec(
		&metrics.CounterOpts{
			Subsystem:      openstackSubsystem,
			Name:           openstackOperationErrorKey,
			Help:           "Cumulative number of openstack Api call errors",
			StabilityLevel: metrics.ALPHA,
		},
		[]string{"request"},
	)
)

var registerOnce sync.Once

func registerMetrics() {
	registerOnce.Do(func() {
		legacyregistry.MustRegister(openstackOperationsLatency)
		legacyregistry.MustRegister(openstackAPIRequestErrors)
	})
}
