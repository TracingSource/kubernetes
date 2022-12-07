package topologymanager

import (
	"k8s.io/kubernetes/pkg/kubelet/lifecycle"
)

type nonePolicy struct{}

var _ Policy = &nonePolicy{}

// PolicyNone policy name.
const PolicyNone string = "none"

// NewNonePolicy returns none policy.
func NewNonePolicy() Policy {
	return &nonePolicy{}
}

func (p *nonePolicy) Name() string {
	return PolicyNone
}

func (p *nonePolicy) canAdmitPodResult(hint *TopologyHint) lifecycle.PodAdmitResult {
	return lifecycle.PodAdmitResult{
		Admit: true,
	}
}

func (p *nonePolicy) Merge(providersHints []map[string][]TopologyHint) (TopologyHint, lifecycle.PodAdmitResult) {
	return TopologyHint{}, p.canAdmitPodResult(nil)
}
