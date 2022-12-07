package storage

import (
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apiserver/pkg/registry/generic"
	genericregistry "k8s.io/apiserver/pkg/registry/generic/registry"
	coordinationapi "k8s.io/kubernetes/pkg/apis/coordination"
	"k8s.io/kubernetes/pkg/printers"
	printersinternal "k8s.io/kubernetes/pkg/printers/internalversion"
	printerstorage "k8s.io/kubernetes/pkg/printers/storage"
	"k8s.io/kubernetes/pkg/registry/coordination/lease"
)

// REST implements a RESTStorage for leases against etcd
type REST struct {
	*genericregistry.Store
}

// NewREST returns a RESTStorage object that will work against leases.
func NewREST(optsGetter generic.RESTOptionsGetter) (*REST, error) {
	store := &genericregistry.Store{
		NewFunc:                  func() runtime.Object { return &coordinationapi.Lease{} },
		NewListFunc:              func() runtime.Object { return &coordinationapi.LeaseList{} },
		DefaultQualifiedResource: coordinationapi.Resource("leases"),

		CreateStrategy: lease.Strategy,
		UpdateStrategy: lease.Strategy,
		DeleteStrategy: lease.Strategy,

		TableConvertor: printerstorage.TableConvertor{TableGenerator: printers.NewTableGenerator().With(printersinternal.AddHandlers)},
	}
	options := &generic.StoreOptions{RESTOptions: optsGetter}
	if err := store.CompleteWithOptions(options); err != nil {
		return nil, err
	}

	return &REST{store}, nil
}
