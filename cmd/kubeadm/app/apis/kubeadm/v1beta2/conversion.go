package v1beta2

import (
	conversion "k8s.io/apimachinery/pkg/conversion"
	kubeadm "k8s.io/kubernetes/cmd/kubeadm/app/apis/kubeadm"
)

func Convert_kubeadm_InitConfiguration_To_v1beta2_InitConfiguration(in *kubeadm.InitConfiguration, out *InitConfiguration, s conversion.Scope) error {
	return autoConvert_kubeadm_InitConfiguration_To_v1beta2_InitConfiguration(in, out, s)
}

func Convert_v1beta2_InitConfiguration_To_kubeadm_InitConfiguration(in *InitConfiguration, out *kubeadm.InitConfiguration, s conversion.Scope) error {
	err := autoConvert_v1beta2_InitConfiguration_To_kubeadm_InitConfiguration(in, out, s)
	if err != nil {
		return err
	}

	// Keep the fuzzer test happy by setting out.ClusterConfiguration to defaults
	clusterCfg := &ClusterConfiguration{}
	SetDefaults_ClusterConfiguration(clusterCfg)
	return Convert_v1beta2_ClusterConfiguration_To_kubeadm_ClusterConfiguration(clusterCfg, &out.ClusterConfiguration, s)
}
