package fuzzer

import (
	fuzz "github.com/google/gofuzz"
	"k8s.io/sample-apiserver/pkg/apis/wardle"

	runtimeserializer "k8s.io/apimachinery/pkg/runtime/serializer"
)

// Funcs returns the fuzzer functions for the apps api group.
var Funcs = func(codecs runtimeserializer.CodecFactory) []interface{} {
	return []interface{}{
		func(s *wardle.FlunderSpec, c fuzz.Continue) {
			c.FuzzNoCustom(s) // fuzz self without calling this function again

			if len(s.FlunderReference) != 0 && len(s.FischerReference) != 0 {
				s.FischerReference = ""
			}
			if len(s.FlunderReference) != 0 {
				s.ReferenceType = wardle.FlunderReferenceType
			} else if len(s.FischerReference) != 0 {
				s.ReferenceType = wardle.FischerReferenceType
			} else {
				s.ReferenceType = ""
			}
		},
	}
}
