package storage

import (
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apiserver/pkg/registry/generic"
	genericregistry "k8s.io/apiserver/pkg/registry/generic/registry"
	"k8s.io/kubernetes/pkg/apis/admissionregistration"
	"k8s.io/kubernetes/pkg/registry/admissionregistration/validatingwebhookconfiguration"
)

// REST implements a RESTStorage for pod disruption budgets against etcd
type REST struct {
	*genericregistry.Store
}

// NewREST returns a RESTStorage object that will work against pod disruption budgets.
func NewREST(optsGetter generic.RESTOptionsGetter) (*REST, error) {
	store := &genericregistry.Store{
		NewFunc:     func() runtime.Object { return &admissionregistration.ValidatingWebhookConfiguration{} },
		NewListFunc: func() runtime.Object { return &admissionregistration.ValidatingWebhookConfigurationList{} },
		ObjectNameFunc: func(obj runtime.Object) (string, error) {
			return obj.(*admissionregistration.ValidatingWebhookConfiguration).Name, nil
		},
		DefaultQualifiedResource: admissionregistration.Resource("validatingwebhookconfigurations"),

		CreateStrategy: validatingwebhookconfiguration.Strategy,
		UpdateStrategy: validatingwebhookconfiguration.Strategy,
		DeleteStrategy: validatingwebhookconfiguration.Strategy,
	}
	options := &generic.StoreOptions{RESTOptions: optsGetter}
	if err := store.CompleteWithOptions(options); err != nil {
		return nil, err
	}
	return &REST{store}, nil
}
