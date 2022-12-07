package container

import (
	"sync"

	"k8s.io/api/core/v1"
)

// RefManager 就是一个 **容器ID:ObjectReference 对象的 map**, 加上一个读写锁.
// ObjectReference 包含了资源的 Kind, NS, Name等信息.
//
// RefManager manages the references for the containers.
// The references are used for reporting events such as creation,
// failure, etc. This manager is thread-safe, no locks are necessary
// for the caller.
type RefManager struct {
	sync.RWMutex
	containerIDToRef map[ContainerID]*v1.ObjectReference
}

// NewRefManager creates and returns a container reference manager
// with empty contents.
func NewRefManager() *RefManager {
	return &RefManager{
		containerIDToRef: make(map[ContainerID]*v1.ObjectReference),
	}
}

// SetRef stores a reference to a pod's container, associating it with the given container ID.
func (c *RefManager) SetRef(id ContainerID, ref *v1.ObjectReference) {
	c.Lock()
	defer c.Unlock()
	c.containerIDToRef[id] = ref
}

// ClearRef forgets the given container id and its associated container reference.
func (c *RefManager) ClearRef(id ContainerID) {
	c.Lock()
	defer c.Unlock()
	delete(c.containerIDToRef, id)
}

// GetRef returns the container reference of the given ID, or (nil, false) if none is stored.
func (c *RefManager) GetRef(id ContainerID) (ref *v1.ObjectReference, ok bool) {
	c.RLock()
	defer c.RUnlock()
	ref, ok = c.containerIDToRef[id]
	return ref, ok
}
