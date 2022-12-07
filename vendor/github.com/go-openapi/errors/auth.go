package errors

import "net/http"

// Unauthenticated returns an unauthenticated error
func Unauthenticated(scheme string) Error {
	return New(http.StatusUnauthorized, "unauthenticated for %s", scheme)
}
