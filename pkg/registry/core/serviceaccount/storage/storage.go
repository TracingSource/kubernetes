package storage

import (
	"time"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apiserver/pkg/authentication/authenticator"
	"k8s.io/apiserver/pkg/registry/generic"
	genericregistry "k8s.io/apiserver/pkg/registry/generic/registry"
	"k8s.io/apiserver/pkg/registry/rest"
	api "k8s.io/kubernetes/pkg/apis/core"
	"k8s.io/kubernetes/pkg/printers"
	printersinternal "k8s.io/kubernetes/pkg/printers/internalversion"
	printerstorage "k8s.io/kubernetes/pkg/printers/storage"
	"k8s.io/kubernetes/pkg/registry/core/serviceaccount"
	token "k8s.io/kubernetes/pkg/serviceaccount"
)

type REST struct {
	*genericregistry.Store
	Token *TokenREST
}

// NewREST returns a RESTStorage object that will work against service accounts.
func NewREST(optsGetter generic.RESTOptionsGetter, issuer token.TokenGenerator, auds authenticator.Audiences, max time.Duration, podStorage, secretStorage *genericregistry.Store) (*REST, error) {
	store := &genericregistry.Store{
		NewFunc:                  func() runtime.Object { return &api.ServiceAccount{} },
		NewListFunc:              func() runtime.Object { return &api.ServiceAccountList{} },
		DefaultQualifiedResource: api.Resource("serviceaccounts"),

		CreateStrategy:      serviceaccount.Strategy,
		UpdateStrategy:      serviceaccount.Strategy,
		DeleteStrategy:      serviceaccount.Strategy,
		ReturnDeletedObject: true,

		TableConvertor: printerstorage.TableConvertor{TableGenerator: printers.NewTableGenerator().With(printersinternal.AddHandlers)},
	}
	options := &generic.StoreOptions{RESTOptions: optsGetter}
	if err := store.CompleteWithOptions(options); err != nil {
		return nil, err
	}

	var trest *TokenREST
	if issuer != nil && podStorage != nil && secretStorage != nil {
		trest = &TokenREST{
			svcaccts:             store,
			pods:                 podStorage,
			secrets:              secretStorage,
			issuer:               issuer,
			auds:                 auds,
			maxExpirationSeconds: int64(max.Seconds()),
		}
	}

	return &REST{
		Store: store,
		Token: trest,
	}, nil
}

// Implement ShortNamesProvider
var _ rest.ShortNamesProvider = &REST{}

// ShortNames implements the ShortNamesProvider interface. Returns a list of short names for a resource.
func (r *REST) ShortNames() []string {
	return []string{"sa"}
}
