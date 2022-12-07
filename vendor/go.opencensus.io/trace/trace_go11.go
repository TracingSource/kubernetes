// +build go1.11

package trace

import (
	"context"
	t "runtime/trace"
)

func startExecutionTracerTask(ctx context.Context, name string) (context.Context, func()) {
	if !t.IsEnabled() {
		// Avoid additional overhead if
		// runtime/trace is not enabled.
		return ctx, func() {}
	}
	nctx, task := t.NewTask(ctx, name)
	return nctx, task.End
}
