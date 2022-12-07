package podsecuritypolicy

import (
	utilfeature "k8s.io/apiserver/pkg/util/feature"
	"k8s.io/kubernetes/pkg/apis/policy"
	"k8s.io/kubernetes/pkg/features"
)

// DropDisabledFields removes disabled fields from the pod security policy spec.
// This should be called from PrepareForCreate/PrepareForUpdate for all resources containing a od security policy spec.
func DropDisabledFields(pspSpec, oldPSPSpec *policy.PodSecurityPolicySpec) {
	if !utilfeature.DefaultFeatureGate.Enabled(features.ProcMountType) && !allowedProcMountTypesInUse(oldPSPSpec) {
		pspSpec.AllowedProcMountTypes = nil
	}
	if !utilfeature.DefaultFeatureGate.Enabled(features.RunAsGroup) && (oldPSPSpec == nil || oldPSPSpec.RunAsGroup == nil) {
		pspSpec.RunAsGroup = nil
	}
	if !utilfeature.DefaultFeatureGate.Enabled(features.Sysctls) && !sysctlsInUse(oldPSPSpec) {
		pspSpec.AllowedUnsafeSysctls = nil
		pspSpec.ForbiddenSysctls = nil
	}
	if !utilfeature.DefaultFeatureGate.Enabled(features.CSIInlineVolume) {
		pspSpec.AllowedCSIDrivers = nil
	}
	if !utilfeature.DefaultFeatureGate.Enabled(features.RuntimeClass) &&
		(oldPSPSpec == nil || oldPSPSpec.RuntimeClass == nil) {
		pspSpec.RuntimeClass = nil
	}
}

func allowedProcMountTypesInUse(oldPSPSpec *policy.PodSecurityPolicySpec) bool {
	if oldPSPSpec == nil {
		return false
	}

	if oldPSPSpec.AllowedProcMountTypes != nil {
		return true
	}

	return false

}

func sysctlsInUse(oldPSPSpec *policy.PodSecurityPolicySpec) bool {
	if oldPSPSpec == nil {
		return false
	}
	if oldPSPSpec.AllowedUnsafeSysctls != nil || oldPSPSpec.ForbiddenSysctls != nil {
		return true
	}
	return false
}
