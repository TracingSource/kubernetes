package controller

import (
	"hash/fnv"
	"sync"

	"github.com/golang/groupcache/lru"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	hashutil "k8s.io/kubernetes/pkg/util/hash"
)

type objectWithMeta interface {
	metav1.Object
}

// keyFunc returns the key of an object, which is used to look up in the cache for it's matching object.
// Since we match objects by namespace and Labels/Selector, so if two objects have the same namespace and labels,
// they will have the same key.
func keyFunc(obj objectWithMeta) uint64 {
	hash := fnv.New32a()
	hashutil.DeepHashObject(hash, &equivalenceLabelObj{
		namespace: obj.GetNamespace(),
		labels:    obj.GetLabels(),
	})
	return uint64(hash.Sum32())
}

type equivalenceLabelObj struct {
	namespace string
	labels    map[string]string
}

// MatchingCache save label and selector matching relationship
type MatchingCache struct {
	mutex sync.RWMutex
	cache *lru.Cache
}

// NewMatchingCache return a NewMatchingCache, which save label and selector matching relationship.
func NewMatchingCache(maxCacheEntries int) *MatchingCache {
	return &MatchingCache{
		cache: lru.New(maxCacheEntries),
	}
}

// Add will add matching information to the cache.
func (c *MatchingCache) Add(labelObj objectWithMeta, selectorObj objectWithMeta) {
	key := keyFunc(labelObj)
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.cache.Add(key, selectorObj)
}

// GetMatchingObject lookup the matching object for a given object.
// Note: the cache information may be invalid since the controller may be deleted or updated,
// we need check in the external request to ensure the cache data is not dirty.
func (c *MatchingCache) GetMatchingObject(labelObj objectWithMeta) (controller interface{}, exists bool) {
	key := keyFunc(labelObj)
	// NOTE: we use Lock() instead of RLock() here because lru's Get() method also modifies state(
	// it need update the least recently usage information). So we can not call it concurrently.
	c.mutex.Lock()
	defer c.mutex.Unlock()
	return c.cache.Get(key)
}

// Update update the cached matching information.
func (c *MatchingCache) Update(labelObj objectWithMeta, selectorObj objectWithMeta) {
	c.Add(labelObj, selectorObj)
}

// InvalidateAll invalidate the whole cache.
func (c *MatchingCache) InvalidateAll() {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.cache = lru.New(c.cache.MaxEntries)
}
