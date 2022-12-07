package metrics

// SchedulingMetrics is a struct for managing scheduling metrics.
type SchedulingMetrics struct {
	PredicateEvaluationLatency  LatencyMetric `json:"predicateEvaluationLatency"`
	PriorityEvaluationLatency   LatencyMetric `json:"priorityEvaluationLatency"`
	PreemptionEvaluationLatency LatencyMetric `json:"preemptionEvaluationLatency"`
	BindingLatency              LatencyMetric `json:"bindingLatency"`
	ThroughputAverage           float64       `json:"throughputAverage"`
	ThroughputPerc50            float64       `json:"throughputPerc50"`
	ThroughputPerc90            float64       `json:"throughputPerc90"`
	ThroughputPerc99            float64       `json:"throughputPerc99"`
}

// SummaryKind returns the summary of scheduling metrics.
func (l *SchedulingMetrics) SummaryKind() string {
	return "SchedulingMetrics"
}

// PrintHumanReadable returns scheduling metrics with JSON format.
func (l *SchedulingMetrics) PrintHumanReadable() string {
	return PrettyPrintJSON(l)
}

// PrintJSON returns scheduling metrics with JSON format.
func (l *SchedulingMetrics) PrintJSON() string {
	return PrettyPrintJSON(l)
}
