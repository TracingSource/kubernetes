package metrics

var interestingAPIServerMetrics = []string{
	"apiserver_request_total",
	"apiserver_request_latency_seconds",
	"apiserver_init_events_total",
}

var interestingControllerManagerMetrics = []string{
	"garbage_collector_attempt_to_delete_queue_latency",
	"garbage_collector_attempt_to_delete_work_duration",
	"garbage_collector_attempt_to_orphan_queue_latency",
	"garbage_collector_attempt_to_orphan_work_duration",
	"garbage_collector_dirty_processing_latency_microseconds",
	"garbage_collector_event_processing_latency_microseconds",
	"garbage_collector_graph_changes_queue_latency",
	"garbage_collector_graph_changes_work_duration",
	"garbage_collector_orphan_processing_latency_microseconds",

	"namespace_queue_latency",
	"namespace_queue_latency_sum",
	"namespace_queue_latency_count",
	"namespace_retries",
	"namespace_work_duration",
	"namespace_work_duration_sum",
	"namespace_work_duration_count",
}

var interestingKubeletMetrics = []string{
	"kubelet_docker_operations_errors_total",
	"kubelet_docker_operations_duration_seconds",
	"kubelet_pod_start_duration_seconds",
	"kubelet_pod_worker_duration_seconds",
	"kubelet_pod_worker_start_duration_seconds",
}

var interestingClusterAutoscalerMetrics = []string{
	"function_duration_seconds",
	"errors_total",
	"evicted_pods_total",
}
