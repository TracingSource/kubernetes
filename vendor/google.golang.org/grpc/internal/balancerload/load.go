/*
 * Copyright 2019 gRPC authors.
 *

 */

// Package balancerload defines APIs to parse server loads in trailers. The
// parsed loads are sent to balancers in DoneInfo.
package balancerload

import (
	"google.golang.org/grpc/metadata"
)

// Parser converts loads from metadata into a concrete type.
type Parser interface {
	// Parse parses loads from metadata.
	Parse(md metadata.MD) interface{}
}

var parser Parser

// SetParser sets the load parser.
//
// Not mutex-protected, should be called before any gRPC functions.
func SetParser(lr Parser) {
	parser = lr
}

// Parse calls parser.Read().
func Parse(md metadata.MD) interface{} {
	if parser == nil {
		return nil
	}
	return parser.Parse(md)
}
