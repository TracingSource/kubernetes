package testutil

import (
	"io"

	"github.com/prometheus/client_golang/prometheus/testutil"

	"k8s.io/component-base/metrics"
)

// CollectAndCompare registers the provided Collector with a newly created
// pedantic Registry. It then does the same as GatherAndCompare, gathering the
// metrics from the pedantic Registry.
func CollectAndCompare(c metrics.Collector, expected io.Reader, metricNames ...string) error {
	return testutil.CollectAndCompare(c, expected, metricNames...)
}

// GatherAndCompare gathers all metrics from the provided Gatherer and compares
// it to an expected output read from the provided Reader in the Prometheus text
// exposition format. If any metricNames are provided, only metrics with those
// names are compared.
func GatherAndCompare(g metrics.Gatherer, expected io.Reader, metricNames ...string) error {
	return testutil.GatherAndCompare(g, expected, metricNames...)
}

// CustomCollectAndCompare registers the provided StableCollector with a newly created
// registry. It then does the same as GatherAndCompare, gathering the
// metrics from the pedantic Registry.
func CustomCollectAndCompare(c metrics.StableCollector, expected io.Reader, metricNames ...string) error {
	registry := metrics.NewKubeRegistry()
	registry.CustomMustRegister(c)

	return GatherAndCompare(registry, expected, metricNames...)
}
