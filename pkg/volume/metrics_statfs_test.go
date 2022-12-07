package volume_test

import (
	"os"
	"testing"

	utiltesting "k8s.io/client-go/util/testing"
	. "k8s.io/kubernetes/pkg/volume"
	volumetest "k8s.io/kubernetes/pkg/volume/testing"
)

func TestGetMetricsStatFS(t *testing.T) {
	metrics := NewMetricsStatFS("")
	actual, err := metrics.GetMetrics()
	expected := &Metrics{}
	if !volumetest.MetricsEqualIgnoreTimestamp(actual, expected) {
		t.Errorf("Expected empty Metrics from uninitialized MetricsStatFS, actual %v", *actual)
	}
	if err == nil {
		t.Errorf("Expected error when calling GetMetrics on uninitialized MetricsStatFS, actual nil")
	}

	metrics = NewMetricsStatFS("/not/a/real/directory")
	actual, err = metrics.GetMetrics()
	if !volumetest.MetricsEqualIgnoreTimestamp(actual, expected) {
		t.Errorf("Expected empty Metrics from incorrectly initialized MetricsStatFS, actual %v", *actual)
	}
	if err == nil {
		t.Errorf("Expected error when calling GetMetrics on incorrectly initialized MetricsStatFS, actual nil")
	}

	tmpDir, err := utiltesting.MkTmpdir("metric_statfs_test")
	if err != nil {
		t.Fatalf("Can't make a tmp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	metrics = NewMetricsStatFS(tmpDir)
	actual, err = metrics.GetMetrics()
	if err != nil {
		t.Errorf("Unexpected error when calling GetMetrics %v", err)
	}

	if a := actual.Capacity.Value(); a <= 0 {
		t.Errorf("Expected Capacity %d to be greater than 0.", a)
	}
	if a := actual.Available.Value(); a <= 0 {
		t.Errorf("Expected Available %d to be greater than 0.", a)
	}

}
