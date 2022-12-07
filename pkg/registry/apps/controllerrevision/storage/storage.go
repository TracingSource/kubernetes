package storage

import (
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apiserver/pkg/registry/generic"
	genericregistry "k8s.io/apiserver/pkg/registry/generic/registry"
	"k8s.io/kubernetes/pkg/apis/apps"
	"k8s.io/kubernetes/pkg/printers"
	printersinternal "k8s.io/kubernetes/pkg/printers/internalversion"
	printerstorage "k8s.io/kubernetes/pkg/printers/storage"
	"k8s.io/kubernetes/pkg/registry/apps/controllerrevision"
)

// REST implements a RESTStorage for ControllerRevision
type REST struct {
	*genericregistry.Store
}

// NewREST returns a RESTStorage object that will work with ControllerRevision objects.
func NewREST(optsGetter generic.RESTOptionsGetter) (*REST, error) {
	store := &genericregistry.Store{
		NewFunc:                  func() runtime.Object { return &apps.ControllerRevision{} },
		NewListFunc:              func() runtime.Object { return &apps.ControllerRevisionList{} },
		DefaultQualifiedResource: apps.Resource("controllerrevisions"),

		CreateStrategy: controllerrevision.Strategy,
		UpdateStrategy: controllerrevision.Strategy,
		DeleteStrategy: controllerrevision.Strategy,

		TableConvertor: printerstorage.TableConvertor{TableGenerator: printers.NewTableGenerator().With(printersinternal.AddHandlers)},
	}
	options := &generic.StoreOptions{RESTOptions: optsGetter}
	if err := store.CompleteWithOptions(options); err != nil {
		return nil, err
	}
	return &REST{store}, nil
}
