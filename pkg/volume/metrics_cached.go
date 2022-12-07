package volume

import (
	"sync"
	"sync/atomic"
)

var _ MetricsProvider = &cachedMetrics{}

// cachedMetrics represents a MetricsProvider that wraps another provider and
// caches the result.
type cachedMetrics struct {
	wrapped       MetricsProvider
	resultError   error
	resultMetrics *Metrics
	once          cacheOnce
}

// NewCachedMetrics creates a new cachedMetrics wrapping another
// MetricsProvider and caching the results.
func NewCachedMetrics(provider MetricsProvider) MetricsProvider {
	return &cachedMetrics{wrapped: provider}
}

// GetMetrics runs the wrapped metrics provider's GetMetrics methd once and
// caches the result. Will not cache result if there is an error.
// See MetricsProvider.GetMetrics
func (md *cachedMetrics) GetMetrics() (*Metrics, error) {
	md.once.cache(func() error {
		md.resultMetrics, md.resultError = md.wrapped.GetMetrics()
		return md.resultError
	})
	return md.resultMetrics, md.resultError
}

// Copied from sync.Once but we don't want to cache the results if there is an
// error
type cacheOnce struct {
	m    sync.Mutex
	done uint32
}

// Copied from sync.Once but we don't want to cache the results if there is an
// error
func (o *cacheOnce) cache(f func() error) {
	if atomic.LoadUint32(&o.done) == 1 {
		return
	}
	// Slow-path.
	o.m.Lock()
	defer o.m.Unlock()
	if o.done == 0 {
		err := f()
		if err == nil {
			atomic.StoreUint32(&o.done, 1)
		}
	}
}
