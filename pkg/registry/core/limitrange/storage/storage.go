package storage

import (
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apiserver/pkg/registry/generic"
	genericregistry "k8s.io/apiserver/pkg/registry/generic/registry"
	"k8s.io/apiserver/pkg/registry/rest"
	api "k8s.io/kubernetes/pkg/apis/core"
	"k8s.io/kubernetes/pkg/registry/core/limitrange"
)

type REST struct {
	*genericregistry.Store
}

// NewREST returns a RESTStorage object that will work against limitranges.
func NewREST(optsGetter generic.RESTOptionsGetter) (*REST, error) {
	store := &genericregistry.Store{
		NewFunc:                  func() runtime.Object { return &api.LimitRange{} },
		NewListFunc:              func() runtime.Object { return &api.LimitRangeList{} },
		DefaultQualifiedResource: api.Resource("limitranges"),

		CreateStrategy: limitrange.Strategy,
		UpdateStrategy: limitrange.Strategy,
		DeleteStrategy: limitrange.Strategy,
		ExportStrategy: limitrange.Strategy,
	}
	options := &generic.StoreOptions{RESTOptions: optsGetter}
	if err := store.CompleteWithOptions(options); err != nil {
		return nil, err
	}
	return &REST{store}, nil
}

// Implement ShortNamesProvider
var _ rest.ShortNamesProvider = &REST{}

// ShortNames implements the ShortNamesProvider interface. Returns a list of short names for a resource.
func (r *REST) ShortNames() []string {
	return []string{"limits"}
}
