package validation

import (
	"fmt"

	"k8s.io/apimachinery/pkg/util/validation/field"
	"k8s.io/kubernetes/pkg/apis/core/validation"
	internalapi "k8s.io/kubernetes/plugin/pkg/admission/podtolerationrestriction/apis/podtolerationrestriction"
)

// ValidateConfiguration validates the configuration.
func ValidateConfiguration(config *internalapi.Configuration) error {
	allErrs := field.ErrorList{}
	fldpath := field.NewPath("podtolerationrestriction")
	allErrs = append(allErrs, validation.ValidateTolerations(config.Default, fldpath.Child("default"))...)
	allErrs = append(allErrs, validation.ValidateTolerations(config.Whitelist, fldpath.Child("whitelist"))...)
	if len(allErrs) > 0 {
		return fmt.Errorf("invalid config: %v", allErrs)
	}
	return nil
}
