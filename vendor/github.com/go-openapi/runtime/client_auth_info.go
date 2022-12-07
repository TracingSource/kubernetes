package runtime

import "github.com/go-openapi/strfmt"

// A ClientAuthInfoWriterFunc converts a function to a request writer interface
type ClientAuthInfoWriterFunc func(ClientRequest, strfmt.Registry) error

// AuthenticateRequest adds authentication data to the request
func (fn ClientAuthInfoWriterFunc) AuthenticateRequest(req ClientRequest, reg strfmt.Registry) error {
	return fn(req, reg)
}

// A ClientAuthInfoWriter implementor knows how to write authentication info to a request
type ClientAuthInfoWriter interface {
	AuthenticateRequest(ClientRequest, strfmt.Registry) error
}
