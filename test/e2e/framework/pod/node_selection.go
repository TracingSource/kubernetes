package pod

import (
	v1 "k8s.io/api/core/v1"
)

// NodeSelection specifies where to run a pod, using a combination of fixed node name,
// node selector and/or affinity.
type NodeSelection struct {
	Name     string
	Selector map[string]string
	Affinity *v1.Affinity
}

// setNodeAffinityRequirement sets affinity with specified operator to nodeName to nodeSelection
func setNodeAffinityRequirement(nodeSelection *NodeSelection, operator v1.NodeSelectorOperator, nodeName string) {
	// Add node-anti-affinity.
	if nodeSelection.Affinity == nil {
		nodeSelection.Affinity = &v1.Affinity{}
	}
	if nodeSelection.Affinity.NodeAffinity == nil {
		nodeSelection.Affinity.NodeAffinity = &v1.NodeAffinity{}
	}
	if nodeSelection.Affinity.NodeAffinity.RequiredDuringSchedulingIgnoredDuringExecution == nil {
		nodeSelection.Affinity.NodeAffinity.RequiredDuringSchedulingIgnoredDuringExecution = &v1.NodeSelector{}
	}
	nodeSelection.Affinity.NodeAffinity.RequiredDuringSchedulingIgnoredDuringExecution.NodeSelectorTerms = append(nodeSelection.Affinity.NodeAffinity.RequiredDuringSchedulingIgnoredDuringExecution.NodeSelectorTerms,
		v1.NodeSelectorTerm{
			MatchFields: []v1.NodeSelectorRequirement{
				{Key: "metadata.name", Operator: operator, Values: []string{nodeName}},
			},
		})
}

// SetAffinity sets affinity to nodeName to nodeSelection
func SetAffinity(nodeSelection *NodeSelection, nodeName string) {
	setNodeAffinityRequirement(nodeSelection, v1.NodeSelectorOpIn, nodeName)
}

// SetAntiAffinity sets anti-affinity to nodeName to nodeSelection
func SetAntiAffinity(nodeSelection *NodeSelection, nodeName string) {
	setNodeAffinityRequirement(nodeSelection, v1.NodeSelectorOpNotIn, nodeName)
}
