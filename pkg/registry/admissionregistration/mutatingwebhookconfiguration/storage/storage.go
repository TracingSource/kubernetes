package storage

import (
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apiserver/pkg/registry/generic"
	genericregistry "k8s.io/apiserver/pkg/registry/generic/registry"
	"k8s.io/kubernetes/pkg/apis/admissionregistration"
	"k8s.io/kubernetes/pkg/registry/admissionregistration/mutatingwebhookconfiguration"
)

// REST implements a RESTStorage for pod disruption budgets against etcd
type REST struct {
	*genericregistry.Store
}

// NewREST returns a RESTStorage object that will work against pod disruption budgets.
func NewREST(optsGetter generic.RESTOptionsGetter) (*REST, error) {
	store := &genericregistry.Store{
		NewFunc:     func() runtime.Object { return &admissionregistration.MutatingWebhookConfiguration{} },
		NewListFunc: func() runtime.Object { return &admissionregistration.MutatingWebhookConfigurationList{} },
		ObjectNameFunc: func(obj runtime.Object) (string, error) {
			return obj.(*admissionregistration.MutatingWebhookConfiguration).Name, nil
		},
		DefaultQualifiedResource: admissionregistration.Resource("mutatingwebhookconfigurations"),

		CreateStrategy: mutatingwebhookconfiguration.Strategy,
		UpdateStrategy: mutatingwebhookconfiguration.Strategy,
		DeleteStrategy: mutatingwebhookconfiguration.Strategy,
	}
	options := &generic.StoreOptions{RESTOptions: optsGetter}
	if err := store.CompleteWithOptions(options); err != nil {
		return nil, err
	}
	return &REST{store}, nil
}
