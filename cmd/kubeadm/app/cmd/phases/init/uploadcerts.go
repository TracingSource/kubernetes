package phases

import (
	"fmt"

	"github.com/pkg/errors"

	"k8s.io/kubernetes/cmd/kubeadm/app/cmd/options"
	"k8s.io/kubernetes/cmd/kubeadm/app/cmd/phases/workflow"
	cmdutil "k8s.io/kubernetes/cmd/kubeadm/app/cmd/util"
	kubeadmconstants "k8s.io/kubernetes/cmd/kubeadm/app/constants"
	"k8s.io/kubernetes/cmd/kubeadm/app/phases/copycerts"
)

// NewUploadCertsPhase returns the uploadCerts phase
func NewUploadCertsPhase() workflow.Phase {
	return workflow.Phase{
		Name:  "upload-certs",
		Short: fmt.Sprintf("Upload certificates to %s", kubeadmconstants.KubeadmCertsSecret),
		Long:  cmdutil.MacroCommandLongDescription,
		Run:   runUploadCerts,
		InheritFlags: []string{
			options.CfgPath,
			options.UploadCerts,
			options.CertificateKey,
			options.SkipCertificateKeyPrint,
		},
	}
}

func runUploadCerts(c workflow.RunData) error {
	data, ok := c.(InitData)
	if !ok {
		return errors.New("upload-certs phase invoked with an invalid data struct")
	}

	if !data.UploadCerts() {
		fmt.Printf("[upload-certs] Skipping phase. Please see --%s\n", options.UploadCerts)
		return nil
	}
	// 读取 /etc/kubernetes/admin.conf 配置文件构建 kube client
	client, err := data.Client()
	if err != nil {
		return err
	}

	if len(data.CertificateKey()) == 0 {
		certificateKey, err := copycerts.CreateCertificateKey()
		if err != nil {
			return err
		}
		data.SetCertificateKey(certificateKey)
	}

	if err := copycerts.UploadCerts(client, data.Cfg(), data.CertificateKey()); err != nil {
		return errors.Wrap(err, "error uploading certs")
	}
	if !data.SkipCertificateKeyPrint() {
		fmt.Printf("[upload-certs] Using certificate key:\n%s\n", data.CertificateKey())
	}
	return nil
}
