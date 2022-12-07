package eventratelimit

import (
	"github.com/hashicorp/golang-lru"

	"k8s.io/client-go/util/flowcontrol"
)

// cache is an interface for caching the limits of a particular type
type cache interface {
	// get the rate limiter associated with the specified key
	get(key interface{}) flowcontrol.RateLimiter
}

// singleCache is a cache that only stores a single, constant item
type singleCache struct {
	// the single rate limiter held by the cache
	rateLimiter flowcontrol.RateLimiter
}

func (c *singleCache) get(key interface{}) flowcontrol.RateLimiter {
	return c.rateLimiter
}

// lruCache is a least-recently-used cache
type lruCache struct {
	// factory to use to create new rate limiters
	rateLimiterFactory func() flowcontrol.RateLimiter
	// the actual LRU cache
	cache *lru.Cache
}

func (c *lruCache) get(key interface{}) flowcontrol.RateLimiter {
	value, found := c.cache.Get(key)
	if !found {
		rateLimter := c.rateLimiterFactory()
		c.cache.Add(key, rateLimter)
		return rateLimter
	}
	return value.(flowcontrol.RateLimiter)
}
