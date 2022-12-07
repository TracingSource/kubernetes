// Package internal provides trace internals.
package internal

// IDGenerator allows custom generators for TraceId and SpanId.
type IDGenerator interface {
	NewTraceID() [16]byte
	NewSpanID() [8]byte
}
