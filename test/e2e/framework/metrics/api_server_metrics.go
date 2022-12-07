package metrics

import "k8s.io/component-base/metrics/testutil"

// APIServerMetrics is metrics for API server
type APIServerMetrics testutil.Metrics

// Equal returns true if all metrics are the same as the arguments.
func (m *APIServerMetrics) Equal(o APIServerMetrics) bool {
	return (*testutil.Metrics)(m).Equal(testutil.Metrics(o))
}

func newAPIServerMetrics() APIServerMetrics {
	result := testutil.NewMetrics()
	return APIServerMetrics(result)
}

func parseAPIServerMetrics(data string) (APIServerMetrics, error) {
	result := newAPIServerMetrics()
	if err := testutil.ParseMetrics(data, (*testutil.Metrics)(&result)); err != nil {
		return APIServerMetrics{}, err
	}
	return result, nil
}

func (g *Grabber) getMetricsFromAPIServer() (string, error) {
	rawOutput, err := g.client.CoreV1().RESTClient().Get().RequestURI("/metrics").Do().Raw()
	if err != nil {
		return "", err
	}
	return string(rawOutput), nil
}
