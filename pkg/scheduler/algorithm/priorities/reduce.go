package priorities

import (
	"k8s.io/api/core/v1"
	framework "k8s.io/kubernetes/pkg/scheduler/framework/v1alpha1"
	schedulerlisters "k8s.io/kubernetes/pkg/scheduler/listers"
)

// NormalizeReduce generates a PriorityReduceFunction that can normalize the result
// scores to [0, maxPriority]. If reverse is set to true, it reverses the scores by
// subtracting it from maxPriority.
func NormalizeReduce(maxPriority int64, reverse bool) PriorityReduceFunction {
	return func(
		_ *v1.Pod,
		_ interface{},
		_ schedulerlisters.SharedLister,
		result framework.NodeScoreList) error {

		var maxCount int64
		for i := range result {
			if result[i].Score > maxCount {
				maxCount = result[i].Score
			}
		}

		if maxCount == 0 {
			if reverse {
				for i := range result {
					result[i].Score = maxPriority
				}
			}
			return nil
		}

		for i := range result {
			score := result[i].Score

			score = maxPriority * score / maxCount
			if reverse {
				score = maxPriority - score
			}

			result[i].Score = score
		}
		return nil
	}
}
