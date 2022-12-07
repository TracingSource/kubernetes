/*
 *
 * Copyright 2018 gRPC authors.
 *

 *
 */

// Package envconfig contains grpc settings configured by environment variables.
package envconfig

import (
	"os"
	"strings"
)

const (
	prefix   = "GRPC_GO_"
	retryStr = prefix + "RETRY"
)

var (
	// Retry is set if retry is explicitly enabled via "GRPC_GO_RETRY=on".
	Retry = strings.EqualFold(os.Getenv(retryStr), "on")
)
