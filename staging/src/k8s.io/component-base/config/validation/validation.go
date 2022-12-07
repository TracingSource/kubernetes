package validation

import (
	"k8s.io/apimachinery/pkg/util/validation/field"
	"k8s.io/component-base/config"
)

// ValidateClientConnectionConfiguration ensures validation of the ClientConnectionConfiguration struct
func ValidateClientConnectionConfiguration(cc *config.ClientConnectionConfiguration, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}
	if cc.Burst < 0 {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("burst"), cc.Burst, "must be non-negative"))
	}
	return allErrs
}

// ValidateLeaderElectionConfiguration ensures validation of the LeaderElectionConfiguration struct
func ValidateLeaderElectionConfiguration(cc *config.LeaderElectionConfiguration, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}
	if !cc.LeaderElect {
		return allErrs
	}
	if cc.LeaseDuration.Duration <= 0 {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("leaseDuration"), cc.LeaseDuration, "must be greater than zero"))
	}
	if cc.RenewDeadline.Duration <= 0 {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("renewDeadline"), cc.RenewDeadline, "must be greater than zero"))
	}
	if cc.RetryPeriod.Duration <= 0 {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("retryPeriod"), cc.RetryPeriod, "must be greater than zero"))
	}
	if cc.LeaseDuration.Duration < cc.RenewDeadline.Duration {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("leaseDuration"), cc.RenewDeadline, "LeaseDuration must be greater than RenewDeadline"))
	}
	if len(cc.ResourceLock) == 0 {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("resourceLock"), cc.ResourceLock, "resourceLock is required"))
	}
	if len(cc.ResourceNamespace) == 0 {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("resourceNamespace"), cc.ResourceNamespace, "resourceNamespace is required"))
	}
	if len(cc.ResourceName) == 0 {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("resourceName"), cc.ResourceName, "resourceName is required"))
	}
	return allErrs
}
