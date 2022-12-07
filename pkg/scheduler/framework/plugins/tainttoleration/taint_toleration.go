package tainttoleration

import (
	"context"
	"fmt"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/kubernetes/pkg/scheduler/algorithm/predicates"
	"k8s.io/kubernetes/pkg/scheduler/algorithm/priorities"
	"k8s.io/kubernetes/pkg/scheduler/framework/plugins/migration"
	framework "k8s.io/kubernetes/pkg/scheduler/framework/v1alpha1"
	"k8s.io/kubernetes/pkg/scheduler/nodeinfo"
)

// TaintToleration is a plugin that checks if a pod tolerates a node's taints.
type TaintToleration struct {
	handle framework.FrameworkHandle
}

var _ framework.FilterPlugin = &TaintToleration{}
var _ framework.ScorePlugin = &TaintToleration{}

// Name is the name of the plugin used in the plugin registry and configurations.
const Name = "TaintToleration"

// Name returns name of the plugin. It is used in logs, etc.
func (pl *TaintToleration) Name() string {
	return Name
}

// Filter invoked at the filter extension point.
func (pl *TaintToleration) Filter(ctx context.Context, state *framework.CycleState, pod *v1.Pod, nodeInfo *nodeinfo.NodeInfo) *framework.Status {
	// Note that PodToleratesNodeTaints doesn't use predicate metadata, hence passing nil here.
	_, reasons, err := predicates.PodToleratesNodeTaints(pod, nil, nodeInfo)
	return migration.PredicateResultToFrameworkStatus(reasons, err)
}

// Score invoked at the Score extension point.
func (pl *TaintToleration) Score(ctx context.Context, state *framework.CycleState, pod *v1.Pod, nodeName string) (int64, *framework.Status) {
	nodeInfo, err := pl.handle.SnapshotSharedLister().NodeInfos().Get(nodeName)
	if err != nil {
		return 0, framework.NewStatus(framework.Error, fmt.Sprintf("getting node %q from Snapshot: %v", nodeName, err))
	}
	meta := migration.PriorityMetadata(state)
	s, err := priorities.ComputeTaintTolerationPriorityMap(pod, meta, nodeInfo)
	return s.Score, migration.ErrorToFrameworkStatus(err)
}

// NormalizeScore invoked after scoring all nodes.
func (pl *TaintToleration) NormalizeScore(ctx context.Context, _ *framework.CycleState, pod *v1.Pod, scores framework.NodeScoreList) *framework.Status {
	// Note that ComputeTaintTolerationPriorityReduce doesn't use priority metadata, hence passing nil here.
	err := priorities.ComputeTaintTolerationPriorityReduce(pod, nil, pl.handle.SnapshotSharedLister(), scores)
	return migration.ErrorToFrameworkStatus(err)
}

// ScoreExtensions of the Score plugin.
func (pl *TaintToleration) ScoreExtensions() framework.ScoreExtensions {
	return pl
}

// New initializes a new plugin and returns it.
func New(_ *runtime.Unknown, h framework.FrameworkHandle) (framework.Plugin, error) {
	return &TaintToleration{handle: h}, nil
}
