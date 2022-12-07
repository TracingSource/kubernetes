/*
Package tag contains OpenCensus tags.

Tags are key-value pairs. Tags provide additional cardinality to
the OpenCensus instrumentation data.

Tags can be propagated on the wire and in the same
process via context.Context. Encode and Decode should be
used to represent tags into their binary propagation form.
*/
package tag // import "go.opencensus.io/tag"
