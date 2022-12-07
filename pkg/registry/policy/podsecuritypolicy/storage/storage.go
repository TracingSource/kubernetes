package storage

import (
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apiserver/pkg/registry/generic"
	genericregistry "k8s.io/apiserver/pkg/registry/generic/registry"
	"k8s.io/kubernetes/pkg/apis/policy"
	"k8s.io/kubernetes/pkg/printers"
	printersinternal "k8s.io/kubernetes/pkg/printers/internalversion"
	printerstorage "k8s.io/kubernetes/pkg/printers/storage"
	"k8s.io/kubernetes/pkg/registry/policy/podsecuritypolicy"
)

// REST implements a RESTStorage for PodSecurityPolicies.
type REST struct {
	*genericregistry.Store
}

// NewREST returns a RESTStorage object that will work against PodSecurityPolicy objects.
func NewREST(optsGetter generic.RESTOptionsGetter) (*REST, error) {
	store := &genericregistry.Store{
		NewFunc:                  func() runtime.Object { return &policy.PodSecurityPolicy{} },
		NewListFunc:              func() runtime.Object { return &policy.PodSecurityPolicyList{} },
		DefaultQualifiedResource: policy.Resource("podsecuritypolicies"),

		CreateStrategy:      podsecuritypolicy.Strategy,
		UpdateStrategy:      podsecuritypolicy.Strategy,
		DeleteStrategy:      podsecuritypolicy.Strategy,
		ReturnDeletedObject: true,

		TableConvertor: printerstorage.TableConvertor{TableGenerator: printers.NewTableGenerator().With(printersinternal.AddHandlers)},
	}
	options := &generic.StoreOptions{RESTOptions: optsGetter}
	if err := store.CompleteWithOptions(options); err != nil {
		return nil, err
	}
	return &REST{store}, nil
}

// ShortNames implements the ShortNamesProvider interface. Returns a list of short names for a resource.
func (r *REST) ShortNames() []string {
	return []string{"psp"}
}
