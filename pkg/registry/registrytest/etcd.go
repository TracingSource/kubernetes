package registrytest

import (
	"testing"

	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apiserver/pkg/server/options"
	serverstorage "k8s.io/apiserver/pkg/server/storage"
	etcd3testing "k8s.io/apiserver/pkg/storage/etcd3/testing"
	"k8s.io/apiserver/pkg/storage/storagebackend"
	"k8s.io/kubernetes/pkg/api/testapi"
	"k8s.io/kubernetes/pkg/kubeapiserver"
)

func NewEtcdStorage(t *testing.T, group string) (*storagebackend.Config, *etcd3testing.EtcdTestServer) {
	server, config := etcd3testing.NewUnsecuredEtcd3TestClientServer(t)
	config.Codec = testapi.Groups[group].StorageCodec()
	return config, server
}

func NewEtcdStorageForResource(t *testing.T, resource schema.GroupResource) (*storagebackend.Config, *etcd3testing.EtcdTestServer) {
	t.Helper()

	server, config := etcd3testing.NewUnsecuredEtcd3TestClientServer(t)

	options := options.NewEtcdOptions(config)
	completedConfig, err := kubeapiserver.NewStorageFactoryConfig().Complete(options)
	if err != nil {
		t.Fatal(err)
	}
	completedConfig.APIResourceConfig = serverstorage.NewResourceConfig()
	factory, err := completedConfig.New()
	if err != nil {
		t.Fatal(err)
	}
	resourceConfig, err := factory.NewConfig(resource)
	if err != nil {
		t.Fatal(err)
	}
	return resourceConfig, server
}
