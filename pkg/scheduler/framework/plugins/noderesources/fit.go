package noderesources

import (
	"context"
	"fmt"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/kubernetes/pkg/scheduler/algorithm/predicates"
	"k8s.io/kubernetes/pkg/scheduler/framework/plugins/migration"
	framework "k8s.io/kubernetes/pkg/scheduler/framework/v1alpha1"
	"k8s.io/kubernetes/pkg/scheduler/nodeinfo"
)

// Fit is a plugin that checks if a node has sufficient resources.
type Fit struct{}

var _ framework.FilterPlugin = &Fit{}

// FitName is the name of the plugin used in the plugin registry and configurations.
const FitName = "NodeResourcesFit"

// Name returns name of the plugin. It is used in logs, etc.
func (f *Fit) Name() string {
	return FitName
}

// Filter invoked at the filter extension point.
func (f *Fit) Filter(ctx context.Context, cycleState *framework.CycleState, pod *v1.Pod, nodeInfo *nodeinfo.NodeInfo) *framework.Status {
	meta, ok := migration.PredicateMetadata(cycleState).(predicates.Metadata)
	if !ok {
		return migration.ErrorToFrameworkStatus(fmt.Errorf("%+v convert to predicates.Metadata error", cycleState))
	}
	_, reasons, err := predicates.PodFitsResources(pod, meta, nodeInfo)
	return migration.PredicateResultToFrameworkStatus(reasons, err)
}

// NewFit initializes a new plugin and returns it.
func NewFit(_ *runtime.Unknown, _ framework.FrameworkHandle) (framework.Plugin, error) {
	return &Fit{}, nil
}
