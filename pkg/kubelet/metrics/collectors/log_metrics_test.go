package collectors

import (
	"strings"
	"testing"

	"k8s.io/component-base/metrics/testutil"
	statsapi "k8s.io/kubernetes/pkg/kubelet/apis/stats/v1alpha1"
)

func TestNoMetricsCollected(t *testing.T) {
	// Refresh Desc to share with different registry
	descLogSize = descLogSize.GetRawDesc()

	collector := &logMetricsCollector{
		podStats: func() ([]statsapi.PodStats, error) {
			return []statsapi.PodStats{}, nil
		},
	}

	if err := testutil.CustomCollectAndCompare(collector, strings.NewReader(""), ""); err != nil {
		t.Fatal(err)
	}
}

func TestMetricsCollected(t *testing.T) {
	// Refresh Desc to share with different registry
	descLogSize = descLogSize.GetRawDesc()

	size := uint64(18)
	collector := &logMetricsCollector{
		podStats: func() ([]statsapi.PodStats, error) {
			return []statsapi.PodStats{
				{
					PodRef: statsapi.PodReference{
						Namespace: "some-namespace",
						Name:      "podName1",
						UID:       "UID_some_id",
					},
					Containers: []statsapi.ContainerStats{
						{
							Name: "containerName1",
							Logs: &statsapi.FsStats{
								UsedBytes: &size,
							},
						},
					},
				},
			}, nil
		},
	}

	err := testutil.CustomCollectAndCompare(collector, strings.NewReader(`
		# HELP kubelet_container_log_filesystem_used_bytes [ALPHA] Bytes used by the container's logs on the filesystem.
		# TYPE kubelet_container_log_filesystem_used_bytes gauge
		kubelet_container_log_filesystem_used_bytes{container="containerName1",namespace="some-namespace",pod="podName1",uid="UID_some_id"} 18
`), "kubelet_container_log_filesystem_used_bytes")
	if err != nil {
		t.Fatal(err)
	}
}
