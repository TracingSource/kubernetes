package scheme

import (
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	kubeletconfig "k8s.io/kubernetes/pkg/kubelet/apis/config"
	kubeletconfigv1beta1 "k8s.io/kubernetes/pkg/kubelet/apis/config/v1beta1"
)

// Utility functions for the Kubelet's kubeletconfig API group

// NewSchemeAndCodecs is a utility function that returns a Scheme and CodecFactory
// that understand the types in the kubeletconfig API group. Passing mutators allows
// for adjusting the behavior of the CodecFactory, for example enable strict decoding.
func NewSchemeAndCodecs(mutators ...serializer.CodecFactoryOptionsMutator) (*runtime.Scheme, *serializer.CodecFactory, error) {
	scheme := runtime.NewScheme()
	if err := kubeletconfig.AddToScheme(scheme); err != nil {
		return nil, nil, err
	}
	if err := kubeletconfigv1beta1.AddToScheme(scheme); err != nil {
		return nil, nil, err
	}
	codecs := serializer.NewCodecFactory(scheme, mutators...)
	return scheme, &codecs, nil
}
