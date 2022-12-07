package metrics

import (
	"sync"
	"time"

	"k8s.io/component-base/metrics"
	"k8s.io/component-base/metrics/legacyregistry"
)

const (
	kubeletSubsystem = "kubelet"
)

var (
	// HTTPRequests tracks the number of the http requests received since the server started.
	HTTPRequests = metrics.NewCounterVec(
		&metrics.CounterOpts{
			Subsystem:      kubeletSubsystem,
			Name:           "http_requests_total",
			Help:           "Number of the http requests received since the server started",
			StabilityLevel: metrics.ALPHA,
		},
		// server_type aims to differentiate the readonly server and the readwrite server.
		// long_running marks whether the request is long-running or not.
		// Currently, long-running requests include exec/attach/portforward/debug.
		[]string{"method", "path", "server_type", "long_running"},
	)
	// HTTPRequestsDuration tracks the duration in seconds to serve http requests.
	HTTPRequestsDuration = metrics.NewHistogramVec(
		&metrics.HistogramOpts{
			Subsystem: kubeletSubsystem,
			Name:      "http_requests_duration_seconds",
			Help:      "Duration in seconds to serve http requests",
			// Use DefBuckets for now, will customize the buckets if necessary.
			Buckets:        metrics.DefBuckets,
			StabilityLevel: metrics.ALPHA,
		},
		[]string{"method", "path", "server_type", "long_running"},
	)
	// HTTPInflightRequests tracks the number of the inflight http requests.
	HTTPInflightRequests = metrics.NewGaugeVec(
		&metrics.GaugeOpts{
			Subsystem:      kubeletSubsystem,
			Name:           "http_inflight_requests",
			Help:           "Number of the inflight http requests",
			StabilityLevel: metrics.ALPHA,
		},
		[]string{"method", "path", "server_type", "long_running"},
	)
)

var registerMetrics sync.Once

// Register all metrics.
func Register() {
	registerMetrics.Do(func() {
		legacyregistry.MustRegister(HTTPRequests)
		legacyregistry.MustRegister(HTTPRequestsDuration)
		legacyregistry.MustRegister(HTTPInflightRequests)
	})
}

// SinceInSeconds gets the time since the specified start in seconds.
func SinceInSeconds(start time.Time) float64 {
	return time.Since(start).Seconds()
}
