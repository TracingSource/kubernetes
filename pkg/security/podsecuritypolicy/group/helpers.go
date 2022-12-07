package group

import (
	"fmt"

	policy "k8s.io/api/policy/v1beta1"
	"k8s.io/apimachinery/pkg/util/validation/field"
	psputil "k8s.io/kubernetes/pkg/security/podsecuritypolicy/util"
)

func ValidateGroupsInRanges(fldPath *field.Path, ranges []policy.IDRange, groups []int64) field.ErrorList {
	allErrs := field.ErrorList{}

	for _, group := range groups {
		if !isGroupInRanges(group, ranges) {
			detail := fmt.Sprintf("group %d must be in the ranges: %v", group, ranges)
			allErrs = append(allErrs, field.Invalid(fldPath, groups, detail))
		}
	}
	return allErrs
}

func isGroupInRanges(group int64, ranges []policy.IDRange) bool {
	for _, rng := range ranges {
		if psputil.GroupFallsInRange(group, rng) {
			return true
		}
	}
	return false
}
