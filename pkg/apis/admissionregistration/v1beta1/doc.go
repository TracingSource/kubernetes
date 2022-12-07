// +k8s:conversion-gen=k8s.io/kubernetes/pkg/apis/admissionregistration
// +k8s:conversion-gen-external-types=k8s.io/api/admissionregistration/v1beta1
// +k8s:defaulter-gen=TypeMeta
// +k8s:defaulter-gen-input=../../../../vendor/k8s.io/api/admissionregistration/v1beta1
// +groupName=admissionregistration.k8s.io

// Package v1beta1 is the v1beta1 version of the API.
// AdmissionConfiguration and AdmissionPluginConfiguration are legacy static admission plugin configuration
// ValidatingWebhookConfiguration, and MutatingWebhookConfiguration are for the
// new dynamic admission controller configuration.
package v1beta1 // import "k8s.io/kubernetes/pkg/apis/admissionregistration/v1beta1"
