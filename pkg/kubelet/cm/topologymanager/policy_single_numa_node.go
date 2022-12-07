package topologymanager

import (
	"k8s.io/kubernetes/pkg/kubelet/lifecycle"
)

type singleNumaNodePolicy struct {
	//List of NUMA Nodes available on the underlying machine
	numaNodes []int
}

var _ Policy = &singleNumaNodePolicy{}

// PolicySingleNumaNode policy name.
const PolicySingleNumaNode string = "single-numa-node"

// NewSingleNumaNodePolicy returns single-numa-node policy.
func NewSingleNumaNodePolicy(numaNodes []int) Policy {
	return &singleNumaNodePolicy{numaNodes: numaNodes}
}

func (p *singleNumaNodePolicy) Name() string {
	return PolicySingleNumaNode
}

func (p *singleNumaNodePolicy) canAdmitPodResult(hint *TopologyHint) lifecycle.PodAdmitResult {
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

func (p *singleNumaNodePolicy) Merge(providersHints []map[string][]TopologyHint) (TopologyHint, lifecycle.PodAdmitResult) {
	hint := mergeProvidersHints(p, p.numaNodes, providersHints)
	admit := p.canAdmitPodResult(&hint)
	return hint, admit
}
