// Package validation contains methods to validate kinds in the
// authentication.k8s.io API group.
package validation

import (
	"time"

	"k8s.io/apimachinery/pkg/util/validation/field"
	"k8s.io/kubernetes/pkg/apis/authentication"
)

// ValidateTokenRequest validates a TokenRequest.
func ValidateTokenRequest(tr *authentication.TokenRequest) field.ErrorList {
	allErrs := field.ErrorList{}
	specPath := field.NewPath("spec")

	const min = 10 * time.Minute
	if tr.Spec.ExpirationSeconds < int64(min.Seconds()) {
		allErrs = append(allErrs, field.Invalid(specPath.Child("expirationSeconds"), tr.Spec.ExpirationSeconds, "may not specify a duration less than 10 minutes"))
	}
	if tr.Spec.ExpirationSeconds > 1<<32 {
		allErrs = append(allErrs, field.Invalid(specPath.Child("expirationSeconds"), tr.Spec.ExpirationSeconds, "may not specify a duration larger than 2^32 seconds"))
	}
	return allErrs
}
