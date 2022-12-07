package storage

import (
	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apiserver/pkg/registry/generic"
	genericregistry "k8s.io/apiserver/pkg/registry/generic/registry"
	"k8s.io/apiserver/pkg/registry/rest"
	"k8s.io/kubernetes/pkg/apis/flowcontrol"
	"k8s.io/kubernetes/pkg/printers"
	printersinternal "k8s.io/kubernetes/pkg/printers/internalversion"
	printerstorage "k8s.io/kubernetes/pkg/printers/storage"
	"k8s.io/kubernetes/pkg/registry/flowcontrol/prioritylevelconfiguration"
)

// PriorityLevelConfigurationStorage implements storage for priority level configuration.
type PriorityLevelConfigurationStorage struct {
	PriorityLevelConfiguration *REST
	Status                     *StatusREST
}

// REST implements a RESTStorage for priority level configuration against etcd
type REST struct {
	*genericregistry.Store
}

// NewREST returns a RESTStorage object that will work against priority level configuration.
func NewREST(optsGetter generic.RESTOptionsGetter) (*REST, *StatusREST, error) {
	store := &genericregistry.Store{
		NewFunc:                  func() runtime.Object { return &flowcontrol.PriorityLevelConfiguration{} },
		NewListFunc:              func() runtime.Object { return &flowcontrol.PriorityLevelConfigurationList{} },
		DefaultQualifiedResource: flowcontrol.Resource("prioritylevelconfigurations"),

		CreateStrategy: prioritylevelconfiguration.Strategy,
		UpdateStrategy: prioritylevelconfiguration.Strategy,
		DeleteStrategy: prioritylevelconfiguration.Strategy,

		TableConvertor: printerstorage.TableConvertor{TableGenerator: printers.NewTableGenerator().With(printersinternal.AddHandlers)},
	}
	options := &generic.StoreOptions{RESTOptions: optsGetter}
	if err := store.CompleteWithOptions(options); err != nil {
		return nil, nil, err
	}

	statusStore := *store
	statusStore.CreateStrategy = nil
	statusStore.UpdateStrategy = prioritylevelconfiguration.StatusStrategy
	statusStore.DeleteStrategy = nil

	return &REST{store}, &StatusREST{store: &statusStore}, nil
}

// StatusREST implements the REST endpoint for changing the status of a priority level configuration.
type StatusREST struct {
	store *genericregistry.Store
}

// New creates a new priority level configuration object.
func (r *StatusREST) New() runtime.Object {
	return &flowcontrol.PriorityLevelConfiguration{}
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
