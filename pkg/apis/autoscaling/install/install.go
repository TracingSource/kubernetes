// Package install installs the experimental API group, making it available as
// an option to all of the API encoding/decoding machinery.
package install

import (
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/kubernetes/pkg/api/legacyscheme"
	"k8s.io/kubernetes/pkg/apis/autoscaling"
	"k8s.io/kubernetes/pkg/apis/autoscaling/v1"
	"k8s.io/kubernetes/pkg/apis/autoscaling/v2beta1"
	"k8s.io/kubernetes/pkg/apis/autoscaling/v2beta2"
)

func init() {
	Install(legacyscheme.Scheme)
}

// Install registers the API group and adds types to a scheme
func Install(scheme *runtime.Scheme) {
	utilruntime.Must(autoscaling.AddToScheme(scheme))
	utilruntime.Must(v2beta2.AddToScheme(scheme))
	utilruntime.Must(v2beta1.AddToScheme(scheme))
	utilruntime.Must(v1.AddToScheme(scheme))
	utilruntime.Must(scheme.SetVersionPriority(v1.SchemeGroupVersion, v2beta1.SchemeGroupVersion, v2beta2.SchemeGroupVersion))
}
