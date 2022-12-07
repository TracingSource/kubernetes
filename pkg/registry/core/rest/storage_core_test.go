package rest

import (
	"testing"

	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apiserver/pkg/server/storage"
	"k8s.io/apiserver/pkg/storage/storagebackend"
)

func TestGetServersToValidate(t *testing.T) {
	servers := componentStatusStorage{fakeStorageFactory{}}.serversToValidate()

	if e, a := 3, len(servers); e != a {
		t.Errorf("expected %v, got %v", e, a)
	}

	for _, server := range []string{"scheduler", "controller-manager", "etcd-0"} {
		if _, ok := servers[server]; !ok {
			t.Errorf("server list missing: %s", server)
		}
	}
}

type fakeStorageFactory struct{}

func (f fakeStorageFactory) NewConfig(groupResource schema.GroupResource) (*storagebackend.Config, error) {
	return nil, nil
}

func (f fakeStorageFactory) ResourcePrefix(groupResource schema.GroupResource) string {
	return ""
}

func (f fakeStorageFactory) Backends() []storage.Backend {
	return []storage.Backend{{Server: "etcd-0"}}
}
