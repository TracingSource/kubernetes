package cloud

import (
	"context"
	"time"
)

const (
	defaultCallTimeout = 1 * time.Hour
)

// ContextWithCallTimeout returns a context with a default timeout, used for generated client calls.
func ContextWithCallTimeout() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), defaultCallTimeout)
}
