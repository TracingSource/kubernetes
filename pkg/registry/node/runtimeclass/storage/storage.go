package storage

import (
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apiserver/pkg/registry/generic"
	genericregistry "k8s.io/apiserver/pkg/registry/generic/registry"
	"k8s.io/kubernetes/pkg/apis/node"
	"k8s.io/kubernetes/pkg/registry/node/runtimeclass"
)

// REST implements a RESTStorage for RuntimeClass against etcd
type REST struct {
	*genericregistry.Store
}

// NewREST returns a RESTStorage object that will work against runtime classes.
func NewREST(optsGetter generic.RESTOptionsGetter) (*REST, error) {
	store := &genericregistry.Store{
		NewFunc:     func() runtime.Object { return &node.RuntimeClass{} },
		NewListFunc: func() runtime.Object { return &node.RuntimeClassList{} },
		ObjectNameFunc: func(obj runtime.Object) (string, error) {
			return obj.(*node.RuntimeClass).Name, nil
		},
		DefaultQualifiedResource: node.Resource("runtimeclasses"),

		CreateStrategy: runtimeclass.Strategy,
		UpdateStrategy: runtimeclass.Strategy,
		DeleteStrategy: runtimeclass.Strategy,
	}
	options := &generic.StoreOptions{RESTOptions: optsGetter}
	if err := store.CompleteWithOptions(options); err != nil {
		return nil, err
	}
	return &REST{store}, nil
}
