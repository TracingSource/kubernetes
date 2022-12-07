package storage

import (
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apiserver/pkg/registry/generic"
	genericregistry "k8s.io/apiserver/pkg/registry/generic/registry"
	"k8s.io/kubernetes/pkg/apis/rbac"
	"k8s.io/kubernetes/pkg/registry/rbac/clusterrole"
)

// REST implements a RESTStorage for ClusterRole
type REST struct {
	*genericregistry.Store
}

// NewREST returns a RESTStorage object that will work against ClusterRole objects.
func NewREST(optsGetter generic.RESTOptionsGetter) (*REST, error) {
	store := &genericregistry.Store{
		NewFunc:                  func() runtime.Object { return &rbac.ClusterRole{} },
		NewListFunc:              func() runtime.Object { return &rbac.ClusterRoleList{} },
		DefaultQualifiedResource: rbac.Resource("clusterroles"),

		CreateStrategy: clusterrole.Strategy,
		UpdateStrategy: clusterrole.Strategy,
		DeleteStrategy: clusterrole.Strategy,
	}
	options := &generic.StoreOptions{RESTOptions: optsGetter}
	if err := store.CompleteWithOptions(options); err != nil {
		return nil, err
	}

	return &REST{store}, nil
}
