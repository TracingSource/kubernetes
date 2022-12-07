package nodeports

import (
	"context"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/kubernetes/pkg/scheduler/algorithm/predicates"
	"k8s.io/kubernetes/pkg/scheduler/framework/plugins/migration"
	framework "k8s.io/kubernetes/pkg/scheduler/framework/v1alpha1"
	"k8s.io/kubernetes/pkg/scheduler/nodeinfo"
)

// NodePorts is a plugin that checks if a node has free ports for the requested pod ports.
type NodePorts struct{}

var _ framework.FilterPlugin = &NodePorts{}

// Name is the name of the plugin used in the plugin registry and configurations.
const Name = "NodePorts"

// Name returns name of the plugin. It is used in logs, etc.
func (pl *NodePorts) Name() string {
	return Name
}

// Filter invoked at the filter extension point.
func (pl *NodePorts) Filter(ctx context.Context, _ *framework.CycleState, pod *v1.Pod, nodeInfo *nodeinfo.NodeInfo) *framework.Status {
	_, reasons, err := predicates.PodFitsHostPorts(pod, nil, nodeInfo)
	return migration.PredicateResultToFrameworkStatus(reasons, err)
}

// New initializes a new plugin and returns it.
func New(_ *runtime.Unknown, _ framework.FrameworkHandle) (framework.Plugin, error) {
	return &NodePorts{}, nil
}
