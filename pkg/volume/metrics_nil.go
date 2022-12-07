package volume

var _ MetricsProvider = &MetricsNil{}

// MetricsNil represents a MetricsProvider that does not support returning
// Metrics.  It serves as a placeholder for Volumes that do not yet support
// metrics.
type MetricsNil struct{}

// GetMetrics returns an empty Metrics and an error.
// See MetricsProvider.GetMetrics
func (*MetricsNil) GetMetrics() (*Metrics, error) {
	return &Metrics{}, NewNotSupportedError()
}
