package metrics

import "k8s.io/component-base/metrics/testutil"

// ControllerManagerMetrics is metrics for controller manager
type ControllerManagerMetrics testutil.Metrics

// Equal returns true if all metrics are the same as the arguments.
func (m *ControllerManagerMetrics) Equal(o ControllerManagerMetrics) bool {
	return (*testutil.Metrics)(m).Equal(testutil.Metrics(o))
}

func newControllerManagerMetrics() ControllerManagerMetrics {
	result := testutil.NewMetrics()
	return ControllerManagerMetrics(result)
}

func parseControllerManagerMetrics(data string) (ControllerManagerMetrics, error) {
	result := newControllerManagerMetrics()
	if err := testutil.ParseMetrics(data, (*testutil.Metrics)(&result)); err != nil {
		return ControllerManagerMetrics{}, err
	}
	return result, nil
}
