package runtime

import (
	"context"
	"net/http"
)

// ClientOperation represents the context for a swagger operation to be submitted to the transport
type ClientOperation struct {
	ID                 string
	Method             string
	PathPattern        string
	ProducesMediaTypes []string
	ConsumesMediaTypes []string
	Schemes            []string
	AuthInfo           ClientAuthInfoWriter
	Params             ClientRequestWriter
	Reader             ClientResponseReader
	Context            context.Context
	Client             *http.Client
}

// A ClientTransport implementor knows how to submit Request objects to some destination
type ClientTransport interface {
	//Submit(string, RequestWriter, ResponseReader, AuthInfoWriter) (interface{}, error)
	Submit(*ClientOperation) (interface{}, error)
}
