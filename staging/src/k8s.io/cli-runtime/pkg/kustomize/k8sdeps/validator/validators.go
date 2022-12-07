// Package validator provides functions to validate labels, annotations, namespace using apimachinery
package validator

import (
	"errors"
	apivalidation "k8s.io/apimachinery/pkg/api/validation"
	v1validation "k8s.io/apimachinery/pkg/apis/meta/v1/validation"
	"k8s.io/apimachinery/pkg/util/validation"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

// KustValidator validates Labels and annotations by apimachinery
type KustValidator struct{}

// NewKustValidator returns a KustValidator object
func NewKustValidator() *KustValidator {
	return &KustValidator{}
}

// MakeAnnotationValidator returns a MapValidatorFunc using apimachinery.
func (v *KustValidator) MakeAnnotationValidator() func(map[string]string) error {
	return func(x map[string]string) error {
		errs := apivalidation.ValidateAnnotations(x, field.NewPath("field"))
		if len(errs) > 0 {
			return errors.New(errs.ToAggregate().Error())
		}
		return nil
	}
}

// MakeLabelValidator returns a MapValidatorFunc using apimachinery.
func (v *KustValidator) MakeLabelValidator() func(map[string]string) error {
	return func(x map[string]string) error {
		errs := v1validation.ValidateLabels(x, field.NewPath("field"))
		if len(errs) > 0 {
			return errors.New(errs.ToAggregate().Error())
		}
		return nil
	}
}

// ValidateNamespace validates a string is a valid namespace using apimachinery.
func (v *KustValidator) ValidateNamespace(s string) []string {
	return validation.IsDNS1123Label(s)
}
