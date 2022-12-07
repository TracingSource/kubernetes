package fuzzer

import (
	fuzz "github.com/google/gofuzz"

	runtimeserializer "k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/kubernetes/pkg/apis/certificates"
)

// Funcs returns the fuzzer functions for the certificates api group.
var Funcs = func(codecs runtimeserializer.CodecFactory) []interface{} {
	return []interface{}{
		func(obj *certificates.CertificateSigningRequestSpec, c fuzz.Continue) {
			c.FuzzNoCustom(obj) // fuzz self without calling this function again
			obj.Usages = []certificates.KeyUsage{certificates.UsageKeyEncipherment}
		},
	}
}
