package volume

import (
	"testing"
)

func TestMetricsNilGetCapacity(t *testing.T) {
	metrics := &MetricsNil{}
	actual, err := metrics.GetMetrics()
	expected := &Metrics{}
	if *actual != *expected {
		t.Errorf("Expected empty Metrics, actual %v", *actual)
	}
	if err == nil {
		t.Errorf("Expected error when calling GetMetrics, actual nil")
	}
}
