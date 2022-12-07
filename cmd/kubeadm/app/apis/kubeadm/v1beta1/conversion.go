package v1beta1

import (
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/conversion"
	"k8s.io/kubernetes/cmd/kubeadm/app/apis/kubeadm"
)

func Convert_kubeadm_InitConfiguration_To_v1beta1_InitConfiguration(in *kubeadm.InitConfiguration, out *InitConfiguration, s conversion.Scope) error {
	if err := autoConvert_kubeadm_InitConfiguration_To_v1beta1_InitConfiguration(in, out, s); err != nil {
		return err
	}

	if in.CertificateKey != "" {
		return errors.New("certificateKey field is not supported by v1beta1 config format")
	}

	return nil
}

func Convert_kubeadm_JoinControlPlane_To_v1beta1_JoinControlPlane(in *kubeadm.JoinControlPlane, out *JoinControlPlane, s conversion.Scope) error {
	if err := autoConvert_kubeadm_JoinControlPlane_To_v1beta1_JoinControlPlane(in, out, s); err != nil {
		return err
	}

	if in.CertificateKey != "" {
		return errors.New("certificateKey field is not supported by v1beta1 config format")
	}

	return nil
}

func Convert_kubeadm_NodeRegistrationOptions_To_v1beta1_NodeRegistrationOptions(in *kubeadm.NodeRegistrationOptions, out *NodeRegistrationOptions, s conversion.Scope) error {
	if err := autoConvert_kubeadm_NodeRegistrationOptions_To_v1beta1_NodeRegistrationOptions(in, out, s); err != nil {
		return err
	}

	if len(in.IgnorePreflightErrors) > 0 {
		return errors.New("ignorePreflightErrors field is not supported by v1beta1 config format")
	}

	return nil
}
