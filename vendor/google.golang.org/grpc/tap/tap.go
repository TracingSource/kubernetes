/*
 *
 * Copyright 2016 gRPC authors.
 *

 *
 */

// Package tap defines the function handles which are executed on the transport
// layer of gRPC-Go and related information. Everything here is EXPERIMENTAL.
package tap

import (
	"context"
)

// Info defines the relevant information needed by the handles.
type Info struct {
	// FullMethodName is the string of grpc method (in the format of
	// /package.service/method).
	FullMethodName string
	// TODO: More to be added.
}

// ServerInHandle defines the function which runs before a new stream is created
// on the server side. If it returns a non-nil error, the stream will not be
// created and a RST_STREAM will be sent back to the client with REFUSED_STREAM.
// The client will receive an RPC error "code = Unavailable, desc = stream
// terminated by RST_STREAM with error code: REFUSED_STREAM".
//
// It's intended to be used in situations where you don't want to waste the
// resources to accept the new stream (e.g. rate-limiting). And the content of
// the error will be ignored and won't be sent back to the client. For other
// general usages, please use interceptors.
//
// Note that it is executed in the per-connection I/O goroutine(s) instead of
// per-RPC goroutine. Therefore, users should NOT have any
// blocking/time-consuming work in this handle. Otherwise all the RPCs would
// slow down. Also, for the same reason, this handle won't be called
// concurrently by gRPC.
type ServerInHandle func(ctx context.Context, info *Info) (context.Context, error)
