package priorities

import (
	"fmt"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/labels"
	v1helper "k8s.io/kubernetes/pkg/apis/core/v1/helper"
	framework "k8s.io/kubernetes/pkg/scheduler/framework/v1alpha1"
	schedulernodeinfo "k8s.io/kubernetes/pkg/scheduler/nodeinfo"
)

// CalculateNodeAffinityPriorityMap prioritizes nodes according to node affinity scheduling preferences
// indicated in PreferredDuringSchedulingIgnoredDuringExecution. Each time a node matches a preferredSchedulingTerm,
// it will get an add of preferredSchedulingTerm.Weight. Thus, the more preferredSchedulingTerms
// the node satisfies and the more the preferredSchedulingTerm that is satisfied weights, the higher
// score the node gets.
func CalculateNodeAffinityPriorityMap(pod *v1.Pod, meta interface{}, nodeInfo *schedulernodeinfo.NodeInfo) (framework.NodeScore, error) {
	node := nodeInfo.Node()
	if node == nil {
		return framework.NodeScore{}, fmt.Errorf("node not found")
	}

	// default is the podspec.
	affinity := pod.Spec.Affinity
	if priorityMeta, ok := meta.(*priorityMetadata); ok {
		// We were able to parse metadata, use affinity from there.
		affinity = priorityMeta.affinity
	}

	var count int32
	// A nil element of PreferredDuringSchedulingIgnoredDuringExecution matches no objects.
	// An element of PreferredDuringSchedulingIgnoredDuringExecution that refers to an
	// empty PreferredSchedulingTerm matches all objects.
	if affinity != nil && affinity.NodeAffinity != nil && affinity.NodeAffinity.PreferredDuringSchedulingIgnoredDuringExecution != nil {
		// Match PreferredDuringSchedulingIgnoredDuringExecution term by term.
		for i := range affinity.NodeAffinity.PreferredDuringSchedulingIgnoredDuringExecution {
			preferredSchedulingTerm := &affinity.NodeAffinity.PreferredDuringSchedulingIgnoredDuringExecution[i]
			if preferredSchedulingTerm.Weight == 0 {
				continue
			}

			// TODO: Avoid computing it for all nodes if this becomes a performance problem.
			nodeSelector, err := v1helper.NodeSelectorRequirementsAsSelector(preferredSchedulingTerm.Preference.MatchExpressions)
			if err != nil {
				return framework.NodeScore{}, err
			}
			if nodeSelector.Matches(labels.Set(node.Labels)) {
				count += preferredSchedulingTerm.Weight
			}
		}
	}

	return framework.NodeScore{
		Name:  node.Name,
		Score: int64(count),
	}, nil
}

// CalculateNodeAffinityPriorityReduce is a reduce function for node affinity priority calculation.
var CalculateNodeAffinityPriorityReduce = NormalizeReduce(framework.MaxNodeScore, false)
