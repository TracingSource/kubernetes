// Copyright 2019, OpenCensus Authors
//


package metricproducer

import (
	"go.opencensus.io/metric/metricdata"
)

// Producer is a source of metrics.
type Producer interface {
	// Read should return the current values of all metrics supported by this
	// metric provider.
	// The returned metrics should be unique for each combination of name and
	// resource.
	Read() []*metricdata.Metric
}
