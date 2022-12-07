package tag

import (
	"context"
)

// FromContext returns the tag map stored in the context.
func FromContext(ctx context.Context) *Map {
	// The returned tag map shouldn't be mutated.
	ts := ctx.Value(mapCtxKey)
	if ts == nil {
		return nil
	}
	return ts.(*Map)
}

// NewContext creates a new context with the given tag map.
// To propagate a tag map to downstream methods and downstream RPCs, add a tag map
// to the current context. NewContext will return a copy of the current context,
// and put the tag map into the returned one.
// If there is already a tag map in the current context, it will be replaced with m.
func NewContext(ctx context.Context, m *Map) context.Context {
	return context.WithValue(ctx, mapCtxKey, m)
}

type ctxKey struct{}

var mapCtxKey = ctxKey{}
