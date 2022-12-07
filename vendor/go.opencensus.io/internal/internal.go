package internal // import "go.opencensus.io/internal"

import (
	"fmt"
	"time"

	opencensus "go.opencensus.io"
)

// UserAgent is the user agent to be added to the outgoing
// requests from the exporters.
var UserAgent = fmt.Sprintf("opencensus-go/%s", opencensus.Version())

// MonotonicEndTime returns the end time at present
// but offset from start, monotonically.
//
// The monotonic clock is used in subtractions hence
// the duration since start added back to start gives
// end as a monotonic time.
// See https://golang.org/pkg/time/#hdr-Monotonic_Clocks
func MonotonicEndTime(start time.Time) time.Time {
	return start.Add(time.Now().Sub(start))
}
