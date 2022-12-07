package external_metrics

import (
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/metrics/pkg/apis/external_metrics/v1beta1"
)

// ExternalMetricsClient is a client for fetching external metrics.
type ExternalMetricsClient interface {
	NamespacedMetricsGetter
}

// NamespacedMetricsGetter provides access to an interface for fetching
// metrics in a particular namespace.
type NamespacedMetricsGetter interface {
	NamespacedMetrics(namespace string) MetricsInterface
}

// MetricsInterface provides access to external metrics.
type MetricsInterface interface {
	// List fetches the metric for the given namespace that maches the given
	// metricSelector.
	List(metricName string, metricSelector labels.Selector) (*v1beta1.ExternalMetricValueList, error)
}
