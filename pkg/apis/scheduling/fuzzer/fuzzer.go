package fuzzer

import (
	"github.com/google/gofuzz"
	runtimeserializer "k8s.io/apimachinery/pkg/runtime/serializer"
	utilfeature "k8s.io/apiserver/pkg/util/feature"
	"k8s.io/kubernetes/pkg/apis/core"
	"k8s.io/kubernetes/pkg/apis/scheduling"
	"k8s.io/kubernetes/pkg/features"
)

// Funcs returns the fuzzer functions for the scheduling api group.
var Funcs = func(codecs runtimeserializer.CodecFactory) []interface{} {
	return []interface{}{
		func(s *scheduling.PriorityClass, c fuzz.Continue) {
			c.FuzzNoCustom(s)
			if utilfeature.DefaultFeatureGate.Enabled(features.NonPreemptingPriority) && s.PreemptionPolicy == nil {
				preemptLowerPriority := core.PreemptLowerPriority
				s.PreemptionPolicy = &preemptLowerPriority
			}
		},
	}
}
