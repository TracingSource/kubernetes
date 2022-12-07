package v1

import kruntime "k8s.io/apimachinery/pkg/runtime"

func addDefaultingFuncs(scheme *kruntime.Scheme) error {
	return RegisterDefaults(scheme)
}

func SetDefaults_Configuration(obj *Configuration) {}
