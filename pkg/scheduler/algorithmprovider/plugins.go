package algorithmprovider

import (
	"k8s.io/kubernetes/pkg/scheduler/algorithmprovider/defaults"
)

// caller:
// 	1. cmd/kube-scheduler/app/server.go -> runCommand()
//  kube-scheduler 启动时被调用.
//
// ApplyFeatureGates applies algorithm by feature gates.
func ApplyFeatureGates() func() {
	return defaults.ApplyFeatureGates()
}
