package priorities

import (
	"fmt"

	v1 "k8s.io/api/core/v1"
	v1helper "k8s.io/kubernetes/pkg/apis/core/v1/helper"
	framework "k8s.io/kubernetes/pkg/scheduler/framework/v1alpha1"
	schedulernodeinfo "k8s.io/kubernetes/pkg/scheduler/nodeinfo"
)

// CountIntolerableTaintsPreferNoSchedule gives the count of intolerable taints of a pod with effect PreferNoSchedule
func countIntolerableTaintsPreferNoSchedule(taints []v1.Taint, tolerations []v1.Toleration) (intolerableTaints int) {
	for _, taint := range taints {
		// check only on taints that have effect PreferNoSchedule
		if taint.Effect != v1.TaintEffectPreferNoSchedule {
			continue
		}

		if !v1helper.TolerationsTolerateTaint(tolerations, &taint) {
			intolerableTaints++
		}
	}
	return
}

// getAllTolerationEffectPreferNoSchedule gets the list of all Tolerations with Effect PreferNoSchedule or with no effect.
func getAllTolerationPreferNoSchedule(tolerations []v1.Toleration) (tolerationList []v1.Toleration) {
	for _, toleration := range tolerations {
		// Empty effect means all effects which includes PreferNoSchedule, so we need to collect it as well.
		if len(toleration.Effect) == 0 || toleration.Effect == v1.TaintEffectPreferNoSchedule {
			tolerationList = append(tolerationList, toleration)
		}
	}
	return
}

// ComputeTaintTolerationPriorityMap prepares the priority list for all the nodes based on the number of intolerable taints on the node
func ComputeTaintTolerationPriorityMap(pod *v1.Pod, meta interface{}, nodeInfo *schedulernodeinfo.NodeInfo) (framework.NodeScore, error) {
	node := nodeInfo.Node()
	if node == nil {
		return framework.NodeScore{}, fmt.Errorf("node not found")
	}
	// To hold all the tolerations with Effect PreferNoSchedule
	var tolerationsPreferNoSchedule []v1.Toleration
	if priorityMeta, ok := meta.(*priorityMetadata); ok {
		tolerationsPreferNoSchedule = priorityMeta.podTolerations

	} else {
		tolerationsPreferNoSchedule = getAllTolerationPreferNoSchedule(pod.Spec.Tolerations)
	}

	return framework.NodeScore{
		Name:  node.Name,
		Score: int64(countIntolerableTaintsPreferNoSchedule(node.Spec.Taints, tolerationsPreferNoSchedule)),
	}, nil
}

// ComputeTaintTolerationPriorityReduce calculates the source of each node based on the number of intolerable taints on the node
var ComputeTaintTolerationPriorityReduce = NormalizeReduce(framework.MaxNodeScore, true)
