package storage

import (
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apiserver/pkg/registry/generic"
	genericregistry "k8s.io/apiserver/pkg/registry/generic/registry"
	"k8s.io/kubernetes/pkg/apis/rbac"
	"k8s.io/kubernetes/pkg/registry/rbac/role"
)

// REST implements a RESTStorage for Role
type REST struct {
	*genericregistry.Store
}

// NewREST returns a RESTStorage object that will work against Role objects.
func NewREST(optsGetter generic.RESTOptionsGetter) (*REST, error) {
	store := &genericregistry.Store{
		NewFunc:                  func() runtime.Object { return &rbac.Role{} },
		NewListFunc:              func() runtime.Object { return &rbac.RoleList{} },
		DefaultQualifiedResource: rbac.Resource("roles"),

		CreateStrategy: role.Strategy,
		UpdateStrategy: role.Strategy,
		DeleteStrategy: role.Strategy,
	}
	options := &generic.StoreOptions{RESTOptions: optsGetter}
	if err := store.CompleteWithOptions(options); err != nil {
		return nil, err
	}

	return &REST{store}, nil
}
