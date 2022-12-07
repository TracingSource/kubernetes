package storage

import (
	"testing"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/kubernetes/pkg/apis/policy"
	// Ensure that policy/v1beta1 package is initialized.
	_ "k8s.io/api/policy/v1beta1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apiserver/pkg/registry/generic"
	genericregistrytest "k8s.io/apiserver/pkg/registry/generic/testing"
	etcd3testing "k8s.io/apiserver/pkg/storage/etcd3/testing"
	"k8s.io/kubernetes/pkg/registry/registrytest"
)

func newStorage(t *testing.T) (*REST, *etcd3testing.EtcdTestServer) {
	etcdStorage, server := registrytest.NewEtcdStorage(t, "policy")
	restOptions := generic.RESTOptions{
		StorageConfig:           etcdStorage,
		Decorator:               generic.UndecoratedStorage,
		DeleteCollectionWorkers: 1,
		ResourcePrefix:          "podsecuritypolicies",
	}
	rest, err := NewREST(restOptions)
	if err != nil {
		t.Fatalf("unexpected error from REST storage: %v", err)
	}
	return rest, server
}

func validNewPodSecurityPolicy() *policy.PodSecurityPolicy {
	return &policy.PodSecurityPolicy{
		ObjectMeta: metav1.ObjectMeta{
			Name: "foo",
		},
		Spec: policy.PodSecurityPolicySpec{
			SELinux: policy.SELinuxStrategyOptions{
				Rule: policy.SELinuxStrategyRunAsAny,
			},
			RunAsUser: policy.RunAsUserStrategyOptions{
				Rule: policy.RunAsUserStrategyRunAsAny,
			},
			RunAsGroup: &policy.RunAsGroupStrategyOptions{
				Rule: policy.RunAsGroupStrategyRunAsAny,
			},
			FSGroup: policy.FSGroupStrategyOptions{
				Rule: policy.FSGroupStrategyRunAsAny,
			},
			SupplementalGroups: policy.SupplementalGroupsStrategyOptions{
				Rule: policy.SupplementalGroupsStrategyRunAsAny,
			},
		},
	}
}

func TestCreate(t *testing.T) {
	storage, server := newStorage(t)
	defer server.Terminate(t)
	defer storage.Store.DestroyFunc()
	test := genericregistrytest.New(t, storage.Store).ClusterScope()
	psp := validNewPodSecurityPolicy()
	psp.ObjectMeta = metav1.ObjectMeta{GenerateName: "foo-"}
	test.TestCreate(
		// valid
		psp,
		// invalid
		&policy.PodSecurityPolicy{
			ObjectMeta: metav1.ObjectMeta{Name: "name with spaces"},
		},
	)
}

func TestUpdate(t *testing.T) {
	storage, server := newStorage(t)
	defer server.Terminate(t)
	defer storage.Store.DestroyFunc()
	test := genericregistrytest.New(t, storage.Store).ClusterScope()
	test.TestUpdate(
		// valid
		validNewPodSecurityPolicy(),
		// updateFunc
		func(obj runtime.Object) runtime.Object {
			object := obj.(*policy.PodSecurityPolicy)
			object.Labels = map[string]string{"a": "b"}
			return object
		},
	)
}

func TestDelete(t *testing.T) {
	storage, server := newStorage(t)
	defer server.Terminate(t)
	defer storage.Store.DestroyFunc()
	test := genericregistrytest.New(t, storage.Store).ClusterScope().ReturnDeletedObject()
	test.TestDelete(validNewPodSecurityPolicy())
}

func TestGet(t *testing.T) {
	storage, server := newStorage(t)
	defer server.Terminate(t)
	defer storage.Store.DestroyFunc()
	test := genericregistrytest.New(t, storage.Store).ClusterScope()
	test.TestGet(validNewPodSecurityPolicy())
}

func TestList(t *testing.T) {
	storage, server := newStorage(t)
	defer server.Terminate(t)
	defer storage.Store.DestroyFunc()
	test := genericregistrytest.New(t, storage.Store).ClusterScope()
	test.TestList(validNewPodSecurityPolicy())
}

func TestWatch(t *testing.T) {
	storage, server := newStorage(t)
	defer server.Terminate(t)
	defer storage.Store.DestroyFunc()
	test := genericregistrytest.New(t, storage.Store).ClusterScope()
	test.TestWatch(
		validNewPodSecurityPolicy(),
		// matching labels
		[]labels.Set{},
		// not matching labels
		[]labels.Set{
			{"foo": "bar"},
		},
		// matching fields
		[]fields.Set{
			{"metadata.name": "foo"},
		},
		// not matching fields
		[]fields.Set{
			{"metadata.name": "bar"},
			{"name": "foo"},
		},
	)
}
