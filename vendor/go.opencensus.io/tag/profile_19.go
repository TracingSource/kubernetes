// Copyright 2018, OpenCensus Authors
//


// +build go1.9

package tag

import (
	"context"
	"runtime/pprof"
)

func do(ctx context.Context, f func(ctx context.Context)) {
	m := FromContext(ctx)
	keyvals := make([]string, 0, 2*len(m.m))
	for k, v := range m.m {
		keyvals = append(keyvals, k.Name(), v.value)
	}
	pprof.Do(ctx, pprof.Labels(keyvals...), f)
}
