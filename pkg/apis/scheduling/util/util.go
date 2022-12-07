package util

import (
	utilfeature "k8s.io/apiserver/pkg/util/feature"
	"k8s.io/kubernetes/pkg/apis/scheduling"
	"k8s.io/kubernetes/pkg/features"
)

// DropDisabledFields removes disabled fields from the PriorityClass object.
func DropDisabledFields(class, oldClass *scheduling.PriorityClass) {
	if !utilfeature.DefaultFeatureGate.Enabled(features.NonPreemptingPriority) && !preemptingPriorityInUse(oldClass) {
		class.PreemptionPolicy = nil
	}
}

func preemptingPriorityInUse(oldClass *scheduling.PriorityClass) bool {
	if oldClass == nil {
		return false
	}
	if oldClass.PreemptionPolicy != nil {
		return true
	}
	return false
}
