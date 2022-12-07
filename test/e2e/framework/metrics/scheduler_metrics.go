package metrics

import "k8s.io/component-base/metrics/testutil"

// SchedulerMetrics is metrics for scheduler
type SchedulerMetrics testutil.Metrics

// Equal returns true if all metrics are the same as the arguments.
func (m *SchedulerMetrics) Equal(o SchedulerMetrics) bool {
	return (*testutil.Metrics)(m).Equal(testutil.Metrics(o))
}

func newSchedulerMetrics() SchedulerMetrics {
	result := testutil.NewMetrics()
	return SchedulerMetrics(result)
}

func parseSchedulerMetrics(data string) (SchedulerMetrics, error) {
	result := newSchedulerMetrics()
	if err := testutil.ParseMetrics(data, (*testutil.Metrics)(&result)); err != nil {
		return SchedulerMetrics{}, err
	}
	return result, nil
}
