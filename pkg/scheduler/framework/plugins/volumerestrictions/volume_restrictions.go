package volumerestrictions

import (
	"context"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/kubernetes/pkg/scheduler/algorithm/predicates"
	"k8s.io/kubernetes/pkg/scheduler/framework/plugins/migration"
	framework "k8s.io/kubernetes/pkg/scheduler/framework/v1alpha1"
	"k8s.io/kubernetes/pkg/scheduler/nodeinfo"
)

// VolumeRestrictions is a plugin that checks volume restrictions.
type VolumeRestrictions struct{}

var _ framework.FilterPlugin = &VolumeRestrictions{}

// Name is the name of the plugin used in the plugin registry and configurations.
const Name = "VolumeRestrictions"

// Name returns name of the plugin. It is used in logs, etc.
func (pl *VolumeRestrictions) Name() string {
	return Name
}

// Filter invoked at the filter extension point.
func (pl *VolumeRestrictions) Filter(ctx context.Context, _ *framework.CycleState, pod *v1.Pod, nodeInfo *nodeinfo.NodeInfo) *framework.Status {
	// metadata is not needed for NoDiskConflict
	_, reasons, err := predicates.NoDiskConflict(pod, nil, nodeInfo)
	return migration.PredicateResultToFrameworkStatus(reasons, err)
}

// New initializes a new plugin and returns it.
func New(_ *runtime.Unknown, _ framework.FrameworkHandle) (framework.Plugin, error) {
	return &VolumeRestrictions{}, nil
}
