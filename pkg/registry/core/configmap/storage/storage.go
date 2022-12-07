package storage

import (
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apiserver/pkg/registry/generic"
	genericregistry "k8s.io/apiserver/pkg/registry/generic/registry"
	"k8s.io/apiserver/pkg/registry/rest"
	"k8s.io/apiserver/pkg/storage"
	api "k8s.io/kubernetes/pkg/apis/core"
	"k8s.io/kubernetes/pkg/printers"
	printersinternal "k8s.io/kubernetes/pkg/printers/internalversion"
	printerstorage "k8s.io/kubernetes/pkg/printers/storage"
	"k8s.io/kubernetes/pkg/registry/core/configmap"
)

// REST implements a RESTStorage for ConfigMap
type REST struct {
	*genericregistry.Store
}

// NewREST returns a RESTStorage object that will work with ConfigMap objects.
func NewREST(optsGetter generic.RESTOptionsGetter) (*REST, error) {
	store := &genericregistry.Store{
		NewFunc:                  func() runtime.Object { return &api.ConfigMap{} },
		NewListFunc:              func() runtime.Object { return &api.ConfigMapList{} },
		DefaultQualifiedResource: api.Resource("configmaps"),

		CreateStrategy: configmap.Strategy,
		UpdateStrategy: configmap.Strategy,
		DeleteStrategy: configmap.Strategy,

		TableConvertor: printerstorage.TableConvertor{TableGenerator: printers.NewTableGenerator().With(printersinternal.AddHandlers)},
	}
	options := &generic.StoreOptions{
		RESTOptions: optsGetter,
		AttrFunc:    configmap.GetAttrs,
		TriggerFunc: map[string]storage.IndexerFunc{"metadata.name": configmap.NameTriggerFunc},
	}
	if err := store.CompleteWithOptions(options); err != nil {
		return nil, err
	}
	return &REST{store}, nil
}

// Implement ShortNamesProvider
var _ rest.ShortNamesProvider = &REST{}

// ShortNames implements the ShortNamesProvider interface. Returns a list of short names for a resource.
func (r *REST) ShortNames() []string {
	return []string{"cm"}
}
