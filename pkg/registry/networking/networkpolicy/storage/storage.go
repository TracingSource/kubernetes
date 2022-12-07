package storage

import (
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apiserver/pkg/registry/generic"
	genericregistry "k8s.io/apiserver/pkg/registry/generic/registry"
	"k8s.io/apiserver/pkg/registry/rest"
	networkingapi "k8s.io/kubernetes/pkg/apis/networking"
	"k8s.io/kubernetes/pkg/printers"
	printersinternal "k8s.io/kubernetes/pkg/printers/internalversion"
	printerstorage "k8s.io/kubernetes/pkg/printers/storage"
	"k8s.io/kubernetes/pkg/registry/networking/networkpolicy"
)

// rest implements a RESTStorage for NetworkPolicies against etcd
type REST struct {
	*genericregistry.Store
}

// NewREST returns a RESTStorage object that will work against NetworkPolicies
func NewREST(optsGetter generic.RESTOptionsGetter) (*REST, error) {
	store := &genericregistry.Store{
		NewFunc:                  func() runtime.Object { return &networkingapi.NetworkPolicy{} },
		NewListFunc:              func() runtime.Object { return &networkingapi.NetworkPolicyList{} },
		DefaultQualifiedResource: networkingapi.Resource("networkpolicies"),

		CreateStrategy: networkpolicy.Strategy,
		UpdateStrategy: networkpolicy.Strategy,
		DeleteStrategy: networkpolicy.Strategy,

		TableConvertor: printerstorage.TableConvertor{TableGenerator: printers.NewTableGenerator().With(printersinternal.AddHandlers)},
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
	return []string{"netpol"}
}
