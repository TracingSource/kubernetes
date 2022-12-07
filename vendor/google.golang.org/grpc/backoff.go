// See internal/backoff package for the backoff implementation. This file is
// kept for the exported types and API backward compatibility.

package grpc

import (
	"time"
)

// DefaultBackoffConfig uses values specified for backoff in
// https://github.com/grpc/grpc/blob/master/doc/connection-backoff.md.
var DefaultBackoffConfig = BackoffConfig{
	MaxDelay: 120 * time.Second,
}

// BackoffConfig defines the parameters for the default gRPC backoff strategy.
type BackoffConfig struct {
	// MaxDelay is the upper bound of backoff delay.
	MaxDelay time.Duration
}
