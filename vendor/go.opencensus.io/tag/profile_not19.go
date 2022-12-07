// +build !go1.9

package tag

import "context"

func do(ctx context.Context, f func(ctx context.Context)) {
	f(ctx)
}
