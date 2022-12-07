package testing

import (
	"time"

	"k8s.io/apimachinery/pkg/types"
	"k8s.io/kubernetes/pkg/kubelet/container"
)

type fakeCache struct {
	runtime container.Runtime
}

func NewFakeCache(runtime container.Runtime) container.Cache {
	return &fakeCache{runtime: runtime}
}

func (c *fakeCache) Get(id types.UID) (*container.PodStatus, error) {
	return c.runtime.GetPodStatus(id, "", "")
}

func (c *fakeCache) GetNewerThan(id types.UID, minTime time.Time) (*container.PodStatus, error) {
	return c.Get(id)
}

func (c *fakeCache) Set(id types.UID, status *container.PodStatus, err error, timestamp time.Time) {
}

func (c *fakeCache) Delete(id types.UID) {
}

func (c *fakeCache) UpdateTime(_ time.Time) {
}
