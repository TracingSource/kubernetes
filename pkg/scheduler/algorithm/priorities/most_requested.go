package priorities

import framework "k8s.io/kubernetes/pkg/scheduler/framework/v1alpha1"

var (
	mostRequestedRatioResources = DefaultRequestedRatioResources
	mostResourcePriority        = &ResourceAllocationPriority{"MostResourceAllocation", mostResourceScorer, mostRequestedRatioResources}

	// MostRequestedPriorityMap is a priority function that favors nodes with most requested resources.
	// It calculates the percentage of memory and CPU requested by pods scheduled on the node, and prioritizes
	// based on the maximum of the average of the fraction of requested to capacity.
	// Details: (cpu(10 * sum(requested) / capacity) + memory(10 * sum(requested) / capacity)) / 2
	MostRequestedPriorityMap = mostResourcePriority.PriorityMap
)

func mostResourceScorer(requested, allocable ResourceToValueMap, includeVolumes bool, requestedVolumes int, allocatableVolumes int) int64 {
	var nodeScore, weightSum int64
	for resource, weight := range mostRequestedRatioResources {
		resourceScore := mostRequestedScore(requested[resource], allocable[resource])
		nodeScore += resourceScore * weight
		weightSum += weight
	}
	return (nodeScore / weightSum)

}

// The used capacity is calculated on a scale of 0-10
// 0 being the lowest priority and 10 being the highest.
// The more resources are used the higher the score is. This function
// is almost a reversed version of least_requested_priority.calculateUnusedScore
// (10 - calculateUnusedScore). The main difference is in rounding. It was added to
// keep the final formula clean and not to modify the widely used (by users
// in their default scheduling policies) calculateUsedScore.
func mostRequestedScore(requested, capacity int64) int64 {
	if capacity == 0 {
		return 0
	}
	if requested > capacity {
		return 0
	}

	return (requested * framework.MaxNodeScore) / capacity
}
