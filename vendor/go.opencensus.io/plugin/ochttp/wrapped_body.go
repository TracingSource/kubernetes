// Copyright 2019, OpenCensus Authors
//


package ochttp

import (
	"io"
)

// wrappedBody returns a wrapped version of the original
// Body and only implements the same combination of additional
// interfaces as the original.
func wrappedBody(wrapper io.ReadCloser, body io.ReadCloser) io.ReadCloser {
	var (
		wr, i0 = body.(io.Writer)
	)
	switch {
	case !i0:
		return struct {
			io.ReadCloser
		}{wrapper}

	case i0:
		return struct {
			io.ReadCloser
			io.Writer
		}{wrapper, wr}
	default:
		return struct {
			io.ReadCloser
		}{wrapper}
	}
}
