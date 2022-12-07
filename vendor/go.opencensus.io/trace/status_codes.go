package trace

// Status codes for use with Span.SetStatus. These correspond to the status
// codes used by gRPC defined here: https://github.com/googleapis/googleapis/blob/master/google/rpc/code.proto
const (
	StatusCodeOK                 = 0
	StatusCodeCancelled          = 1
	StatusCodeUnknown            = 2
	StatusCodeInvalidArgument    = 3
	StatusCodeDeadlineExceeded   = 4
	StatusCodeNotFound           = 5
	StatusCodeAlreadyExists      = 6
	StatusCodePermissionDenied   = 7
	StatusCodeResourceExhausted  = 8
	StatusCodeFailedPrecondition = 9
	StatusCodeAborted            = 10
	StatusCodeOutOfRange         = 11
	StatusCodeUnimplemented      = 12
	StatusCodeInternal           = 13
	StatusCodeUnavailable        = 14
	StatusCodeDataLoss           = 15
	StatusCodeUnauthenticated    = 16
)
