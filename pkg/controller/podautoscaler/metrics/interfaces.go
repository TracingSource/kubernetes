package metrics

import (
	"time"

	autoscaling "k8s.io/api/autoscaling/v2beta2"
	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/labels"
)

// PodMetric contains pod metric value (the metric values are expected to be the metric as a milli-value)
type PodMetric struct {
	Timestamp time.Time
	Window    time.Duration
	Value     int64
}

// PodMetricsInfo contains pod metrics as a map from pod names to PodMetricsInfo
type PodMetricsInfo map[string]PodMetric

// MetricsClient knows how to query a remote interface to retrieve container-level
// resource metrics as well as pod-level arbitrary metrics
type MetricsClient interface {
	// GetResourceMetric gets the given resource metric (and an associated oldest timestamp)
	// for all pods matching the specified selector in the given namespace
	GetResourceMetric(resource v1.ResourceName, namespace string, selector labels.Selector) (PodMetricsInfo, time.Time, error)

	// GetRawMetric gets the given metric (and an associated oldest timestamp)
	// for all pods matching the specified selector in the given namespace
	GetRawMetric(metricName string, namespace string, selector labels.Selector, metricSelector labels.Selector) (PodMetricsInfo, time.Time, error)

	// GetObjectMetric gets the given metric (and an associated timestamp) for the given
	// object in the given namespace
	GetObjectMetric(metricName string, namespace string, objectRef *autoscaling.CrossVersionObjectReference, metricSelector labels.Selector) (int64, time.Time, error)

	// GetExternalMetric gets all the values of a given external metric
	// that match the specified selector.
	GetExternalMetric(metricName string, namespace string, selector labels.Selector) ([]int64, time.Time, error)
}
