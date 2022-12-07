package validation

import (
	apimachineryvalidation "k8s.io/apimachinery/pkg/api/validation"
	unversionedvalidation "k8s.io/apimachinery/pkg/apis/meta/v1/validation"
	"k8s.io/apimachinery/pkg/util/validation/field"
	apivalidation "k8s.io/kubernetes/pkg/apis/core/validation"
	"k8s.io/kubernetes/pkg/apis/settings"
)

// ValidatePodPresetName can be used to check whether the given PodPreset name is valid.
// Prefix indicates this name will be used as part of generation, in which case
// trailing dashes are allowed.
func ValidatePodPresetName(name string, prefix bool) []string {
	// TODO: Validate that there's name for the suffix inserted by the pods.
	// Currently this is just "-index". In the future we may allow a user
	// specified list of suffixes and we need  to validate the longest one.
	return apimachineryvalidation.NameIsDNSSubdomain(name, prefix)
}

// ValidatePodPresetSpec tests if required fields in the PodPreset spec are set.
func ValidatePodPresetSpec(spec *settings.PodPresetSpec, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	allErrs = append(allErrs, unversionedvalidation.ValidateLabelSelector(&spec.Selector, fldPath.Child("selector"))...)

	if spec.Env == nil && spec.EnvFrom == nil && spec.VolumeMounts == nil && spec.Volumes == nil {
		allErrs = append(allErrs, field.Required(fldPath.Child("volumes", "env", "envFrom", "volumeMounts"), "must specify at least one"))
	}

	vols, vErrs := apivalidation.ValidateVolumes(spec.Volumes, fldPath.Child("volumes"))
	allErrs = append(allErrs, vErrs...)
	allErrs = append(allErrs, apivalidation.ValidateEnv(spec.Env, fldPath.Child("env"))...)
	allErrs = append(allErrs, apivalidation.ValidateEnvFrom(spec.EnvFrom, fldPath.Child("envFrom"))...)
	allErrs = append(allErrs, apivalidation.ValidateVolumeMounts(spec.VolumeMounts, nil, vols, nil, fldPath.Child("volumeMounts"))...)

	return allErrs
}

// ValidatePodPreset validates a PodPreset.
func ValidatePodPreset(pip *settings.PodPreset) field.ErrorList {
	allErrs := apivalidation.ValidateObjectMeta(&pip.ObjectMeta, true, ValidatePodPresetName, field.NewPath("metadata"))
	allErrs = append(allErrs, ValidatePodPresetSpec(&pip.Spec, field.NewPath("spec"))...)
	return allErrs
}

// ValidatePodPresetUpdate tests if required fields in the PodPreset are set.
func ValidatePodPresetUpdate(pip, oldPip *settings.PodPreset) field.ErrorList {
	allErrs := apivalidation.ValidateObjectMetaUpdate(&pip.ObjectMeta, &oldPip.ObjectMeta, field.NewPath("metadata"))
	allErrs = append(allErrs, ValidatePodPresetSpec(&pip.Spec, field.NewPath("spec"))...)

	return allErrs
}
