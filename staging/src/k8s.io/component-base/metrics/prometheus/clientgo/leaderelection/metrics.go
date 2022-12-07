package leaderelection

import (
	"k8s.io/client-go/tools/leaderelection"
	k8smetrics "k8s.io/component-base/metrics"
	"k8s.io/component-base/metrics/legacyregistry"
)

var (
	leaderGauge = k8smetrics.NewGaugeVec(&k8smetrics.GaugeOpts{
		Name: "leader_election_master_status",
		Help: "Gauge of if the reporting system is master of the relevant lease, 0 indicates backup, 1 indicates master. 'name' is the string used to identify the lease. Please make sure to group by name.",
	}, []string{"name"})
)

func init() {
	legacyregistry.MustRegister(leaderGauge)
	leaderelection.SetProvider(prometheusMetricsProvider{})
}

type prometheusMetricsProvider struct{}

func (prometheusMetricsProvider) NewLeaderMetric() leaderelection.SwitchMetric {
	return &switchAdapter{gauge: leaderGauge}
}

type switchAdapter struct {
	gauge *k8smetrics.GaugeVec
}

func (s *switchAdapter) On(name string) {
	s.gauge.WithLabelValues(name).Set(1.0)
}

func (s *switchAdapter) Off(name string) {
	s.gauge.WithLabelValues(name).Set(0.0)
}
