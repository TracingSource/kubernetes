package metrics

import (
	"os"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/procfs"

	"k8s.io/klog"
)

var processStartTime = prometheus.NewGaugeVec(
	prometheus.GaugeOpts{
		Name: "process_start_time_seconds",
		Help: "Start time of the process since unix epoch in seconds.",
	},
	[]string{},
)

// Registerer is an interface expected by RegisterProcessStartTime in order to register the metric
type Registerer interface {
	Register(prometheus.Collector) error
}

// RegisterProcessStartTime registers the process_start_time_seconds to
// a prometheus registry. This metric needs to be included to ensure counter
// data fidelity.
func RegisterProcessStartTime(registrationFunc func(prometheus.Collector) error) error {
	start, err := getProcessStart()
	if err != nil {
		klog.Errorf("Could not get process start time, %v", err)
		start = float64(time.Now().Unix())
	}
	processStartTime.WithLabelValues().Set(start)
	return registrationFunc(processStartTime)
}

func getProcessStart() (float64, error) {
	pid := os.Getpid()
	p, err := procfs.NewProc(pid)
	if err != nil {
		return 0, err
	}

	if stat, err := p.NewStat(); err == nil {
		return stat.StartTime()
	}
	return 0, err
}
