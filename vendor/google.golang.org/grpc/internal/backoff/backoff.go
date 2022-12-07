// Package backoff implement the backoff strategy for gRPC.
//
// This is kept in internal until the gRPC project decides whether or not to
// allow alternative backoff strategies.
package backoff

import (
	"time"

	"google.golang.org/grpc/internal/grpcrand"
)

// Strategy defines the methodology for backing off after a grpc connection
// failure.
//
type Strategy interface {
	// Backoff returns the amount of time to wait before the next retry given
	// the number of consecutive failures.
	Backoff(retries int) time.Duration
}

const (
	// baseDelay is the amount of time to wait before retrying after the first
	// failure.
	baseDelay = 1.0 * time.Second
	// factor is applied to the backoff after each retry.
	factor = 1.6
	// jitter provides a range to randomize backoff delays.
	jitter = 0.2
)

// Exponential implements exponential backoff algorithm as defined in
// https://github.com/grpc/grpc/blob/master/doc/connection-backoff.md.
type Exponential struct {
	// MaxDelay is the upper bound of backoff delay.
	MaxDelay time.Duration
}

// Backoff returns the amount of time to wait before the next retry given the
// number of retries.
func (bc Exponential) Backoff(retries int) time.Duration {
	if retries == 0 {
		return baseDelay
	}
	backoff, max := float64(baseDelay), float64(bc.MaxDelay)
	for backoff < max && retries > 0 {
		backoff *= factor
		retries--
	}
	if backoff > max {
		backoff = max
	}
	// Randomize backoff delays so that if a cluster of requests start at
	// the same time, they won't operate in lockstep.
	backoff *= 1 + jitter*(grpcrand.Float64()*2-1)
	if backoff < 0 {
		return 0
	}
	return time.Duration(backoff)
}
