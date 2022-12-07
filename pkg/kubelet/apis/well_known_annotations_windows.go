// +build windows

package apis

import (
	utilfeature "k8s.io/apiserver/pkg/util/feature"
	"k8s.io/kubernetes/pkg/features"
)

const (
	// HypervIsolationAnnotationKey is used to run windows containers with hyperv isolation.
	// Refer https://aka.ms/hyperv-container.
	HypervIsolationAnnotationKey = "experimental.windows.kubernetes.io/isolation-type"
	// HypervIsolationValue is used to run windows containers with hyperv isolation.
	// Refer https://aka.ms/hyperv-container.
	HypervIsolationValue = "hyperv"
)

// ShouldIsolatedByHyperV returns true if a windows container should be run with hyperv isolation.
func ShouldIsolatedByHyperV(annotations map[string]string) bool {
	if !utilfeature.DefaultFeatureGate.Enabled(features.HyperVContainer) {
		return false
	}

	v, ok := annotations[HypervIsolationAnnotationKey]
	return ok && v == HypervIsolationValue
}
