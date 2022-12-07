package storage

import (
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apiserver/pkg/registry/generic"
	genericregistry "k8s.io/apiserver/pkg/registry/generic/registry"
	"k8s.io/kubernetes/pkg/apis/discovery"
	"k8s.io/kubernetes/pkg/printers"
	printersinternal "k8s.io/kubernetes/pkg/printers/internalversion"
	printerstorage "k8s.io/kubernetes/pkg/printers/storage"
	"k8s.io/kubernetes/pkg/registry/discovery/endpointslice"
)

// REST implements a RESTStorage for EndpointSlice against etcd
type REST struct {
	*genericregistry.Store
}

// NewREST returns a RESTStorage object that will work against endpoint slices.
func NewREST(optsGetter generic.RESTOptionsGetter) (*REST, error) {
	store := &genericregistry.Store{
		NewFunc:                  func() runtime.Object { return &discovery.EndpointSlice{} },
		NewListFunc:              func() runtime.Object { return &discovery.EndpointSliceList{} },
		DefaultQualifiedResource: discovery.Resource("endpointslices"),

		CreateStrategy: endpointslice.Strategy,
		UpdateStrategy: endpointslice.Strategy,
		DeleteStrategy: endpointslice.Strategy,

		TableConvertor: printerstorage.TableConvertor{TableGenerator: printers.NewTableGenerator().With(printersinternal.AddHandlers)},
	}
	options := &generic.StoreOptions{RESTOptions: optsGetter}
	if err := store.CompleteWithOptions(options); err != nil {
		return nil, err
	}
	return &REST{store}, nil
}
