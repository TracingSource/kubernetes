package storage

import (
	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apiserver/pkg/registry/generic"
	genericregistry "k8s.io/apiserver/pkg/registry/generic/registry"
	"k8s.io/apiserver/pkg/registry/rest"
	api "k8s.io/kubernetes/pkg/apis/core"
	"k8s.io/kubernetes/pkg/registry/core/resourcequota"
)

type REST struct {
	*genericregistry.Store
}

// NewREST returns a RESTStorage object that will work against resource quotas.
func NewREST(optsGetter generic.RESTOptionsGetter) (*REST, *StatusREST, error) {
	store := &genericregistry.Store{
		NewFunc:                  func() runtime.Object { return &api.ResourceQuota{} },
		NewListFunc:              func() runtime.Object { return &api.ResourceQuotaList{} },
		DefaultQualifiedResource: api.Resource("resourcequotas"),

		CreateStrategy:      resourcequota.Strategy,
		UpdateStrategy:      resourcequota.Strategy,
		DeleteStrategy:      resourcequota.Strategy,
		ReturnDeletedObject: true,
	}
	options := &generic.StoreOptions{RESTOptions: optsGetter}
	if err := store.CompleteWithOptions(options); err != nil {
		return nil, nil, err
	}

	statusStore := *store
	statusStore.UpdateStrategy = resourcequota.StatusStrategy

	return &REST{store}, &StatusREST{store: &statusStore}, nil
}

// Implement ShortNamesProvider
var _ rest.ShortNamesProvider = &REST{}

// ShortNames implements the ShortNamesProvider interface. Returns a list of short names for a resource.
func (r *REST) ShortNames() []string {
	return []string{"quota"}
}

// StatusREST implements the REST endpoint for changing the status of a resourcequota.
type StatusREST struct {
	store *genericregistry.Store
}

func (r *StatusREST) New() runtime.Object {
	return &api.ResourceQuota{}
}

// Get retrieves the object from the storage. It is required to support Patch.
func (r *StatusREST) Get(ctx context.Context, name string, options *metav1.GetOptions) (runtime.Object, error) {
	return r.store.Get(ctx, name, options)
}

// Update alters the status subset of an object.
func (r *StatusREST) Update(ctx context.Context, name string, objInfo rest.UpdatedObjectInfo, createValidation rest.ValidateObjectFunc, updateValidation rest.ValidateObjectUpdateFunc, forceAllowCreate bool, options *metav1.UpdateOptions) (runtime.Object, bool, error) {
	// We are explicitly setting forceAllowCreate to false in the call to the underlying storage because
	// subresources should never allow create on update.
	return r.store.Update(ctx, name, objInfo, createValidation, updateValidation, false, options)
}
