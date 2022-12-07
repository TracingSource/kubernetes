package stats

import (
	"fmt"

	"k8s.io/kubernetes/pkg/volume"
)

type fakeLogMetrics struct {
	fakeStats map[string]*volume.Metrics
}

func NewFakeLogMetricsService(stats map[string]*volume.Metrics) LogMetricsService {
	return &fakeLogMetrics{fakeStats: stats}
}

func (l *fakeLogMetrics) createLogMetricsProvider(path string) volume.MetricsProvider {
	return NewFakeMetricsDu(path, l.fakeStats[path])
}

type fakeMetricsDu struct {
	fakeStats *volume.Metrics
}

func NewFakeMetricsDu(path string, stats *volume.Metrics) volume.MetricsProvider {
	return &fakeMetricsDu{fakeStats: stats}
}

func (f *fakeMetricsDu) GetMetrics() (*volume.Metrics, error) {
	if f.fakeStats == nil {
		return nil, fmt.Errorf("no stats provided")
	}
	return f.fakeStats, nil
}
