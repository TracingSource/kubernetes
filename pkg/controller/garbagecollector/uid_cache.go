package garbagecollector

import (
	"sync"

	"github.com/golang/groupcache/lru"
	"k8s.io/apimachinery/pkg/types"
)

// UIDCache is an LRU cache for uid.
type UIDCache struct {
	mutex sync.Mutex
	cache *lru.Cache
}

// NewUIDCache returns a UIDCache.
func NewUIDCache(maxCacheEntries int) *UIDCache {
	return &UIDCache{
		cache: lru.New(maxCacheEntries),
	}
}

// Add adds a uid to the cache.
func (c *UIDCache) Add(uid types.UID) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.cache.Add(uid, nil)
}

// Has returns if a uid is in the cache.
func (c *UIDCache) Has(uid types.UID) bool {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	_, found := c.cache.Get(uid)
	return found
}
