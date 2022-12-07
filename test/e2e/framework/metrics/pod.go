package metrics

import (
	"time"
)

// LatencyMetric is a struct for dashboard metrics.
type LatencyMetric struct {
	Perc50  time.Duration `json:"Perc50"`
	Perc90  time.Duration `json:"Perc90"`
	Perc99  time.Duration `json:"Perc99"`
	Perc100 time.Duration `json:"Perc100"`
}
