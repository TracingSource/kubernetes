package validation

import (
	"k8s.io/apimachinery/pkg/util/validation/field"

	resourcequotaapi "k8s.io/kubernetes/plugin/pkg/admission/resourcequota/apis/resourcequota"
)

// ValidateConfiguration validates the configuration.
func ValidateConfiguration(config *resourcequotaapi.Configuration) field.ErrorList {
	allErrs := field.ErrorList{}
	fldPath := field.NewPath("limitedResources")
	for i, limitedResource := range config.LimitedResources {
		idxPath := fldPath.Index(i)
		if len(limitedResource.Resource) == 0 {
			allErrs = append(allErrs, field.Required(idxPath.Child("resource"), ""))
		}
	}
	return allErrs
}
