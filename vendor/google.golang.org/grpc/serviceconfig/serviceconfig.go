/*
 *
 * Copyright 2019 gRPC authors.
 *

 *
 */

// Package serviceconfig defines types and methods for operating on gRPC
// service configs.
//
// This package is EXPERIMENTAL.
package serviceconfig

import (
	"google.golang.org/grpc/internal"
)

// Config represents an opaque data structure holding a service config.
type Config interface {
	isConfig()
}

// LoadBalancingConfig represents an opaque data structure holding a load
// balancer config.
type LoadBalancingConfig interface {
	isLoadBalancingConfig()
}

// Parse parses the JSON service config provided into an internal form or
// returns an error if the config is invalid.
func Parse(ServiceConfigJSON string) (Config, error) {
	c, err := internal.ParseServiceConfig(ServiceConfigJSON)
	if err != nil {
		return nil, err
	}
	return c.(Config), err
}
