package spec

import "sync"

// ResolutionCache a cache for resolving urls
type ResolutionCache interface {
	Get(string) (interface{}, bool)
	Set(string, interface{})
}

type simpleCache struct {
	lock  sync.RWMutex
	store map[string]interface{}
}

// Get retrieves a cached URI
func (s *simpleCache) Get(uri string) (interface{}, bool) {
	debugLog("getting %q from resolution cache", uri)
	s.lock.RLock()
	v, ok := s.store[uri]
	debugLog("got %q from resolution cache: %t", uri, ok)

	s.lock.RUnlock()
	return v, ok
}

// Set caches a URI
func (s *simpleCache) Set(uri string, data interface{}) {
	s.lock.Lock()
	s.store[uri] = data
	s.lock.Unlock()
}

var resCache ResolutionCache

func init() {
	resCache = initResolutionCache()
}

// initResolutionCache initializes the URI resolution cache
func initResolutionCache() ResolutionCache {
	return &simpleCache{store: map[string]interface{}{
		"http://swagger.io/v2/schema.json":       MustLoadSwagger20Schema(),
		"http://json-schema.org/draft-04/schema": MustLoadJSONSchemaDraft04(),
	}}
}
