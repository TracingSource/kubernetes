package install

import (
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/sample-apiserver/pkg/apis/wardle"
	"k8s.io/sample-apiserver/pkg/apis/wardle/v1alpha1"
	"k8s.io/sample-apiserver/pkg/apis/wardle/v1beta1"
)

// Install registers the API group and adds types to a scheme
func Install(scheme *runtime.Scheme) {
	utilruntime.Must(wardle.AddToScheme(scheme))
	utilruntime.Must(v1beta1.AddToScheme(scheme))
	utilruntime.Must(v1alpha1.AddToScheme(scheme))
	utilruntime.Must(scheme.SetVersionPriority(v1beta1.SchemeGroupVersion, v1alpha1.SchemeGroupVersion))
}
