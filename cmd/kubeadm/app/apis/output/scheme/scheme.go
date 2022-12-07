package scheme

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/kubernetes/cmd/kubeadm/app/apis/output"
	"k8s.io/kubernetes/cmd/kubeadm/app/apis/output/v1alpha1"
)

// Scheme is the runtime.Scheme to which all kubeadm api types are registered.
var Scheme = runtime.NewScheme()

// Codecs provides access to encoding and decoding for the scheme.
var Codecs = serializer.NewCodecFactory(Scheme)

func init() {
	metav1.AddToGroupVersion(Scheme, schema.GroupVersion{Version: "v1"})
	AddToScheme(Scheme)
}

// AddToScheme builds the kubeadm scheme using all known versions of the kubeadm api.
func AddToScheme(scheme *runtime.Scheme) {
	utilruntime.Must(output.AddToScheme(scheme))
	utilruntime.Must(v1alpha1.AddToScheme(scheme))
	utilruntime.Must(scheme.SetVersionPriority(v1alpha1.SchemeGroupVersion))
}
