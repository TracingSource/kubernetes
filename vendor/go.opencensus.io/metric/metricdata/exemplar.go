package metricdata

import (
	"time"
)

// Exemplars keys.
const (
	AttachmentKeySpanContext = "SpanContext"
)

// Exemplar is an example data point associated with each bucket of a
// distribution type aggregation.
//
// Their purpose is to provide an example of the kind of thing
// (request, RPC, trace span, etc.) that resulted in that measurement.
type Exemplar struct {
	Value       float64     // the value that was recorded
	Timestamp   time.Time   // the time the value was recorded
	Attachments Attachments // attachments (if any)
}

// Attachments is a map of extra values associated with a recorded data point.
type Attachments map[string]interface{}
