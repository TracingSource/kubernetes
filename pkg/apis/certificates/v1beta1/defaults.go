package v1beta1

import (
	certificatesv1beta1 "k8s.io/api/certificates/v1beta1"
	"k8s.io/apimachinery/pkg/runtime"
)

func addDefaultingFuncs(scheme *runtime.Scheme) error {
	return RegisterDefaults(scheme)
}
func SetDefaults_CertificateSigningRequestSpec(obj *certificatesv1beta1.CertificateSigningRequestSpec) {
	if obj.Usages == nil {
		obj.Usages = []certificatesv1beta1.KeyUsage{certificatesv1beta1.UsageDigitalSignature, certificatesv1beta1.UsageKeyEncipherment}
	}
}
