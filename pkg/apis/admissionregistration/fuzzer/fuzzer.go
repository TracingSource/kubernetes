package fuzzer

import (
	fuzz "github.com/google/gofuzz"

	runtimeserializer "k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/kubernetes/pkg/apis/admissionregistration"
)

// Funcs returns the fuzzer functions for the admissionregistration api group.
var Funcs = func(codecs runtimeserializer.CodecFactory) []interface{} {
	return []interface{}{
		func(obj *admissionregistration.Rule, c fuzz.Continue) {
			c.FuzzNoCustom(obj) // fuzz self without calling this function again
			if obj.Scope == nil {
				s := admissionregistration.AllScopes
				obj.Scope = &s
			}
		},
		func(obj *admissionregistration.ValidatingWebhook, c fuzz.Continue) {
			c.FuzzNoCustom(obj) // fuzz self without calling this function again
			p := admissionregistration.FailurePolicyType("Fail")
			obj.FailurePolicy = &p
			m := admissionregistration.MatchPolicyType("Exact")
			obj.MatchPolicy = &m
			s := admissionregistration.SideEffectClassUnknown
			obj.SideEffects = &s
			if obj.TimeoutSeconds == nil {
				i := int32(30)
				obj.TimeoutSeconds = &i
			}
			obj.AdmissionReviewVersions = []string{"v1beta1"}
		},
		func(obj *admissionregistration.MutatingWebhook, c fuzz.Continue) {
			c.FuzzNoCustom(obj) // fuzz self without calling this function again
			p := admissionregistration.FailurePolicyType("Fail")
			obj.FailurePolicy = &p
			m := admissionregistration.MatchPolicyType("Exact")
			obj.MatchPolicy = &m
			s := admissionregistration.SideEffectClassUnknown
			obj.SideEffects = &s
			n := admissionregistration.NeverReinvocationPolicy
			obj.ReinvocationPolicy = &n
			if obj.TimeoutSeconds == nil {
				i := int32(30)
				obj.TimeoutSeconds = &i
			}
			obj.AdmissionReviewVersions = []string{"v1beta1"}
		},
	}
}
