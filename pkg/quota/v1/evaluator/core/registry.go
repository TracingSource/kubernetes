package core

import (
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/clock"
	quota "k8s.io/kubernetes/pkg/quota/v1"
	"k8s.io/kubernetes/pkg/quota/v1/generic"
)

// legacyObjectCountAliases are what we used to do simple object counting quota with mapped to alias
var legacyObjectCountAliases = map[schema.GroupVersionResource]corev1.ResourceName{
	corev1.SchemeGroupVersion.WithResource("configmaps"):             corev1.ResourceConfigMaps,
	corev1.SchemeGroupVersion.WithResource("resourcequotas"):         corev1.ResourceQuotas,
	corev1.SchemeGroupVersion.WithResource("replicationcontrollers"): corev1.ResourceReplicationControllers,
	corev1.SchemeGroupVersion.WithResource("secrets"):                corev1.ResourceSecrets,
}

// NewEvaluators returns the list of static evaluators that manage more than counts
func NewEvaluators(f quota.ListerForResourceFunc) []quota.Evaluator {
	// these evaluators have special logic
	result := []quota.Evaluator{
		NewPodEvaluator(f, clock.RealClock{}),
		NewServiceEvaluator(f),
		NewPersistentVolumeClaimEvaluator(f),
	}
	// these evaluators require an alias for backwards compatibility
	for gvr, alias := range legacyObjectCountAliases {
		result = append(result,
			generic.NewObjectCountEvaluator(gvr.GroupResource(), generic.ListResourceUsingListerFunc(f, gvr), alias))
	}
	return result
}
