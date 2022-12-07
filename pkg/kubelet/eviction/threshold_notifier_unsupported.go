// +build !linux

package eviction

import "k8s.io/klog"

// NewCgroupNotifier creates a cgroup notifier that does nothing because cgroups do not exist on non-linux systems.
func NewCgroupNotifier(path, attribute string, threshold int64) (CgroupNotifier, error) {
	klog.V(5).Infof("cgroup notifications not supported")
	return &unsupportedThresholdNotifier{}, nil
}

type unsupportedThresholdNotifier struct{}

func (*unsupportedThresholdNotifier) Start(_ chan<- struct{}) {}

func (*unsupportedThresholdNotifier) Stop() {}
