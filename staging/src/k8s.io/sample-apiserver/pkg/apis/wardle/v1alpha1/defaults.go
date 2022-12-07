package v1alpha1

import (
	"k8s.io/apimachinery/pkg/runtime"
)

func addDefaultingFuncs(scheme *runtime.Scheme) error {
	return RegisterDefaults(scheme)
}

func SetDefaults_FlunderSpec(obj *FlunderSpec) {
	if (obj.ReferenceType == nil || len(*obj.ReferenceType) == 0) && len(obj.Reference) != 0 {
		t := FlunderReferenceType
		obj.ReferenceType = &t
	}
}
