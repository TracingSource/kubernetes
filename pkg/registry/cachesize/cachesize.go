package cachesize

import (
	"k8s.io/apimachinery/pkg/runtime/schema"
)

// NewHeuristicWatchCacheSizes returns a map of suggested watch cache sizes based on total
// memory.
func NewHeuristicWatchCacheSizes(expectedRAMCapacityMB int) map[schema.GroupResource]int {
	// From our documentation, we officially recommend 120GB machines for
	// 2000 nodes, and we scale from that point. Thus we assume ~60MB of
	// capacity per node.
	// TODO: Revisit this heuristics
	clusterSize := expectedRAMCapacityMB / 60

	// We should specify cache size for a given resource only if it
	// is supposed to have non-default value.
	//
	// TODO: Figure out which resource we should have non-default value.
	watchCacheSizes := make(map[schema.GroupResource]int)
	watchCacheSizes[schema.GroupResource{Resource: "replicationcontrollers"}] = maxInt(5*clusterSize, 100)
	watchCacheSizes[schema.GroupResource{Resource: "endpoints"}] = maxInt(10*clusterSize, 1000)
	watchCacheSizes[schema.GroupResource{Resource: "endpointslices", Group: "discovery.k8s.io"}] = maxInt(10*clusterSize, 1000)
	watchCacheSizes[schema.GroupResource{Resource: "nodes"}] = maxInt(5*clusterSize, 1000)
	watchCacheSizes[schema.GroupResource{Resource: "pods"}] = maxInt(50*clusterSize, 1000)
	watchCacheSizes[schema.GroupResource{Resource: "services"}] = maxInt(5*clusterSize, 1000)
	watchCacheSizes[schema.GroupResource{Resource: "events"}] = 0
	watchCacheSizes[schema.GroupResource{Resource: "apiservices", Group: "apiregistration.k8s.io"}] = maxInt(5*clusterSize, 1000)
	watchCacheSizes[schema.GroupResource{Resource: "leases", Group: "coordination.k8s.io"}] = maxInt(5*clusterSize, 1000)
	return watchCacheSizes
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}
