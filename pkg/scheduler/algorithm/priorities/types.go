package priorities

import (
	v1 "k8s.io/api/core/v1"
	framework "k8s.io/kubernetes/pkg/scheduler/framework/v1alpha1"
	schedulerlisters "k8s.io/kubernetes/pkg/scheduler/listers"
	schedulernodeinfo "k8s.io/kubernetes/pkg/scheduler/nodeinfo"
)

// PriorityMapFunction is a function that computes per-node results for a given node.
// TODO: Figure out the exact API of this method.
// TODO: Change interface{} to a specific type.
type PriorityMapFunction func(pod *v1.Pod, meta interface{}, nodeInfo *schedulernodeinfo.NodeInfo) (framework.NodeScore, error)

// PriorityReduceFunction is a function that aggregated per-node results and computes
// final scores for all nodes.
// TODO: Figure out the exact API of this method.
// TODO: Change interface{} to a specific type.
type PriorityReduceFunction func(pod *v1.Pod, meta interface{}, sharedLister schedulerlisters.SharedLister, result framework.NodeScoreList) error

// MetadataProducer is a function that computes metadata for a given pod. This
// is now used for only for priority functions. For predicates please use PredicateMetadataProducer.
type MetadataProducer func(pod *v1.Pod, filteredNodes []*v1.Node, sharedLister schedulerlisters.SharedLister) interface{}

// PriorityConfig is a config used for a priority function.
type PriorityConfig struct {
	Name   string
	Map    PriorityMapFunction
	Reduce PriorityReduceFunction
	Weight int64
}

// EmptyMetadataProducer returns a no-op MetadataProducer type.
func EmptyMetadataProducer(pod *v1.Pod, filteredNodes []*v1.Node, sharedLister schedulerlisters.SharedLister) interface{} {
	return nil
}
