package storage

import (
	"testing"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apiserver/pkg/registry/generic"
	genericregistrytest "k8s.io/apiserver/pkg/registry/generic/testing"
	etcd3testing "k8s.io/apiserver/pkg/storage/etcd3/testing"
	api "k8s.io/kubernetes/pkg/apis/core"
	"k8s.io/kubernetes/pkg/registry/registrytest"
)

var testTTL uint64 = 60

func newStorage(t *testing.T) (*REST, *etcd3testing.EtcdTestServer) {
	etcdStorage, server := registrytest.NewEtcdStorage(t, "")
	restOptions := generic.RESTOptions{
		StorageConfig:           etcdStorage,
		Decorator:               generic.UndecoratedStorage,
		DeleteCollectionWorkers: 1,
		ResourcePrefix:          "events",
	}
	rest, err := NewREST(restOptions, testTTL)
	if err != nil {
		t.Fatalf("unexpected error from REST storage: %v", err)
	}
	return rest, server
}

func validNewEvent(namespace string) *api.Event {
	return &api.Event{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "foo",
			Namespace: namespace,
		},
		Reason: "forTesting",
		InvolvedObject: api.ObjectReference{
			Name:      "bar",
			Namespace: namespace,
		},
	}
}

func TestCreate(t *testing.T) {
	storage, server := newStorage(t)
	defer server.Terminate(t)
	defer storage.Store.DestroyFunc()
	test := genericregistrytest.New(t, storage.Store)
	event := validNewEvent(test.TestNamespace())
	event.ObjectMeta = metav1.ObjectMeta{}
	test.TestCreate(
		// valid
		event,
		// invalid
		&api.Event{},
	)
}

func TestUpdate(t *testing.T) {
	storage, server := newStorage(t)
	defer server.Terminate(t)
	defer storage.Store.DestroyFunc()
	test := genericregistrytest.New(t, storage.Store).AllowCreateOnUpdate()
	test.TestUpdate(
		// valid
		validNewEvent(test.TestNamespace()),
		// valid updateFunc
		func(obj runtime.Object) runtime.Object {
			object := obj.(*api.Event)
			object.Reason = "forDifferentTesting"
			return object
		},
		// invalid updateFunc
		func(obj runtime.Object) runtime.Object {
			object := obj.(*api.Event)
			object.InvolvedObject.Namespace = "different-namespace"
			return object
		},
	)
}

func TestDelete(t *testing.T) {
	storage, server := newStorage(t)
	defer server.Terminate(t)
	defer storage.Store.DestroyFunc()
	test := genericregistrytest.New(t, storage.Store)
	test.TestDelete(validNewEvent(test.TestNamespace()))
}

func TestShortNames(t *testing.T) {
	storage, server := newStorage(t)
	defer server.Terminate(t)
	defer storage.Store.DestroyFunc()
	expected := []string{"ev"}
	registrytest.AssertShortNames(t, storage, expected)
}
