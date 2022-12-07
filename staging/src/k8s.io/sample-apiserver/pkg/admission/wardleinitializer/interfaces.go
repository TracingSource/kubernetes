package wardleinitializer

import (
	"k8s.io/apiserver/pkg/admission"
	informers "k8s.io/sample-apiserver/pkg/generated/informers/externalversions"
)

// WantsInternalWardleInformerFactory defines a function which sets InformerFactory for admission plugins that need it
type WantsInternalWardleInformerFactory interface {
	SetInternalWardleInformerFactory(informers.SharedInformerFactory)
	admission.InitializationValidator
}
