package algorithmprovider

import (
	"k8s.io/kubernetes/pkg/scheduler/algorithmprovider/defaults"
)

// ApplyFeatureGates applies algorithm by feature gates.
func ApplyFeatureGates() func() {
	return defaults.ApplyFeatureGates()
}
