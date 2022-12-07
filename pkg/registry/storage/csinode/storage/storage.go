package storage

import (
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apiserver/pkg/registry/generic"
	genericregistry "k8s.io/apiserver/pkg/registry/generic/registry"
	storageapi "k8s.io/kubernetes/pkg/apis/storage"
	"k8s.io/kubernetes/pkg/registry/storage/csinode"
)

// CSINodeStorage includes storage for CSINodes and all subresources
type CSINodeStorage struct {
	CSINode *REST
}

// REST object that will work for CSINodes
type REST struct {
	*genericregistry.Store
}

// NewStorage returns a RESTStorage object that will work against CSINodes
func NewStorage(optsGetter generic.RESTOptionsGetter) (*CSINodeStorage, error) {
	store := &genericregistry.Store{
		NewFunc:                  func() runtime.Object { return &storageapi.CSINode{} },
		NewListFunc:              func() runtime.Object { return &storageapi.CSINodeList{} },
		DefaultQualifiedResource: storageapi.Resource("csinodes"),

		CreateStrategy:      csinode.Strategy,
		UpdateStrategy:      csinode.Strategy,
		DeleteStrategy:      csinode.Strategy,
		ReturnDeletedObject: true,
	}
	options := &generic.StoreOptions{RESTOptions: optsGetter}
	if err := store.CompleteWithOptions(options); err != nil {
		return nil, err
	}

	return &CSINodeStorage{
		CSINode: &REST{store},
	}, nil
}
