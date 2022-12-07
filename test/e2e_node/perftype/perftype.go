package perftype

// ResourceSeries defines the time series of the resource usage.
type ResourceSeries struct {
	Timestamp            []int64           `json:"ts"`
	CPUUsageInMilliCores []int64           `json:"cpu"`
	MemoryRSSInMegaBytes []int64           `json:"memory"`
	Units                map[string]string `json:"unit"`
}

// NodeTimeSeries defines the time series of the operations and the resource
// usage.
type NodeTimeSeries struct {
	OperationData map[string][]int64         `json:"op_series,omitempty"`
	ResourceData  map[string]*ResourceSeries `json:"resource_series,omitempty"`
	Labels        map[string]string          `json:"labels"`
	Version       string                     `json:"version"`
}
