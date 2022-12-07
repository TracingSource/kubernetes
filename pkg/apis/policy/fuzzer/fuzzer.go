package fuzzer

import (
	fuzz "github.com/google/gofuzz"

	runtimeserializer "k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/kubernetes/pkg/apis/policy"
)

// Funcs returns the fuzzer functions for the policy api group.
var Funcs = func(codecs runtimeserializer.CodecFactory) []interface{} {
	return []interface{}{
		func(s *policy.PodDisruptionBudgetStatus, c fuzz.Continue) {
			c.FuzzNoCustom(s) // fuzz self without calling this function again
			s.PodDisruptionsAllowed = int32(c.Rand.Intn(2))
		},
		func(psp *policy.PodSecurityPolicySpec, c fuzz.Continue) {
			c.FuzzNoCustom(psp) // fuzz self without calling this function again

			runAsUserRules := []policy.RunAsUserStrategy{
				policy.RunAsUserStrategyMustRunAsNonRoot,
				policy.RunAsUserStrategyMustRunAs,
				policy.RunAsUserStrategyRunAsAny,
			}
			psp.RunAsUser.Rule = runAsUserRules[c.Rand.Intn(len(runAsUserRules))]

			runAsGroupRules := []policy.RunAsGroupStrategy{
				policy.RunAsGroupStrategyMustRunAs,
				policy.RunAsGroupStrategyRunAsAny,
				policy.RunAsGroupStrategyMayRunAs,
			}
			psp.RunAsGroup = &policy.RunAsGroupStrategyOptions{}
			psp.RunAsGroup.Rule = runAsGroupRules[c.Rand.Intn(len(runAsGroupRules))]
			seLinuxRules := []policy.SELinuxStrategy{
				policy.SELinuxStrategyMustRunAs,
				policy.SELinuxStrategyRunAsAny,
			}
			psp.SELinux.Rule = seLinuxRules[c.Rand.Intn(len(seLinuxRules))]

			supplementalGroupsRules := []policy.SupplementalGroupsStrategyType{
				policy.SupplementalGroupsStrategyRunAsAny,
				policy.SupplementalGroupsStrategyMayRunAs,
				policy.SupplementalGroupsStrategyMustRunAs,
			}
			psp.SupplementalGroups.Rule = supplementalGroupsRules[c.Rand.Intn(len(supplementalGroupsRules))]

			fsGroupRules := []policy.FSGroupStrategyType{
				policy.FSGroupStrategyMustRunAs,
				policy.FSGroupStrategyMayRunAs,
				policy.FSGroupStrategyRunAsAny,
			}
			psp.FSGroup.Rule = fsGroupRules[c.Rand.Intn(len(fsGroupRules))]
		},
	}
}
