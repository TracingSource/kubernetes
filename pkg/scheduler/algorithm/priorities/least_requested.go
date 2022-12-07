package priorities

import framework "k8s.io/kubernetes/pkg/scheduler/framework/v1alpha1"

var (
	leastRequestedRatioResources = DefaultRequestedRatioResources
	leastResourcePriority        = &ResourceAllocationPriority{"LeastResourceAllocation", leastResourceScorer, leastRequestedRatioResources}

	// LeastRequestedPriorityMap is a priority function that favors nodes with fewer requested resources.
	// It calculates the percentage of memory and CPU requested by pods scheduled on the node, and
	// prioritizes based on the minimum of the average of the fraction of requested to capacity.
	//
	// Details:
	// (cpu((capacity-sum(requested))*10/capacity) + memory((capacity-sum(requested))*10/capacity))/2
	LeastRequestedPriorityMap = leastResourcePriority.PriorityMap
)

func leastResourceScorer(requested, allocable ResourceToValueMap, includeVolumes bool, requestedVolumes int, allocatableVolumes int) int64 {
	var nodeScore, weightSum int64
	for resource, weight := range leastRequestedRatioResources {
		resourceScore := leastRequestedScore(requested[resource], allocable[resource])
		nodeScore += resourceScore * weight
		weightSum += weight
	}
	return nodeScore / weightSum
}

// The unused capacity is calculated on a scale of 0-10
// 0 being the lowest priority and 10 being the highest.
// The more unused resources the higher the score is.
func leastRequestedScore(requested, capacity int64) int64 {
	if capacity == 0 {
		return 0
	}
	if requested > capacity {
		return 0
	}

	return ((capacity - requested) * int64(framework.MaxNodeScore)) / capacity
}
