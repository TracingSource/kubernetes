package fuzzer

import (
	fuzz "github.com/google/gofuzz"
	runtimeserializer "k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/kubernetes/pkg/apis/networking"
)

// Funcs returns the fuzzer functions for the networking api group.
var Funcs = func(codecs runtimeserializer.CodecFactory) []interface{} {
	return []interface{}{
		func(np *networking.NetworkPolicyPeer, c fuzz.Continue) {
			c.FuzzNoCustom(np) // fuzz self without calling this function again
			// TODO: Implement a fuzzer to generate valid keys, values and operators for
			// selector requirements.
			if np.IPBlock != nil {
				np.IPBlock = &networking.IPBlock{
					CIDR:   "192.168.1.0/24",
					Except: []string{"192.168.1.1/24", "192.168.1.2/24"},
				}
			}
		},
		func(np *networking.NetworkPolicy, c fuzz.Continue) {
			c.FuzzNoCustom(np) // fuzz self without calling this function again
			// TODO: Implement a fuzzer to generate valid keys, values and operators for
			// selector requirements.
			if len(np.Spec.PolicyTypes) == 0 {
				np.Spec.PolicyTypes = []networking.PolicyType{networking.PolicyTypeIngress}
			}
		},
	}
}
