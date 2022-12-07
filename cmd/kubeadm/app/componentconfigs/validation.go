package componentconfigs

import (
	"k8s.io/apimachinery/pkg/util/validation/field"
	"k8s.io/klog"
	kubeadmapi "k8s.io/kubernetes/cmd/kubeadm/app/apis/kubeadm"
)

// NoValidator returns a dummy validator function when no validation method is available for the component
func NoValidator(component string) func(*kubeadmapi.ClusterConfiguration, *field.Path) field.ErrorList {
	return func(_ *kubeadmapi.ClusterConfiguration, _ *field.Path) field.ErrorList {
		klog.Warningf("Cannot validate %s config - no validator is available", component)
		return field.ErrorList{}
	}
}
