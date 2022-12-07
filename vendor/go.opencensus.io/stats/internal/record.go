package internal

import (
	"go.opencensus.io/tag"
)

// DefaultRecorder will be called for each Record call.
var DefaultRecorder func(tags *tag.Map, measurement interface{}, attachments map[string]interface{})

// SubscriptionReporter reports when a view subscribed with a measure.
var SubscriptionReporter func(measure string)
