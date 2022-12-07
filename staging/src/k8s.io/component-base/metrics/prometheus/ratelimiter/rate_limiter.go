package ratelimiter

import (
	"fmt"
	"sync"

	"k8s.io/client-go/util/flowcontrol"
	"k8s.io/component-base/metrics"
	"k8s.io/component-base/metrics/legacyregistry"
)

var (
	metricsLock        sync.Mutex
	rateLimiterMetrics = make(map[string]*rateLimiterMetric)
)

type rateLimiterMetric struct {
	metric metrics.GaugeMetric
	stopCh chan struct{}
}

func registerRateLimiterMetric(ownerName string) error {
	metricsLock.Lock()
	defer metricsLock.Unlock()

	if _, ok := rateLimiterMetrics[ownerName]; ok {
		// only register once in Prometheus. We happen to see an ownerName reused in parallel integration tests.
		return nil
	}
	metric := metrics.NewGauge(&metrics.GaugeOpts{
		Name:           "rate_limiter_use",
		Subsystem:      ownerName,
		Help:           fmt.Sprintf("A metric measuring the saturation of the rate limiter for %v", ownerName),
		StabilityLevel: metrics.ALPHA,
	})
	if err := legacyregistry.Register(metric); err != nil {
		return fmt.Errorf("error registering rate limiter usage metric: %v", err)
	}
	stopCh := make(chan struct{})
	rateLimiterMetrics[ownerName] = &rateLimiterMetric{
		metric: metric,
		stopCh: stopCh,
	}
	return nil
}

// RegisterMetricAndTrackRateLimiterUsage registers a metric ownerName_rate_limiter_use in prometheus to track
// how much used rateLimiter is and starts a goroutine that updates this metric every updatePeriod
func RegisterMetricAndTrackRateLimiterUsage(ownerName string, rateLimiter flowcontrol.RateLimiter) error {
	if err := registerRateLimiterMetric(ownerName); err != nil {
		return err
	}
	// TODO: determine how to track rate limiter saturation
	// See discussion at https://go-review.googlesource.com/c/time/+/29958#message-4caffc11669cadd90e2da4c05122cfec50ea6a22
	// go wait.Until(func() {
	//   metricsLock.Lock()
	//   defer metricsLock.Unlock()
	//   rateLimiterMetrics[ownerName].metric.Set()
	// }, updatePeriod, rateLimiterMetrics[ownerName].stopCh)
	return nil
}
