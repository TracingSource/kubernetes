package v1beta1

import (
	v1 "k8s.io/api/core/v1"
	discoveryv1beta1 "k8s.io/api/discovery/v1beta1"
	"k8s.io/apimachinery/pkg/runtime"
)

var (
	defaultPortName = ""
	defaultProtocol = v1.ProtocolTCP
)

func addDefaultingFuncs(scheme *runtime.Scheme) error {
	return RegisterDefaults(scheme)
}

func SetDefaults_EndpointPort(obj *discoveryv1beta1.EndpointPort) {
	if obj.Name == nil {
		obj.Name = &defaultPortName
	}

	if obj.Protocol == nil {
		obj.Protocol = &defaultProtocol
	}
}
