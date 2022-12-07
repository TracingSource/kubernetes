package scheme

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/kubernetes/cmd/kubeadm/app/apis/kubeadm"
	"k8s.io/kubernetes/cmd/kubeadm/app/apis/kubeadm/v1beta1"
	"k8s.io/kubernetes/cmd/kubeadm/app/apis/kubeadm/v1beta2"
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
	utilruntime.Must(kubeadm.AddToScheme(scheme))
	utilruntime.Must(v1beta1.AddToScheme(scheme))
	utilruntime.Must(v1beta2.AddToScheme(scheme))
	utilruntime.Must(scheme.SetVersionPriority(v1beta2.SchemeGroupVersion))
}
