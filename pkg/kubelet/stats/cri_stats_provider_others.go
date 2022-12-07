// +build !windows

package stats

import (
	statsapi "k8s.io/kubernetes/pkg/kubelet/apis/stats/v1alpha1"
)

// listContainerNetworkStats returns the network stats of all the running containers.
// It should return (nil, nil) for platforms other than Windows.
func (p *criStatsProvider) listContainerNetworkStats() (map[string]*statsapi.NetworkStats, error) {
	return nil, nil
}
