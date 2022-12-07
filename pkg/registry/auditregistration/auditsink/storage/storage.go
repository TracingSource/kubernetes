package storage

import (
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apiserver/pkg/registry/generic"
	genericregistry "k8s.io/apiserver/pkg/registry/generic/registry"
	"k8s.io/kubernetes/pkg/apis/auditregistration"
	auditstrategy "k8s.io/kubernetes/pkg/registry/auditregistration/auditsink"
)

// REST implements a RESTStorage for audit sink against etcd
type REST struct {
	*genericregistry.Store
}

// NewREST returns a RESTStorage object that will work against audit sinks
func NewREST(optsGetter generic.RESTOptionsGetter) (*REST, error) {
	store := &genericregistry.Store{
		NewFunc:     func() runtime.Object { return &auditregistration.AuditSink{} },
		NewListFunc: func() runtime.Object { return &auditregistration.AuditSinkList{} },
		ObjectNameFunc: func(obj runtime.Object) (string, error) {
			return obj.(*auditregistration.AuditSink).Name, nil
		},
		DefaultQualifiedResource: auditregistration.Resource("auditsinks"),

		CreateStrategy: auditstrategy.Strategy,
		UpdateStrategy: auditstrategy.Strategy,
		DeleteStrategy: auditstrategy.Strategy,
	}
	options := &generic.StoreOptions{RESTOptions: optsGetter}
	if err := store.CompleteWithOptions(options); err != nil {
		return nil, err
	}
	return &REST{store}, nil
}
