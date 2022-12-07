package streaming

import (
	"net/http"
	"strconv"

	"google.golang.org/grpc/codes"
	grpcstatus "google.golang.org/grpc/status"
)

// NewErrorStreamingDisabled creates an error for disabled streaming method.
func NewErrorStreamingDisabled(method string) error {
	return grpcstatus.Errorf(codes.NotFound, "streaming method %s disabled", method)
}

// NewErrorTooManyInFlight creates an error for exceeding the maximum number of in-flight requests.
func NewErrorTooManyInFlight() error {
	return grpcstatus.Error(codes.ResourceExhausted, "maximum number of in-flight requests exceeded")
}

// WriteError translates a CRI streaming error into an appropriate HTTP response.
func WriteError(err error, w http.ResponseWriter) error {
	s, _ := grpcstatus.FromError(err)
	var status int
	switch s.Code() {
	case codes.NotFound:
		status = http.StatusNotFound
	case codes.ResourceExhausted:
		// We only expect to hit this if there is a DoS, so we just wait the full TTL.
		// If this is ever hit in steady-state operations, consider increasing the maxInFlight requests,
		// or plumbing through the time to next expiration.
		w.Header().Set("Retry-After", strconv.Itoa(int(cacheTTL.Seconds())))
		status = http.StatusTooManyRequests
	default:
		status = http.StatusInternalServerError
	}
	w.WriteHeader(status)
	_, writeErr := w.Write([]byte(err.Error()))
	return writeErr
}
