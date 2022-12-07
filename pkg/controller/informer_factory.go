package controller

import (
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/metadata/metadatainformer"
)

// InformerFactory creates informers for each group version resource.
type InformerFactory interface {
	ForResource(resource schema.GroupVersionResource) (informers.GenericInformer, error)
	Start(stopCh <-chan struct{})
}

type informerFactory struct {
	typedInformerFactory    informers.SharedInformerFactory
	metadataInformerFactory metadatainformer.SharedInformerFactory
}

func (i *informerFactory) ForResource(
	resource schema.GroupVersionResource,
) (informers.GenericInformer, error) {
	informer, err := i.typedInformerFactory.ForResource(resource)
	if err != nil {
		return i.metadataInformerFactory.ForResource(resource), nil
	}
	return informer, nil
}

// Start ...
//
// caller:
// 	1. cmd/kube-controller-manager/app/controllermanager.go -> run()
// 	kcm 启动时, 完成选主后由主进程调用
//
func (i *informerFactory) Start(stopCh <-chan struct{}) {
	i.typedInformerFactory.Start(stopCh)
	i.metadataInformerFactory.Start(stopCh)
}

// NewInformerFactory ...
//
// caller: 
// 	1. cmd/kube-controller-manager/app/controllermanager.go -> CreateControllerContext()
//
// NewInformerFactory creates a new InformerFactory which works with both typed
// resources and metadata-only resources
func NewInformerFactory(
	typedInformerFactory informers.SharedInformerFactory, 
	metadataInformerFactory metadatainformer.SharedInformerFactory,
) InformerFactory {
	return &informerFactory{
		typedInformerFactory:    typedInformerFactory,
		metadataInformerFactory: metadataInformerFactory,
	}
}
