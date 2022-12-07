package common

import "log"

// Runner is a client or server to run.
type Runner interface {
	// NewOptions returns a new empty options structure to be populated
	// by from the JSON -options argument.
	NewOptions() interface{}
	// Run the client or server, taking in options. This execute the
	// test code.
	Run(logger *log.Logger, options interface{}) error
}
