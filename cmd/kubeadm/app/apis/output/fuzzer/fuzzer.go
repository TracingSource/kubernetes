package fuzzer

import (
	"time"

	fuzz "github.com/google/gofuzz"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtimeserializer "k8s.io/apimachinery/pkg/runtime/serializer"
	kubeadmapiv1beta2 "k8s.io/kubernetes/cmd/kubeadm/app/apis/kubeadm/v1beta2"
	"k8s.io/kubernetes/cmd/kubeadm/app/apis/output"
)

// Funcs returns the fuzzer functions for the kubeadm apis.
func Funcs(codecs runtimeserializer.CodecFactory) []interface{} {
	return []interface{}{
		fuzzBootstrapToken,
	}
}

func fuzzBootstrapToken(obj *output.BootstrapToken, c fuzz.Continue) {
	c.FuzzNoCustom(obj)

	obj.Token = &kubeadmapiv1beta2.BootstrapTokenString{ID: "uvxdac", Secret: "fq35fuyue3kd4gda"}
	obj.Description = ""
	obj.TTL = &metav1.Duration{Duration: time.Hour * 24}
	obj.Usages = []string{"authentication", "signing"}
	obj.Groups = []string{"system:bootstrappers:kubeadm:default-node-token"}
}
