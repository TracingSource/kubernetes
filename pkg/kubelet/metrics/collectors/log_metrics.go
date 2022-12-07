package collectors

import (
	"k8s.io/component-base/metrics"
	"k8s.io/klog"
	statsapi "k8s.io/kubernetes/pkg/kubelet/apis/stats/v1alpha1"
)

var (
	descLogSize = metrics.NewDesc(
		"kubelet_container_log_filesystem_used_bytes",
		"Bytes used by the container's logs on the filesystem.",
		[]string{
			"uid",
			"namespace",
			"pod",
			"container",
		}, nil,
		metrics.ALPHA,
		"",
	)
)

type logMetricsCollector struct {
	metrics.BaseStableCollector

	podStats func() ([]statsapi.PodStats, error)
}

// Check if logMetricsCollector implements necessary interface
var _ metrics.StableCollector = &logMetricsCollector{}

// NewLogMetricsCollector implements the metrics.StableCollector interface and
// exposes metrics about container's log volume size.
func NewLogMetricsCollector(podStats func() ([]statsapi.PodStats, error)) metrics.StableCollector {
	return &logMetricsCollector{
		podStats: podStats,
	}
}

// DescribeWithStability implements the metrics.StableCollector interface.
func (c *logMetricsCollector) DescribeWithStability(ch chan<- *metrics.Desc) {
	ch <- descLogSize
}

// CollectWithStability implements the metrics.StableCollector interface.
func (c *logMetricsCollector) CollectWithStability(ch chan<- metrics.Metric) {
	podStats, err := c.podStats()
	if err != nil {
		klog.Errorf("failed to get pod stats: %v", err)
		return
	}

	for _, ps := range podStats {
		for _, c := range ps.Containers {
			if c.Logs != nil && c.Logs.UsedBytes != nil {
				ch <- metrics.NewLazyConstMetric(
					descLogSize,
					metrics.GaugeValue,
					float64(*c.Logs.UsedBytes),
					ps.PodRef.UID,
					ps.PodRef.Namespace,
					ps.PodRef.Name,
					c.Name,
				)
			}
		}
	}
}
