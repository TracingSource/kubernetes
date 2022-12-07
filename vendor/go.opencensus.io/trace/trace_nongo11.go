// +build !go1.11

package trace

import (
	"context"
)

func startExecutionTracerTask(ctx context.Context, name string) (context.Context, func()) {
	return ctx, func() {}
}
