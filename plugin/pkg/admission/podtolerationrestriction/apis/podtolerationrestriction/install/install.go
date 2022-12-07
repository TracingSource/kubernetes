// Package install installs the experimental API group, making it available as
// an option to all of the API encoding/decoding machinery.
package install

import (
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	internalapi "k8s.io/kubernetes/plugin/pkg/admission/podtolerationrestriction/apis/podtolerationrestriction"
	versionedapi "k8s.io/kubernetes/plugin/pkg/admission/podtolerationrestriction/apis/podtolerationrestriction/v1alpha1"
)

// Install registers the API group and adds types to a scheme
func Install(scheme *runtime.Scheme) {
	utilruntime.Must(internalapi.AddToScheme(scheme))
	utilruntime.Must(versionedapi.AddToScheme(scheme))
	utilruntime.Must(scheme.SetVersionPriority(versionedapi.SchemeGroupVersion))
}
