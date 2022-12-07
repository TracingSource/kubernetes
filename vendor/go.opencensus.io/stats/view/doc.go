// Package view contains support for collecting and exposing aggregates over stats.
//
// In order to collect measurements, views need to be defined and registered.
// A view allows recorded measurements to be filtered and aggregated.
//
// All recorded measurements can be grouped by a list of tags.
//
// OpenCensus provides several aggregation methods: Count, Distribution and Sum.
//
// Count only counts the number of measurement points recorded.
// Distribution provides statistical summary of the aggregated data by counting
// how many recorded measurements fall into each bucket.
// Sum adds up the measurement values.
// LastValue just keeps track of the most recently recorded measurement value.
// All aggregations are cumulative.
//
// Views can be registerd and unregistered at any time during program execution.
//
// Libraries can define views but it is recommended that in most cases registering
// views be left up to applications.
//
// Exporting
//
// Collected and aggregated data can be exported to a metric collection
// backend by registering its exporter.
//
// Multiple exporters can be registered to upload the data to various
// different back ends.
package view // import "go.opencensus.io/stats/view"

// TODO(acetechnologist): Add a link to the language independent OpenCensus
// spec when it is available.
