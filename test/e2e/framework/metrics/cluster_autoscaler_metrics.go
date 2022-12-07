package metrics

import "k8s.io/component-base/metrics/testutil"

// ClusterAutoscalerMetrics is metrics for cluster autoscaler
type ClusterAutoscalerMetrics testutil.Metrics

// Equal returns true if all metrics are the same as the arguments.
func (m *ClusterAutoscalerMetrics) Equal(o ClusterAutoscalerMetrics) bool {
	return (*testutil.Metrics)(m).Equal(testutil.Metrics(o))
}

func newClusterAutoscalerMetrics() ClusterAutoscalerMetrics {
	result := testutil.NewMetrics()
	return ClusterAutoscalerMetrics(result)
}

func parseClusterAutoscalerMetrics(data string) (ClusterAutoscalerMetrics, error) {
	result := newClusterAutoscalerMetrics()
	if err := testutil.ParseMetrics(data, (*testutil.Metrics)(&result)); err != nil {
		return ClusterAutoscalerMetrics{}, err
	}
	return result, nil
}
