package nodevolumelimits

import (
	"context"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/kubernetes/pkg/scheduler/algorithm/predicates"
	"k8s.io/kubernetes/pkg/scheduler/framework/plugins/migration"
	framework "k8s.io/kubernetes/pkg/scheduler/framework/v1alpha1"
	"k8s.io/kubernetes/pkg/scheduler/nodeinfo"
)

// GCEPDLimits is a plugin that checks node volume limits.
type GCEPDLimits struct {
	predicate predicates.FitPredicate
}

var _ framework.FilterPlugin = &GCEPDLimits{}

// GCEPDName is the name of the plugin used in the plugin registry and configurations.
const GCEPDName = "GCEPDLimits"

// Name returns name of the plugin. It is used in logs, etc.
func (pl *GCEPDLimits) Name() string {
	return GCEPDName
}

// Filter invoked at the filter extension point.
func (pl *GCEPDLimits) Filter(ctx context.Context, _ *framework.CycleState, pod *v1.Pod, nodeInfo *nodeinfo.NodeInfo) *framework.Status {
	// metadata is not needed
	_, reasons, err := pl.predicate(pod, nil, nodeInfo)
	return migration.PredicateResultToFrameworkStatus(reasons, err)
}

// NewGCEPD returns function that initializes a new plugin and returns it.
func NewGCEPD(_ *runtime.Unknown, handle framework.FrameworkHandle) (framework.Plugin, error) {
	informerFactory := handle.SharedInformerFactory()
	pvLister := informerFactory.Core().V1().PersistentVolumes().Lister()
	pvcLister := informerFactory.Core().V1().PersistentVolumeClaims().Lister()
	scLister := informerFactory.Storage().V1().StorageClasses().Lister()

	return &GCEPDLimits{
		predicate: predicates.NewMaxPDVolumeCountPredicate(predicates.GCEPDVolumeFilterType, getCSINodeListerIfEnabled(informerFactory), scLister, pvLister, pvcLister),
	}, nil
}
