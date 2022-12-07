package v1

import (
	"k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

func addDefaultingFuncs(scheme *runtime.Scheme) error {
	return RegisterDefaults(scheme)
}

func SetDefaults_NetworkPolicyPort(obj *networkingv1.NetworkPolicyPort) {
	// Default any undefined Protocol fields to TCP.
	if obj.Protocol == nil {
		proto := v1.ProtocolTCP
		obj.Protocol = &proto
	}
}

func SetDefaults_NetworkPolicy(obj *networkingv1.NetworkPolicy) {
	if len(obj.Spec.PolicyTypes) == 0 {
		// Any policy that does not specify policyTypes implies at least "Ingress".
		obj.Spec.PolicyTypes = []networkingv1.PolicyType{networkingv1.PolicyTypeIngress}
		if len(obj.Spec.Egress) != 0 {
			obj.Spec.PolicyTypes = append(obj.Spec.PolicyTypes, networkingv1.PolicyTypeEgress)
		}
	}
}
