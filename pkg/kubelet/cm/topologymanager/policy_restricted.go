package topologymanager

import (
	"k8s.io/kubernetes/pkg/kubelet/lifecycle"
)

type restrictedPolicy struct {
	bestEffortPolicy
}

var _ Policy = &restrictedPolicy{}

// PolicyRestricted policy name.
const PolicyRestricted string = "restricted"

// NewRestrictedPolicy returns restricted policy.
func NewRestrictedPolicy(numaNodes []int) Policy {
	return &restrictedPolicy{bestEffortPolicy{numaNodes: numaNodes}}
}

func (p *restrictedPolicy) Name() string {
	return PolicyRestricted
}

func (p *restrictedPolicy) canAdmitPodResult(hint *TopologyHint) lifecycle.PodAdmitResult {
	if !hint.Preferred {
		return lifecycle.PodAdmitResult{
			Admit:   false,
			Reason:  "Topology Affinity Error",
			Message: "Resources cannot be allocated with Topology Locality",
		}
	}
	return lifecycle.PodAdmitResult{
		Admit: true,
	}
}

func (p *restrictedPolicy) Merge(providersHints []map[string][]TopologyHint) (TopologyHint, lifecycle.PodAdmitResult) {
	hint := mergeProvidersHints(p, p.numaNodes, providersHints)
	admit := p.canAdmitPodResult(&hint)
	return hint, admit
}
