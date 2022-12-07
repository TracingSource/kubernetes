// Package install installs the experimental API group, making it available as
// an option to all of the API encoding/decoding machinery.
package install

import (
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/kubernetes/pkg/api/legacyscheme"
	"k8s.io/kubernetes/pkg/apis/flowcontrol"
	flowcontrolv1alpha1 "k8s.io/kubernetes/pkg/apis/flowcontrol/v1alpha1"
)

func init() {
	Install(legacyscheme.Scheme)
}

// Install registers the API group and adds types to a scheme
func Install(scheme *runtime.Scheme) {
	utilruntime.Must(flowcontrol.AddToScheme(scheme))
	utilruntime.Must(flowcontrolv1alpha1.AddToScheme(scheme))
	utilruntime.Must(scheme.SetVersionPriority(flowcontrolv1alpha1.SchemeGroupVersion))
}
