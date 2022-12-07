/*
 *
 * Copyright 2018 gRPC authors.
 *

 *
 */

// This file contains exported variables/functions that are exported for testing
// only.
//
// An ideal way for this would be to put those in a *_test.go but in binarylog
// package. But this doesn't work with staticcheck with go module. Error was:
// "MdToMetadataProto not declared by package binarylog". This could be caused
// by the way staticcheck looks for files for a certain package, which doesn't
// support *_test.go files.
//
// Move those to binary_test.go when staticcheck is fixed.

package binarylog

var (
	// AllLogger is a logger that logs all headers/messages for all RPCs. It's
	// for testing only.
	AllLogger = NewLoggerFromConfigString("*")
	// MdToMetadataProto converts metadata to a binary logging proto message.
	// It's for testing only.
	MdToMetadataProto = mdToMetadataProto
	// AddrToProto converts an address to a binary logging proto message. It's
	// for testing only.
	AddrToProto = addrToProto
)
