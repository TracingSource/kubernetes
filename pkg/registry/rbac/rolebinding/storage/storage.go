package storage

import (
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apiserver/pkg/registry/generic"
	genericregistry "k8s.io/apiserver/pkg/registry/generic/registry"
	"k8s.io/kubernetes/pkg/apis/rbac"
	"k8s.io/kubernetes/pkg/printers"
	printersinternal "k8s.io/kubernetes/pkg/printers/internalversion"
	printerstorage "k8s.io/kubernetes/pkg/printers/storage"
	"k8s.io/kubernetes/pkg/registry/rbac/rolebinding"
)

// REST implements a RESTStorage for RoleBinding
type REST struct {
	*genericregistry.Store
}

// NewREST returns a RESTStorage object that will work against RoleBinding objects.
func NewREST(optsGetter generic.RESTOptionsGetter) (*REST, error) {
	store := &genericregistry.Store{
		NewFunc:                  func() runtime.Object { return &rbac.RoleBinding{} },
		NewListFunc:              func() runtime.Object { return &rbac.RoleBindingList{} },
		DefaultQualifiedResource: rbac.Resource("rolebindings"),

		CreateStrategy: rolebinding.Strategy,
		UpdateStrategy: rolebinding.Strategy,
		DeleteStrategy: rolebinding.Strategy,

		TableConvertor: printerstorage.TableConvertor{TableGenerator: printers.NewTableGenerator().With(printersinternal.AddHandlers)},
	}
	options := &generic.StoreOptions{RESTOptions: optsGetter}
	if err := store.CompleteWithOptions(options); err != nil {
		return nil, err
	}

	return &REST{store}, nil
}
