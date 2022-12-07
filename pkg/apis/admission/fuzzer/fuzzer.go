package fuzzer

import (
	fuzz "github.com/google/gofuzz"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	runtimeserializer "k8s.io/apimachinery/pkg/runtime/serializer"
)

// Funcs returns the fuzzer functions for the admission api group.
var Funcs = func(codecs runtimeserializer.CodecFactory) []interface{} {
	return []interface{}{
		func(s *runtime.RawExtension, c fuzz.Continue) {
			u := &unstructured.Unstructured{Object: map[string]interface{}{
				"apiVersion": "unknown.group/unknown",
				"kind":       "Something",
				"somekey":    "somevalue",
			}}
			s.Object = u
		},
	}
}
