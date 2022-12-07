package container

// TestRunTimeCache embeds runtimeCache with some additional methods for testing.
// It must be declared in the container package to have visibility to runtimeCache.
// It cannot be in a "..._test.go" file in order for runtime_cache_test.go to have cross-package visibility to it.
// (cross-package declarations in test files cannot be used from dot imports if this package is vendored)
type TestRuntimeCache struct {
	runtimeCache
}

func (r *TestRuntimeCache) UpdateCacheWithLock() error {
	r.Lock()
	defer r.Unlock()
	return r.updateCache()
}

func (r *TestRuntimeCache) GetCachedPods() []*Pod {
	r.Lock()
	defer r.Unlock()
	return r.pods
}

func NewTestRuntimeCache(getter podsGetter) *TestRuntimeCache {
	return &TestRuntimeCache{
		runtimeCache: runtimeCache{
			getter: getter,
		},
	}
}
