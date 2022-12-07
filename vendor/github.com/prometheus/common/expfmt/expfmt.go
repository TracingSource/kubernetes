// Package expfmt contains tools for reading and writing Prometheus metrics.
package expfmt

// Format specifies the HTTP content type of the different wire protocols.
type Format string

// Constants to assemble the Content-Type values for the different wire protocols.
const (
	TextVersion   = "0.0.4"
	ProtoType     = `application/vnd.google.protobuf`
	ProtoProtocol = `io.prometheus.client.MetricFamily`
	ProtoFmt      = ProtoType + "; proto=" + ProtoProtocol + ";"

	// The Content-Type values for the different wire protocols.
	FmtUnknown      Format = `<unknown>`
	FmtText         Format = `text/plain; version=` + TextVersion + `; charset=utf-8`
	FmtProtoDelim   Format = ProtoFmt + ` encoding=delimited`
	FmtProtoText    Format = ProtoFmt + ` encoding=text`
	FmtProtoCompact Format = ProtoFmt + ` encoding=compact-text`
)

const (
	hdrContentType = "Content-Type"
	hdrAccept      = "Accept"
)
