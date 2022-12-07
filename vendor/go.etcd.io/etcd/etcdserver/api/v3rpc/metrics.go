package v3rpc

import "github.com/prometheus/client_golang/prometheus"

var (
	sentBytes = prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: "etcd",
		Subsystem: "network",
		Name:      "client_grpc_sent_bytes_total",
		Help:      "The total number of bytes sent to grpc clients.",
	})

	receivedBytes = prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: "etcd",
		Subsystem: "network",
		Name:      "client_grpc_received_bytes_total",
		Help:      "The total number of bytes received from grpc clients.",
	})

	streamFailures = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: "etcd",
		Subsystem: "network",
		Name:      "server_stream_failures_total",
		Help:      "The total number of stream failures from the local server.",
	},
		[]string{"Type", "API"},
	)
)

func init() {
	prometheus.MustRegister(sentBytes)
	prometheus.MustRegister(receivedBytes)
	prometheus.MustRegister(streamFailures)
}
