package etcdhttp

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"go.etcd.io/etcd/etcdserver"
	"go.etcd.io/etcd/etcdserver/etcdserverpb"
	"go.etcd.io/etcd/raft"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	PathMetrics = "/metrics"
	PathHealth  = "/health"
)

// HandleMetricsHealth registers metrics and health handlers.
func HandleMetricsHealth(mux *http.ServeMux, srv etcdserver.ServerV2) {
	mux.Handle(PathMetrics, promhttp.Handler())
	mux.Handle(PathHealth, NewHealthHandler(func() Health { return checkHealth(srv) }))
}

// HandlePrometheus registers prometheus handler on '/metrics'.
func HandlePrometheus(mux *http.ServeMux) {
	mux.Handle(PathMetrics, promhttp.Handler())
}

// NewHealthHandler handles '/health' requests.
func NewHealthHandler(hfunc func() Health) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.Header().Set("Allow", http.MethodGet)
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}
		h := hfunc()
		d, _ := json.Marshal(h)
		if h.Health != "true" {
			http.Error(w, string(d), http.StatusServiceUnavailable)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write(d)
	}
}

var (
	healthSuccess = prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: "etcd",
		Subsystem: "server",
		Name:      "health_success",
		Help:      "The total number of successful health checks",
	})
	healthFailed = prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: "etcd",
		Subsystem: "server",
		Name:      "health_failures",
		Help:      "The total number of failed health checks",
	})
)

func init() {
	prometheus.MustRegister(healthSuccess)
	prometheus.MustRegister(healthFailed)
}

// Health defines etcd server health status.
// TODO: remove manual parsing in etcdctl cluster-health
type Health struct {
	Health string `json:"health"`
}

// TODO: server NOSPACE, etcdserver.ErrNoLeader in health API

func checkHealth(srv etcdserver.ServerV2) Health {
	h := Health{Health: "true"}

	as := srv.Alarms()
	if len(as) > 0 {
		h.Health = "false"
	}

	if h.Health == "true" {
		if uint64(srv.Leader()) == raft.None {
			h.Health = "false"
		}
	}

	if h.Health == "true" {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		_, err := srv.Do(ctx, etcdserverpb.Request{Method: "QGET"})
		cancel()
		if err != nil {
			h.Health = "false"
		}
	}

	if h.Health == "true" {
		healthSuccess.Inc()
	} else {
		healthFailed.Inc()
	}
	return h
}
