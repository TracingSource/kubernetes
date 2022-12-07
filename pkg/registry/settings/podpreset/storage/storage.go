package storage

import (
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apiserver/pkg/registry/generic"
	genericregistry "k8s.io/apiserver/pkg/registry/generic/registry"
	settingsapi "k8s.io/kubernetes/pkg/apis/settings"
	"k8s.io/kubernetes/pkg/registry/settings/podpreset"
)

// rest implements a RESTStorage for replication controllers against etcd
type REST struct {
	*genericregistry.Store
}

// NewREST returns a RESTStorage object that will work against replication controllers.
func NewREST(optsGetter generic.RESTOptionsGetter) (*REST, error) {
	store := &genericregistry.Store{
		NewFunc:                  func() runtime.Object { return &settingsapi.PodPreset{} },
		NewListFunc:              func() runtime.Object { return &settingsapi.PodPresetList{} },
		DefaultQualifiedResource: settingsapi.Resource("podpresets"),

		CreateStrategy: podpreset.Strategy,
		UpdateStrategy: podpreset.Strategy,
		DeleteStrategy: podpreset.Strategy,
	}
	options := &generic.StoreOptions{RESTOptions: optsGetter}
	if err := store.CompleteWithOptions(options); err != nil {
		return nil, err
	}

	return &REST{store}, nil
}
