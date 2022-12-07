package volumebinding

import (
	"context"

	v1 "k8s.io/api/core/v1"
	"k8s.io/kubernetes/pkg/scheduler/algorithm/predicates"
	"k8s.io/kubernetes/pkg/scheduler/framework/plugins/migration"
	framework "k8s.io/kubernetes/pkg/scheduler/framework/v1alpha1"
	schedulernodeinfo "k8s.io/kubernetes/pkg/scheduler/nodeinfo"
	"k8s.io/kubernetes/pkg/scheduler/volumebinder"
)

// VolumeBinding is a plugin that binds pod volumes in scheduling.
type VolumeBinding struct {
	predicate predicates.FitPredicate
}

var _ framework.FilterPlugin = &VolumeBinding{}

// Name is the name of the plugin used in Registry and configurations.
const Name = "VolumeBinding"

// Name returns name of the plugin. It is used in logs, etc.
func (pl *VolumeBinding) Name() string {
	return Name
}

// Filter invoked at the filter extension point.
func (pl *VolumeBinding) Filter(ctx context.Context, cs *framework.CycleState, pod *v1.Pod, nodeInfo *schedulernodeinfo.NodeInfo) *framework.Status {
	_, reasons, err := pl.predicate(pod, nil, nodeInfo)
	return migration.PredicateResultToFrameworkStatus(reasons, err)
}

// NewFromVolumeBinder initializes a new plugin with volume binder and returns it.
func NewFromVolumeBinder(volumeBinder *volumebinder.VolumeBinder) framework.Plugin {
	return &VolumeBinding{
		predicate: predicates.NewVolumeBindingPredicate(volumeBinder),
	}
}
