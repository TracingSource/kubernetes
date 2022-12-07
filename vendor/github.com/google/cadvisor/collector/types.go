package collector

import (
	"time"

	"github.com/google/cadvisor/info/v1"
)

// TODO(vmarmol): Export to a custom metrics type when that is available.

// Metric collector.
type Collector interface {
	// Collect metrics from this collector.
	// Returns the next time this collector should be collected from.
	// Next collection time is always returned, even when an error occurs.
	// A collection time of zero means no more collection.
	Collect(map[string][]v1.MetricVal) (time.Time, map[string][]v1.MetricVal, error)

	// Return spec for all metrics associated with this collector
	GetSpec() []v1.MetricSpec

	// Name of this collector.
	Name() string
}

// Manages and runs collectors.
type CollectorManager interface {
	// Register a collector.
	RegisterCollector(collector Collector) error

	// Collect from collectors that are ready and return the next time
	// at which a collector will be ready to collect from.
	// Next collection time is always returned, even when an error occurs.
	// A collection time of zero means no more collection.
	Collect() (time.Time, map[string][]v1.MetricVal, error)

	// Get metric spec from all registered collectors.
	GetSpec() ([]v1.MetricSpec, error)
}
