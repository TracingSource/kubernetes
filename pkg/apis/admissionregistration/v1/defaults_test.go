package v1_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	v1 "k8s.io/api/admissionregistration/v1"
	apiequality "k8s.io/apimachinery/pkg/api/equality"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/kubernetes/pkg/api/legacyscheme"
	_ "k8s.io/kubernetes/pkg/apis/admissionregistration/install"
	utilpointer "k8s.io/utils/pointer"
)

func TestDefaultAdmissionWebhook(t *testing.T) {
	fail := v1.Fail
	equivalent := v1.Equivalent
	never := v1.NeverReinvocationPolicy
	ten := int32(10)
	allScopes := v1.AllScopes

	tests := []struct {
		name     string
		original runtime.Object
		expected runtime.Object
	}{
		{
			name: "ValidatingWebhookConfiguration",
			original: &v1.ValidatingWebhookConfiguration{
				Webhooks: []v1.ValidatingWebhook{{}},
			},
			expected: &v1.ValidatingWebhookConfiguration{
				Webhooks: []v1.ValidatingWebhook{{
					FailurePolicy:     &fail,
					MatchPolicy:       &equivalent,
					TimeoutSeconds:    &ten,
					NamespaceSelector: &metav1.LabelSelector{},
					ObjectSelector:    &metav1.LabelSelector{},
				}},
			},
		},
		{
			name: "MutatingWebhookConfiguration",
			original: &v1.MutatingWebhookConfiguration{
				Webhooks: []v1.MutatingWebhook{{}},
			},
			expected: &v1.MutatingWebhookConfiguration{
				Webhooks: []v1.MutatingWebhook{{
					FailurePolicy:      &fail,
					MatchPolicy:        &equivalent,
					ReinvocationPolicy: &never,
					TimeoutSeconds:     &ten,
					NamespaceSelector:  &metav1.LabelSelector{},
					ObjectSelector:     &metav1.LabelSelector{},
				}},
			},
		},
		{
			name: "scope=*",
			original: &v1.MutatingWebhookConfiguration{
				Webhooks: []v1.MutatingWebhook{{
					Rules: []v1.RuleWithOperations{{}},
				}},
			},
			expected: &v1.MutatingWebhookConfiguration{
				Webhooks: []v1.MutatingWebhook{{
					Rules: []v1.RuleWithOperations{{Rule: v1.Rule{
						Scope: &allScopes, // defaulted
					}}},
					FailurePolicy:      &fail,
					MatchPolicy:        &equivalent,
					ReinvocationPolicy: &never,
					TimeoutSeconds:     &ten,
					NamespaceSelector:  &metav1.LabelSelector{},
					ObjectSelector:     &metav1.LabelSelector{},
				}},
			},
		},
		{
			name: "port=443",
			original: &v1.MutatingWebhookConfiguration{
				Webhooks: []v1.MutatingWebhook{{
					ClientConfig: v1.WebhookClientConfig{
						Service: &v1.ServiceReference{},
					},
				}},
			},
			expected: &v1.MutatingWebhookConfiguration{
				Webhooks: []v1.MutatingWebhook{{
					ClientConfig: v1.WebhookClientConfig{
						Service: &v1.ServiceReference{
							Port: utilpointer.Int32Ptr(443), // defaulted
						},
					},
					FailurePolicy:      &fail,
					MatchPolicy:        &equivalent,
					ReinvocationPolicy: &never,
					TimeoutSeconds:     &ten,
					NamespaceSelector:  &metav1.LabelSelector{},
					ObjectSelector:     &metav1.LabelSelector{},
				}},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			original := test.original
			expected := test.expected
			legacyscheme.Scheme.Default(original)
			if !apiequality.Semantic.DeepEqual(original, expected) {
				t.Error(cmp.Diff(expected, original))
			}
		})
	}
}
