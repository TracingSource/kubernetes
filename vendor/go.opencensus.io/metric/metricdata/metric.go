package metricdata

import (
	"time"

	"go.opencensus.io/resource"
)

// Descriptor holds metadata about a metric.
type Descriptor struct {
	Name        string     // full name of the metric
	Description string     // human-readable description
	Unit        Unit       // units for the measure
	Type        Type       // type of measure
	LabelKeys   []LabelKey // label keys
}

// Metric represents a quantity measured against a resource with different
// label value combinations.
type Metric struct {
	Descriptor Descriptor         // metric descriptor
	Resource   *resource.Resource // resource against which this was measured
	TimeSeries []*TimeSeries      // one time series for each combination of label values
}

// TimeSeries is a sequence of points associated with a combination of label
// values.
type TimeSeries struct {
	LabelValues []LabelValue // label values, same order as keys in the metric descriptor
	Points      []Point      // points sequence
	StartTime   time.Time    // time we started recording this time series
}
