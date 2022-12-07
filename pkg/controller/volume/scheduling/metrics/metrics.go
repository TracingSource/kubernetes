package metrics

import (
	"k8s.io/component-base/metrics"
	"k8s.io/component-base/metrics/legacyregistry"
)

// VolumeSchedulerSubsystem - subsystem name used by scheduler
const VolumeSchedulerSubsystem = "scheduler_volume"

var (
	// VolumeBindingRequestSchedulerBinderCache tracks the number of volume binder cache operations.
	VolumeBindingRequestSchedulerBinderCache = metrics.NewCounterVec(
		&metrics.CounterOpts{
			Subsystem:      VolumeSchedulerSubsystem,
			Name:           "binder_cache_requests_total",
			Help:           "Total number for request volume binding cache",
			StabilityLevel: metrics.ALPHA,
		},
		[]string{"operation"},
	)
	// VolumeSchedulingStageLatency tracks the latency of volume scheduling operations.
	VolumeSchedulingStageLatency = metrics.NewHistogramVec(
		&metrics.HistogramOpts{
			Subsystem:      VolumeSchedulerSubsystem,
			Name:           "scheduling_duration_seconds",
			Help:           "Volume scheduling stage latency",
			Buckets:        metrics.ExponentialBuckets(1000, 2, 15),
			StabilityLevel: metrics.ALPHA,
		},
		[]string{"operation"},
	)
	// VolumeSchedulingStageFailed tracks the number of failed volume scheduling operations.
	VolumeSchedulingStageFailed = metrics.NewCounterVec(
		&metrics.CounterOpts{
			Subsystem:      VolumeSchedulerSubsystem,
			Name:           "scheduling_stage_error_total",
			Help:           "Volume scheduling stage error count",
			StabilityLevel: metrics.ALPHA,
		},
		[]string{"operation"},
	)
)

// RegisterVolumeSchedulingMetrics is used for scheduler, because the volume binding cache is a library
// used by scheduler process.
func RegisterVolumeSchedulingMetrics() {
	legacyregistry.MustRegister(VolumeBindingRequestSchedulerBinderCache)
	legacyregistry.MustRegister(VolumeSchedulingStageLatency)
	legacyregistry.MustRegister(VolumeSchedulingStageFailed)
}
