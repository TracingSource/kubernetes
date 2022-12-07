// +build !linux

package pidlimit

import (
	statsapi "k8s.io/kubernetes/pkg/kubelet/apis/stats/v1alpha1"
)

// Stats provides basic information about max and current process count
func Stats() (*statsapi.RlimitStats, error) {
	return nil, nil
}
