package scheme

import (
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/kubernetes/pkg/proxy/apis/config"
	"k8s.io/kubernetes/pkg/proxy/apis/config/v1alpha1"
)

var (
	// Scheme defines methods for serializing and deserializing API objects.
	Scheme = runtime.NewScheme()
	// Codecs provides methods for retrieving codecs and serializers for specific
	// versions and content types.
	Codecs = serializer.NewCodecFactory(Scheme, serializer.EnableStrict)
)

func init() {
	AddToScheme(Scheme)
}

// AddToScheme adds the types of this group into the given scheme.
func AddToScheme(scheme *runtime.Scheme) {
	utilruntime.Must(v1alpha1.AddToScheme(scheme))
	utilruntime.Must(config.AddToScheme(scheme))
}
