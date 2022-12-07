package cache

import (
	"time"

	expirationcache "k8s.io/client-go/tools/cache"
)

// ObjectCache is a simple wrapper of expiration cache that
// 1. use string type key
// 2. has an updater to get value directly if it is expired
// 3. then update the cache
type ObjectCache struct {
	cache   expirationcache.Store
	updater func() (interface{}, error)
}

// objectEntry is an object with string type key.
type objectEntry struct {
	key string
	obj interface{}
}

// NewObjectCache creates ObjectCache with an updater.
// updater returns an object to cache.
func NewObjectCache(f func() (interface{}, error), ttl time.Duration) *ObjectCache {
	return &ObjectCache{
		updater: f,
		cache:   expirationcache.NewTTLStore(stringKeyFunc, ttl),
	}
}

// stringKeyFunc is a string as cache key function
func stringKeyFunc(obj interface{}) (string, error) {
	key := obj.(objectEntry).key
	return key, nil
}

// Get gets cached objectEntry by using a unique string as the key.
func (c *ObjectCache) Get(key string) (interface{}, error) {
	value, ok, err := c.cache.Get(objectEntry{key: key})
	if err != nil {
		return nil, err
	}
	if !ok {
		obj, err := c.updater()
		if err != nil {
			return nil, err
		}
		err = c.cache.Add(objectEntry{
			key: key,
			obj: obj,
		})
		if err != nil {
			return nil, err
		}
		return obj, nil
	}
	return value.(objectEntry).obj, nil
}

// Add adds objectEntry by using a unique string as the key.
func (c *ObjectCache) Add(key string, obj interface{}) error {
	err := c.cache.Add(objectEntry{key: key, obj: obj})
	if err != nil {
		return err
	}
	return nil
}
