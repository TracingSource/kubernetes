package phases

import (
	"k8s.io/component-base/version"
	kubeadmapiv1beta2 "k8s.io/kubernetes/cmd/kubeadm/app/apis/kubeadm/v1beta2"
)

// SetKubernetesVersion gets the current Kubeadm version and sets it as KubeadmVersion in the config,
// unless it's already set to a value different from the default.
func SetKubernetesVersion(cfg *kubeadmapiv1beta2.ClusterConfiguration) {

	if cfg.KubernetesVersion != kubeadmapiv1beta2.DefaultKubernetesVersion && cfg.KubernetesVersion != "" {
		return
	}
	cfg.KubernetesVersion = version.Get().String()
}
