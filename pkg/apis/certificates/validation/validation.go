package validation

import (
	"fmt"

	"k8s.io/apimachinery/pkg/util/validation/field"
	"k8s.io/kubernetes/pkg/apis/certificates"
	apivalidation "k8s.io/kubernetes/pkg/apis/core/validation"
)

// validateCSR validates the signature and formatting of a base64-wrapped,
// PEM-encoded PKCS#10 certificate signing request. If this is invalid, we must
// not accept the CSR for further processing.
func validateCSR(obj *certificates.CertificateSigningRequest) error {
	csr, err := certificates.ParseCSR(obj)
	if err != nil {
		return err
	}
	// check that the signature is valid
	err = csr.CheckSignature()
	if err != nil {
		return err
	}
	return nil
}

// We don't care what you call your certificate requests.
func ValidateCertificateRequestName(name string, prefix bool) []string {
	return nil
}

func ValidateCertificateSigningRequest(csr *certificates.CertificateSigningRequest) field.ErrorList {
	isNamespaced := false
	allErrs := apivalidation.ValidateObjectMeta(&csr.ObjectMeta, isNamespaced, ValidateCertificateRequestName, field.NewPath("metadata"))
	err := validateCSR(csr)

	specPath := field.NewPath("spec")

	if err != nil {
		allErrs = append(allErrs, field.Invalid(specPath.Child("request"), csr.Spec.Request, fmt.Sprintf("%v", err)))
	}
	if len(csr.Spec.Usages) == 0 {
		allErrs = append(allErrs, field.Required(specPath.Child("usages"), "usages must be provided"))
	}
	return allErrs
}

func ValidateCertificateSigningRequestUpdate(newCSR, oldCSR *certificates.CertificateSigningRequest) field.ErrorList {
	validationErrorList := ValidateCertificateSigningRequest(newCSR)
	metaUpdateErrorList := apivalidation.ValidateObjectMetaUpdate(&newCSR.ObjectMeta, &oldCSR.ObjectMeta, field.NewPath("metadata"))
	return append(validationErrorList, metaUpdateErrorList...)
}
