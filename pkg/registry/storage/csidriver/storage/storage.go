package storage

import (
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apiserver/pkg/registry/generic"
	genericregistry "k8s.io/apiserver/pkg/registry/generic/registry"
	storageapi "k8s.io/kubernetes/pkg/apis/storage"
	"k8s.io/kubernetes/pkg/registry/storage/csidriver"
)

// CSIDriverStorage includes storage for CSIDrivers and all subresources
type CSIDriverStorage struct {
	CSIDriver *REST
}

// REST object that will work for CSIDrivers
type REST struct {
	*genericregistry.Store
}

// NewStorage returns a RESTStorage object that will work against CSIDrivers
func NewStorage(optsGetter generic.RESTOptionsGetter) (*CSIDriverStorage, error) {
	store := &genericregistry.Store{
		NewFunc:                  func() runtime.Object { return &storageapi.CSIDriver{} },
		NewListFunc:              func() runtime.Object { return &storageapi.CSIDriverList{} },
		DefaultQualifiedResource: storageapi.Resource("csidrivers"),

		CreateStrategy:      csidriver.Strategy,
		UpdateStrategy:      csidriver.Strategy,
		DeleteStrategy:      csidriver.Strategy,
		ReturnDeletedObject: true,
	}
	options := &generic.StoreOptions{RESTOptions: optsGetter}
	if err := store.CompleteWithOptions(options); err != nil {
		return nil, err
	}

	return &CSIDriverStorage{
		CSIDriver: &REST{store},
	}, nil
}
