package external_metrics

import (
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// a list of values for a given metric for some set labels
type ExternalMetricValueList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`

	// value of the metric matching a given set of labels
	Items []ExternalMetricValue `json:"items"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// a metric value for external metric
// A single metric value is identified by metric name and a set of string labels.
// For one metric there can be multiple values with different sets of labels.
type ExternalMetricValue struct {
	metav1.TypeMeta `json:",inline"`

	// the name of the metric
	MetricName string `json:"metricName"`

	// a set of labels that identify a single time series for the metric
	MetricLabels map[string]string `json:"metricLabels"`

	// indicates the time at which the metrics were produced
	Timestamp metav1.Time `json:"timestamp"`

	// indicates the window ([Timestamp-Window, Timestamp]) from
	// which these metrics were calculated, when returning rate
	// metrics calculated from cumulative metrics (or zero for
	// non-calculated instantaneous metrics).
	WindowSeconds *int64 `json:"window,omitempty"`

	// the value of the metric
	Value resource.Quantity `json:"value"`
}
