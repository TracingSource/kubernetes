package main

import (
	"k8s.io/component-base/metrics"
)

const (
	counterMetricType   = "Counter"
	gaugeMetricType     = "Gauge"
	histogramMetricType = "Histogram"
)

type metric struct {
	Name              string    `yaml:"name"`
	Subsystem         string    `yaml:"subsystem,omitempty"`
	Namespace         string    `yaml:"namespace,omitempty"`
	Help              string    `yaml:"help,omitempty"`
	Type              string    `yaml:"type,omitempty"`
	DeprecatedVersion string    `yaml:"deprecatedVersion,omitempty"`
	StabilityLevel    string    `yaml:"stabilityLevel,omitempty"`
	Labels            []string  `yaml:"labels,omitempty"`
	Buckets           []float64 `yaml:"buckets,omitempty"`
}

func (m metric) buildFQName() string {
	return metrics.BuildFQName(m.Namespace, m.Subsystem, m.Name)
}

type byFQName []metric

func (ms byFQName) Len() int { return len(ms) }
func (ms byFQName) Less(i, j int) bool {
	return ms[i].buildFQName() < ms[j].buildFQName()
}
func (ms byFQName) Swap(i, j int) {
	ms[i], ms[j] = ms[j], ms[i]
}
