package metrics

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

// ValueType is an enumeration of metric types that represent a simple value.
type ValueType int

// Possible values for the ValueType enum.
const (
	_ ValueType = iota
	CounterValue
	GaugeValue
	UntypedValue
)

func (vt *ValueType) toPromValueType() prometheus.ValueType {
	return prometheus.ValueType(*vt)
}

// NewLazyConstMetric is a helper of MustNewConstMetric.
//
// Note: If the metrics described by the desc is hidden, the metrics will not be created.
func NewLazyConstMetric(desc *Desc, valueType ValueType, value float64, labelValues ...string) Metric {
	if desc.IsHidden() {
		return nil
	}
	return prometheus.MustNewConstMetric(desc.toPrometheusDesc(), valueType.toPromValueType(), value, labelValues...)
}

// NewLazyMetricWithTimestamp is a helper of NewMetricWithTimestamp.
//
// Warning: the Metric 'm' must be the one created by NewLazyConstMetric(),
//          otherwise, no stability guarantees would be offered.
func NewLazyMetricWithTimestamp(t time.Time, m Metric) Metric {
	if m == nil {
		return nil
	}

	return prometheus.NewMetricWithTimestamp(t, m)
}
