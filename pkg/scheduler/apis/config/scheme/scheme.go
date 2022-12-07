package scheme

import (
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	kubeschedulerconfig "k8s.io/kubernetes/pkg/scheduler/apis/config"
	kubeschedulerconfigv1 "k8s.io/kubernetes/pkg/scheduler/apis/config/v1"
	kubeschedulerconfigv1alpha1 "k8s.io/kubernetes/pkg/scheduler/apis/config/v1alpha1"
)

var (
	// Scheme is the runtime.Scheme to which all kubescheduler api types are registered.
	Scheme = runtime.NewScheme()

	// Codecs provides access to encoding and decoding for the scheme.
	Codecs = serializer.NewCodecFactory(Scheme, serializer.EnableStrict)
)

func init() {
	AddToScheme(Scheme)
}

// AddToScheme builds the kubescheduler scheme using all known versions of the kubescheduler api.
func AddToScheme(scheme *runtime.Scheme) {
	utilruntime.Must(kubeschedulerconfig.AddToScheme(Scheme))
	utilruntime.Must(kubeschedulerconfigv1.AddToScheme(Scheme))
	utilruntime.Must(kubeschedulerconfigv1alpha1.AddToScheme(Scheme))
	utilruntime.Must(scheme.SetVersionPriority(kubeschedulerconfigv1alpha1.SchemeGroupVersion))
}
