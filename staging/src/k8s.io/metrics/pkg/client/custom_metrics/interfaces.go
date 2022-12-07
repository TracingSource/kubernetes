package custom_metrics

import (
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/metrics/pkg/apis/custom_metrics/v1beta2"
)

// CustomMetricsClient is a client for fetching metrics
// describing both root-scoped and namespaced resources.
type CustomMetricsClient interface {
	RootScopedMetricsGetter
	NamespacedMetricsGetter
}

// RootScopedMetricsGetter provides access to an interface for fetching
// metrics describing root-scoped objects.  Note that metrics describing
// a namespace are simply considered a special case of root-scoped metrics.
type RootScopedMetricsGetter interface {
	RootScopedMetrics() MetricsInterface
}

// NamespacedMetricsGetter provides access to an interface for fetching
// metrics describing resources in a particular namespace.
type NamespacedMetricsGetter interface {
	NamespacedMetrics(namespace string) MetricsInterface
}

// MetricsInterface provides access to metrics describing Kubernetes objects.
type MetricsInterface interface {
	// GetForObject fetchs the given metric describing the given object.
	GetForObject(groupKind schema.GroupKind, name string, metricName string, metricSelector labels.Selector) (*v1beta2.MetricValue, error)

	// GetForObjects fetches the given metric describing all objects of the given
	// type matching the given label selector (or simply all objects of the given type
	// if the selector is nil).
	GetForObjects(groupKind schema.GroupKind, selector labels.Selector, metricName string, metricSelector labels.Selector) (*v1beta2.MetricValueList, error)
}
