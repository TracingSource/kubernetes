package storage

import (
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apiserver/pkg/registry/generic"
	genericregistry "k8s.io/apiserver/pkg/registry/generic/registry"
	"k8s.io/apiserver/pkg/storage"
	api "k8s.io/kubernetes/pkg/apis/core"
	"k8s.io/kubernetes/pkg/printers"
	printersinternal "k8s.io/kubernetes/pkg/printers/internalversion"
	printerstorage "k8s.io/kubernetes/pkg/printers/storage"
	"k8s.io/kubernetes/pkg/registry/core/secret"
)

// REST defines the RESTStorage object that will work against secrets.
type REST struct {
	*genericregistry.Store
}

// NewREST returns a RESTStorage object that will work against secrets.
func NewREST(optsGetter generic.RESTOptionsGetter) (*REST, error) {
	store := &genericregistry.Store{
		NewFunc:                  func() runtime.Object { return &api.Secret{} },
		NewListFunc:              func() runtime.Object { return &api.SecretList{} },
		PredicateFunc:            secret.Matcher,
		DefaultQualifiedResource: api.Resource("secrets"),

		CreateStrategy: secret.Strategy,
		UpdateStrategy: secret.Strategy,
		DeleteStrategy: secret.Strategy,
		ExportStrategy: secret.Strategy,

		TableConvertor: printerstorage.TableConvertor{TableGenerator: printers.NewTableGenerator().With(printersinternal.AddHandlers)},
	}
	options := &generic.StoreOptions{
		RESTOptions: optsGetter,
		AttrFunc:    secret.GetAttrs,
		TriggerFunc: map[string]storage.IndexerFunc{"metadata.name": secret.NameTriggerFunc},
	}
	if err := store.CompleteWithOptions(options); err != nil {
		return nil, err
	}
	return &REST{store}, nil
}
