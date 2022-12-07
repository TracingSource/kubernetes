package topologymanager

import (
	"k8s.io/kubernetes/pkg/kubelet/lifecycle"
)

// Policy interface for Topology Manager Pod Admit Result
type Policy interface {
	// Returns Policy Name
	Name() string
	// Returns a merged TopologyHint based on input from hint providers
	// and a Pod Admit Handler Response based on hints and policy type
	Merge(providersHints []map[string][]TopologyHint) (TopologyHint, lifecycle.PodAdmitResult)
}
